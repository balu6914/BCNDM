package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	mgo "gopkg.in/mgo.v2"

	"datapace"
	authapi "datapace/auth/api/grpc"
	"datapace/executions"
	"datapace/executions/api"
	grpcapi "datapace/executions/api/grpc"
	httpapi "datapace/executions/api/http"
	"datapace/executions/kubeflow"
	"datapace/executions/mongo"
	"datapace/executions/wwh"
	"datapace/logger"
)

const (
	envHTTPPort         = "DATAPACE_EXECUTIONS_HTTP_PORT"
	envGRPCPort         = "DATAPACE_EXECUTIONS_GRPC_PORT"
	envDBURL            = "DATAPACE_EXECUTIONS_DB_URL"
	envDBUser           = "DATAPACE_EXECUTIONS_DB_USER"
	envDBPass           = "DATAPACE_EXECUTIONS_DB_PASS"
	envDBName           = "DATAPACE_EXECUTIONS_DB_NAME"
	envAuthURL          = "DATAPACE_AUTH_URL"
	envWWHCatalogURL    = "DATAPACE_WWH_CATALOG_URL"
	envWWHDaemonURL     = "DATAPACE_WWH_DAEMON_URL"
	envWWHToken         = "DATAPACE_WWH_TOKEN"
	envWWHUsername      = "DATAPACE_WWH_USERNAME"
	envWWHPassword      = "DATAPACE_WWH_PASSWORD"
	envKFFlag           = "DATAPACE_KUBEFLOW_ACTIVE"
	envKFURL            = "DATAPACE_KUBEFLOW_URL"
	envKFStatusInterval = "DATAPACE_KUBEFLOW_STATUS_INTERVAL"

	defHTTPPort         = "8080"
	defGRPCPort         = "8081"
	defDBURL            = "0.0.0.0"
	defDBUser           = ""
	defDBPass           = ""
	defDBName           = "executions"
	defAuthURL          = "localhost:8081"
	defWWHCatalogURL    = "http://localhost:31222"
	defWWHDaemonURL     = "http://localhost:32222"
	defWWHToken         = ""
	defWWHUsername      = ""
	defWWHPassword      = ""
	defKFFlag           = "true"
	defKFURL            = "https://ar.k9s.datapace.io"
	defKFStatusInterval = "10" // in seconds

	dbConnectTimeout = 5000
	dbSocketTimeout  = 5000
)

type config struct {
	httpPort         string
	grpcPort         string
	dbURL            string
	dbUser           string
	dbPass           string
	dbName           string
	dbConnectTimeout int
	dbSocketTimeout  int
	authURL          string
	wwhCatalogURL    string
	wwhDaemonURL     string
	wwhToken         string
	wwhUsername      string
	wwhPassword      string
	kfFlag           bool
	kfURL            string
	kfInterval       time.Duration
}

func main() {
	cfg := loadConfig()

	logger := logger.New(os.Stdout)

	ms := connectToDB(cfg, logger)
	defer ms.Close()

	authConn := newGRPCConn(cfg.authURL, logger)
	defer authConn.Close()

	ac := authapi.NewClient(authConn)

	svc := newService(cfg, ms, logger)

	errs := make(chan error, 2)

	go func(errs chan error) {
		if err := svc.ProcessEvents(); err != nil {
			errs <- err
		}
	}(errs)

	go startHTTPServer(svc, ac, cfg.httpPort, logger, errs)

	go startGRPCServer(svc, cfg.grpcPort, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err := <-errs
	logger.Error(fmt.Sprintf("Executions service terminated: %s", err))
}

func loadConfig() config {
	kubeflowFlag := datapace.Env(envKFFlag, defKFFlag)
	kfFlag, err := strconv.ParseBool(kubeflowFlag)
	if err != nil {
		log.Fatal(err.Error())
	}

	kubeflowStatusInterval := datapace.Env(envKFStatusInterval, defKFStatusInterval)
	kfStatusInterval, err := strconv.ParseInt(kubeflowStatusInterval, 10, 64)
	if err != nil {
		log.Fatal(err.Error())
	}

	return config{
		httpPort:         datapace.Env(envHTTPPort, defHTTPPort),
		grpcPort:         datapace.Env(envGRPCPort, defGRPCPort),
		dbURL:            datapace.Env(envDBURL, defDBURL),
		dbUser:           datapace.Env(envDBUser, defDBUser),
		dbPass:           datapace.Env(envDBPass, defDBPass),
		dbName:           datapace.Env(envDBName, defDBName),
		dbConnectTimeout: dbConnectTimeout,
		dbSocketTimeout:  dbSocketTimeout,
		authURL:          datapace.Env(envAuthURL, defAuthURL),
		wwhCatalogURL:    datapace.Env(envWWHCatalogURL, defWWHCatalogURL),
		wwhDaemonURL:     datapace.Env(envWWHDaemonURL, defWWHDaemonURL),
		wwhToken:         datapace.Env(envWWHToken, defWWHToken),
		wwhUsername:      datapace.Env(envWWHUsername, defWWHUsername),
		wwhPassword:      datapace.Env(envWWHPassword, defWWHPassword),
		kfFlag:           kfFlag,
		kfURL:            datapace.Env(envKFURL, defKFURL),
		kfInterval:       time.Duration(kfStatusInterval) * time.Second,
	}
}

func connectToDB(cfg config, logger logger.Logger) *mgo.Session {
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

func newGRPCConn(authURL string, logger logger.Logger) *grpc.ClientConn {
	conn, err := grpc.Dial(authURL, grpc.WithInsecure())
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to auth service: %s", err))
		os.Exit(1)
	}

	return conn
}

func newService(cfg config, ms *mgo.Session, logger logger.Logger) executions.Service {
	execs := mongo.NewExecutionRepository(ms)
	algos := mongo.NewAlgorithmRepository(ms)
	data := mongo.NewDatasetRepository(ms)
	var ai executions.AIService
	if cfg.kfFlag {
		ai = kubeflow.New(cfg.kfURL, cfg.kfInterval, logger)
	} else {
		ai = wwh.NewAIService(cfg.wwhCatalogURL, cfg.wwhDaemonURL, cfg.wwhToken, cfg.wwhUsername, cfg.wwhPassword)
	}
	svc := executions.NewService(execs, algos, data, ai)
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

func startHTTPServer(svc executions.Service, ac datapace.AuthServiceClient, port string, logger logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("Executions HTTP service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, httpapi.MakeHandler(svc, ac))
}

func startGRPCServer(svc executions.Service, port string, logger logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", p)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to listen on port %s: %s", port, err))
	}

	server := grpc.NewServer()
	datapace.RegisterExecutionsServiceServer(server, grpcapi.NewServer(svc))
	logger.Info(fmt.Sprintf("Executions gRPC service started, exposed port %s", port))
	errs <- server.Serve(listener)
}
