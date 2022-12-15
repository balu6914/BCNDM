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

	"github.com/datapace/datapace/executions/argo"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	mgo "gopkg.in/mgo.v2"

	"github.com/datapace/datapace"

	authapi "github.com/datapace/datapace/auth/api/grpc"
	"github.com/datapace/datapace/executions"
	"github.com/datapace/datapace/executions/api"
	grpcapi "github.com/datapace/datapace/executions/api/grpc"
	httpapi "github.com/datapace/datapace/executions/api/http"
	"github.com/datapace/datapace/executions/kubeflow"
	"github.com/datapace/datapace/executions/mongo"
	"github.com/datapace/datapace/executions/subscriptions"
	"github.com/datapace/datapace/executions/wwh"
	"github.com/datapace/datapace/logger"
	authproto "github.com/datapace/datapace/proto/auth"
	executionproto "github.com/datapace/datapace/proto/executions"
)

const (
	envHTTPPort           = "DATAPACE_EXECUTIONS_HTTP_PORT"
	envGRPCPort           = "DATAPACE_EXECUTIONS_GRPC_PORT"
	envDBURL              = "DATAPACE_EXECUTIONS_DB_URL"
	envDBUser             = "DATAPACE_EXECUTIONS_DB_USER"
	envDBPass             = "DATAPACE_EXECUTIONS_DB_PASS"
	envDBName             = "DATAPACE_EXECUTIONS_DB_NAME"
	envAuthURL            = "DATAPACE_AUTH_URL"
	envSubscriptionsURL   = "DATAPACE_SUBSCRIPTIONS_URL"
	envStreamsURL         = "DATAPACE_STREAMS_URL"
	envWWHCatalogURL      = "DATAPACE_WWH_CATALOG_URL"
	envWWHDaemonURL       = "DATAPACE_WWH_DAEMON_URL"
	envWWHToken           = "DATAPACE_WWH_TOKEN"
	envWWHUsername        = "DATAPACE_WWH_USERNAME"
	envWWHPassword        = "DATAPACE_WWH_PASSWORD"
	envAISystem           = "DATAPACE_AI_SYSTEM"
	envKFURL              = "DATAPACE_KUBEFLOW_URL"
	envKFStatusInterval   = "DATAPACE_KUBEFLOW_STATUS_INTERVAL"
	envArgoURL            = "DATAPACE_ARGO_URL"
	envArgoStatusInterval = "DATAPACE_KUBEFLOW_STATUS_INTERVAL"

	defHTTPPort           = "8080"
	defGRPCPort           = "8081"
	defDBURL              = "0.0.0.0"
	defDBUser             = ""
	defDBPass             = ""
	defDBName             = "executions"
	defAuthURL            = "localhost:8081"
	defSubscriptionsURL   = "localhost:8086"
	defStreamsURL         = "localhost:8084"
	defWWHCatalogURL      = "http://localhost:31222"
	defWWHDaemonURL       = "http://localhost:32222"
	defWWHToken           = ""
	defWWHUsername        = ""
	defWWHPassword        = ""
	defAISystem           = "kubeflow"
	defKFURL              = "https://ar.k9s.datapace.io"
	defKFStatusInterval   = "10" // in seconds
	defArgoURL            = "https://argo.datapace.io"
	defArgoStatusInterval = "10" // in seconds

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
	subscriptionsURL string
	streamsURL       string
	wwhCatalogURL    string
	wwhDaemonURL     string
	wwhToken         string
	wwhUsername      string
	wwhPassword      string
	aiSystem         string
	kfURL            string
	kfInterval       time.Duration
	argoURL          string
	argoInterval     time.Duration
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
	kubeflowStatusInterval := datapace.Env(envKFStatusInterval, defKFStatusInterval)
	kfStatusInterval, err := strconv.ParseInt(kubeflowStatusInterval, 10, 64)
	if err != nil {
		log.Fatal(err.Error())
	}

	argoServerStatusInterval := datapace.Env(envArgoStatusInterval, defArgoStatusInterval)
	argoStatusInterval, err := strconv.ParseInt(argoServerStatusInterval, 10, 64)
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
		subscriptionsURL: datapace.Env(envSubscriptionsURL, defSubscriptionsURL),
		streamsURL:       datapace.Env(envStreamsURL, defStreamsURL),
		wwhCatalogURL:    datapace.Env(envWWHCatalogURL, defWWHCatalogURL),
		wwhDaemonURL:     datapace.Env(envWWHDaemonURL, defWWHDaemonURL),
		wwhToken:         datapace.Env(envWWHToken, defWWHToken),
		wwhUsername:      datapace.Env(envWWHUsername, defWWHUsername),
		wwhPassword:      datapace.Env(envWWHPassword, defWWHPassword),
		aiSystem:         datapace.Env(envAISystem, defAISystem),
		kfURL:            datapace.Env(envKFURL, defKFURL),
		kfInterval:       time.Duration(kfStatusInterval) * time.Second,
		argoURL:          datapace.Env(envArgoURL, defArgoURL),
		argoInterval:     time.Duration(argoStatusInterval) * time.Second,
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
	paths := subscriptions.New(cfg.subscriptionsURL)
	//paths := streams.New(cfg.streamsURL, cfg )
	var ai executions.AIService
	switch cfg.aiSystem {
	case "kubeflow":
		ai = kubeflow.New(cfg.kfURL, cfg.kfInterval, logger)
	case "wwf":
		ai = wwh.NewAIService(cfg.wwhCatalogURL, cfg.wwhDaemonURL, cfg.wwhToken, cfg.wwhUsername, cfg.wwhPassword)
	case "argo":
		ai = argo.New(cfg.argoURL, cfg.argoInterval, logger)
	default:
		logger.Error(fmt.Sprintf("Invalid AI system: %s", cfg.aiSystem))
		os.Exit(1)
	}
	svc := executions.NewService(execs, algos, data, ai, paths)
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

func startHTTPServer(svc executions.Service, ac authproto.AuthServiceClient, port string, logger logger.Logger, errs chan error) {
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
	executionproto.RegisterExecutionsServiceServer(server, grpcapi.NewServer(svc))
	logger.Info(fmt.Sprintf("Executions gRPC service started, exposed port %s", port))
	errs <- server.Serve(listener)
}
