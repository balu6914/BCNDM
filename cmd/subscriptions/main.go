package main

import (
	"fmt"
	"monetasa"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	authapi "monetasa/auth/api/grpc"
	log "monetasa/logger"
	"monetasa/subscriptions"
	"monetasa/subscriptions/api"
	"monetasa/subscriptions/mongo"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"

	"google.golang.org/grpc"
)

const (
	envSubsPort      = "MONETASA_SUBSCRIPTIONS_PORT"
	envAuthURL       = "MONETASA_AUTH_URL"
	envMongoURL      = "MONETASA_SUBSCRIPTIONS_DB_URL"
	envMongoUser     = "MONETASA_SUBSCRIPTIONS_DB_USER"
	envMongoPass     = "MONETASA_SUBSCRIPTIONS_DB_PASS"
	envMongoDatabase = "MONETASA_SUBSCRIPTIONS_DB_NAME"

	defSubsPort      = "8080"
	defSubsURL       = "0.0.0.0"
	defMongoURL      = "0.0.0.0"
	defMongoUser     = ""
	defMongoPass     = ""
	defMongoDatabase = "subscriptions"
	defAuthURL       = "localhost:8080"

	defMongoConnectTimeout = 5000
	defMongoSocketTimeout  = 5000
)

type config struct {
	SubsPort            string
	MongoURL            string
	MongoUser           string
	MongoPass           string
	MongoDatabase       string
	MongoConnectTimeout int
	MongoSocketTimeout  int
	AuthURL             string
}

func main() {
	cfg := config{
		SubsPort:            monetasa.Env(envSubsPort, defSubsPort),
		MongoURL:            monetasa.Env(envMongoURL, defMongoURL),
		MongoUser:           monetasa.Env(envMongoUser, defMongoUser),
		MongoPass:           monetasa.Env(envMongoPass, defMongoPass),
		MongoDatabase:       monetasa.Env(envMongoDatabase, defMongoDatabase),
		MongoConnectTimeout: defMongoConnectTimeout,
		MongoSocketTimeout:  defMongoSocketTimeout,
		AuthURL:             monetasa.Env(envAuthURL, defAuthURL),
	}

	logger := log.New(os.Stdout)

	ms, err := mongo.Connect(
		cfg.MongoURL,
		cfg.MongoConnectTimeout,
		cfg.MongoSocketTimeout,
		cfg.MongoDatabase,
		cfg.MongoUser,
		cfg.MongoPass,
	)
	if err != nil {
		logger.Error("Failed to connect to Mongo.")
		os.Exit(1)
	}
	defer ms.Close()

	conn, err := grpc.Dial(defAuthURL, grpc.WithInsecure())
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to users service: %s", err))
		os.Exit(1)
	}
	defer conn.Close()

	ac := authapi.NewClient(conn)

	subr := mongo.NewRepository(ms)
	svc := subscriptions.New(subr)

	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "subscriptions",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "subscriptions",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)

	errs := make(chan error, 2)

	go func() {
		p := fmt.Sprintf(":%s", cfg.SubsPort)
		logger.Info(fmt.Sprintf("Subscriptions service started, exposed port %s", cfg.SubsPort))
		errs <- http.ListenAndServe(p, api.MakeHandler(svc, ac))
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("Subscriptions service terminated: %s", err))
}
