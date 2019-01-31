package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	mgo "gopkg.in/mgo.v2"

	"datapace"
	authapi "datapace/auth/api/grpc"
	"datapace/executions"
	"datapace/executions/api"
	"datapace/executions/mongo"
	log "datapace/logger"
)

const (
	envHTTPPort = "DATAPACE_EXECUTIONS_HTTP_PORT"
	envDBURL    = "DATAPACE_EXECUTIONS_DB_URL"
	envDBUser   = "DATAPACE_EXECUTIONS_DB_USER"
	envDBPass   = "DATAPACE_EXECUTIONS_DB_PASS"
	envDBName   = "DATAPACE_EXECUTIONS_DB_NAME"
	envAuthURL  = "DATAPACE_AUTH_URL"

	defHTTPPort = "8080"
	defDBURL    = "0.0.0.0"
	defDBUser   = ""
	defDBPass   = ""
	defDBName   = "executions"
	defAuthURL  = "localhost:8081"

	dbConnectTimeout = 5000
	dbSocketTimeout  = 5000
)

type config struct {
	httpPort         string
	dbURL            string
	dbUser           string
	dbPass           string
	dbName           string
	dbConnectTimeout int
	dbSocketTimeout  int
	authURL          string
}

func main() {
	cfg := loadConfig()

	logger := log.New(os.Stdout)

	ms := connectToDB(cfg, logger)
	defer ms.Close()

	authConn := newGRPCConn(cfg.authURL, logger)
	defer authConn.Close()

	ac := authapi.NewClient(authConn)

	svc := newService(cfg, ms, logger)

	errs := make(chan error, 2)

	go startHTTPServer(svc, ac, cfg.httpPort, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err := <-errs
	logger.Error(fmt.Sprintf("Executions service terminated: %s", err))
}

func loadConfig() config {
	return config{
		httpPort:         datapace.Env(envHTTPPort, defHTTPPort),
		dbURL:            datapace.Env(envDBURL, defDBURL),
		dbUser:           datapace.Env(envDBUser, defDBUser),
		dbPass:           datapace.Env(envDBPass, defDBPass),
		dbName:           datapace.Env(envDBName, defDBName),
		dbConnectTimeout: dbConnectTimeout,
		dbSocketTimeout:  dbSocketTimeout,
		authURL:          datapace.Env(envAuthURL, defAuthURL),
	}
}

func connectToDB(cfg config, logger log.Logger) *mgo.Session {
	ms, err := mongo.Connect(
		cfg.dbURL,
		cfg.dbConnectTimeout,
		cfg.dbSocketTimeout,
		cfg.dbName,
		cfg.dbUser,
		cfg.dbPass,
	)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to Mongo: %s", err))
		os.Exit(1)
	}

	return ms
}

func newGRPCConn(authURL string, logger log.Logger) *grpc.ClientConn {
	conn, err := grpc.Dial(authURL, grpc.WithInsecure())
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to auth service: %s", err))
		os.Exit(1)
	}

	return conn
}

func newService(cfg config, ms *mgo.Session, logger log.Logger) executions.Service {
	repo := mongo.New(ms)
	svc := executions.NewService(repo)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "executions",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "executions",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)
	return svc
}

func startHTTPServer(svc executions.Service, ac datapace.AuthServiceClient, port string, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("Executions HTTP service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, api.MakeHandler(svc, ac))
}
