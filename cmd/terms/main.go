package main

import (
	"fmt"
	"github.com/datapace/datapace"
	authapi "github.com/datapace/datapace/auth/api/grpc"
	log "github.com/datapace/datapace/logger"
	authproto "github.com/datapace/datapace/proto/auth"
	termsproto "github.com/datapace/datapace/proto/terms"
	"github.com/datapace/datapace/terms"
	"github.com/datapace/datapace/terms/api"
	grpcapi "github.com/datapace/datapace/terms/api/grpc"
	"github.com/datapace/datapace/terms/mongo"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"gopkg.in/mgo.v2"
	"net"
	"os"
)

const (
	envHTTPPort = "DATAPACE_TERMS_HTTP_PORT"
	envGRPCPort = "DATAPACE_TERMS_GRPC_PORT"
	envDBURL    = "DATAPACE_TERMS_DB_URL"
	envDBUser   = "DATAPACE_TERMS_DB_USER"
	envDBPass   = "DATAPACE_TERMS_DB_PASS"
	envDBName   = "DATAPACE_TERMS_DB_NAME"
	envAuthURL  = "DATAPACE_AUTH_URL"

	defHTTPPort      = "8080"
	defGRPCPort      = "8081"
	defDBURL         = "0.0.0.0"
	defDBUser        = ""
	defDBPass        = ""
	defDBName        = "terms"
	defAuthURL       = "localhost:8081"
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
}

func main() {
	cfg := loadConfig()

	logger := log.New(os.Stdout)

	ms := connectToDB(cfg, logger)
	defer ms.Close()
	aconn := newGRPCConn(cfg.authURL, logger)
	defer aconn.Close()

	auth := authapi.NewClient(aconn)
	svc := newService(auth, ms, logger)
	errs := make(chan error, 2)
	go startGRPCServer(svc, cfg.grpcPort, logger, errs)
	err := <-errs
	logger.Error(fmt.Sprintf("Terms service terminated: %s", err))
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

func loadConfig() config {
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
	}
}

func newGRPCConn(addr string, logger log.Logger) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to GRPC service on address %s: %s", addr, err))
		os.Exit(1)
	}
	return conn
}

func newService(auth authproto.AuthServiceClient, ms *mgo.Session, logger log.Logger) terms.Service {

	repo := mongo.NewTermsRepository(ms)
	svc := terms.New(auth, repo)

	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "terms",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "terms",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)

	return svc
}

func startGRPCServer(svc terms.Service, port string, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", p)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to listen on port %s: %s", port, err))
	}

	server := grpc.NewServer()
	termsproto.RegisterTermsServiceServer(server, grpcapi.NewServer(svc))
	logger.Info(fmt.Sprintf("Terms gRPC service started, exposed port %s", port))
	errs <- server.Serve(listener)
}
