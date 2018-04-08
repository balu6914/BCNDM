package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"monetasa/auth"
	"monetasa/auth/mongo"
	"monetasa/auth/api"
	"monetasa/auth/mongo"
	"monetasa/auth/bcrypt"
	"monetasa/auth/jwt"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	port           					int    = 8080
	defMongoURL            	string = "0.0.0.0"
	defMongoUser       			string = ""
	defMongoPass       			string = ""
	defMongoDatabase       	string = "auth"
	defMongoPort           	int    = 27017
	defMongoConnectTimeout 	int    = 5000
	defMongoSocketTimeout  	int    = 5000

	envMongoURL string = "MONETASA_AUTH_MONGO_URL"
	envAuthURL string = "MONETASA_AUTH_URL"
)

type config struct {
	Port        int
	AuthURL     string
	PostgresURL string

	MongoURL            string
	MongoUser           string
	MongoPass           string
	MongoDatabase       string
	MongoPort           int
	MongoConnectTimeout int
	MongoSocketTimeout  int
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func main() {
	cfg := config{
		Port: 								port,
		AuthURL: 							getenv(envMongoURL, defMongoURL),
		MongoURL:							getenv(envAuthURL, defAuthURL),
		MongoUser: 						defMongoUser,
		MongoPass: 						defMongoPass,
		MongoDatabase: 				defMongoDatabase,
		MongoPort: 						defMongoPort,
		MongoConnectTimeout:	defMongoConnectTimeout,
		MongoSocketTimeout: 	defMongoSocketTimeout,
	}

	var logger log.Logger
	logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	ms, err := mongo.Connect(cfg.MongoURL, cfg.MongoConnectTimeout, cfg.MongoSocketTimeout,
														cfg.MongoDatabase, cfg.MongoUser, cfg.MongoPass)
	if err != nil {
		logger.Error("Failed to connect to Mongo.", zap.Error(err))
		os.Exit(1)
	}
	defer ms.Close()

	users := mongo.NewUserRepository(ms)
	hasher := bcrypt.New()
	idp := jwt.New(cfg.Secret)

	svc := manager.New(users, hasher, idp)

	var svc auth.Service
	svc = api.NewLoggingService(logger, svc)

	fields := []string{"method"}
	svc = api.NewMetricService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "auth",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fields),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "auth",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fields),
		svc,
	)

	errs := make(chan error, 2)

	go func() {
		p := fmt.Sprintf(":%d", cfg.Port)
		errs <- http.ListenAndServe(p, api.MakeHandler(svc))
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}
