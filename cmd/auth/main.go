package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	mgo "gopkg.in/mgo.v2"

	"datapace"
	"datapace/auth"
	"datapace/auth/api"
	grpcapi "datapace/auth/api/grpc"
	httpapi "datapace/auth/api/http"
	"datapace/auth/bcrypt"
	"datapace/auth/jwt"
	"datapace/auth/mongo"
	"datapace/auth/transactions"
	log "datapace/logger"
	transactionsapi "datapace/transactions/api/grpc"
)

const (
	envHTTPPort        = "DATAPACE_AUTH_HTTP_PORT"
	envGRPCPort        = "DATAPACE_AUTH_GRPC_PORT"
	envDBURL           = "DATAPACE_AUTH_DB_URL"
	envDBUser          = "DATAPACE_AUTH_DB_USER"
	envDBPass          = "DATAPACE_AUTH_DB_PASS"
	envDBName          = "DATAPACE_AUTH_DB_NAME"
	envTransactionsURL = "DATAPACE_TRANSACTIONS_URL"
	envSecret          = "DATAPACE_AUTH_SECRET"

	defHTTPPort        = "8080"
	defGRPCPort        = "8081"
	defDBURL           = "0.0.0.0"
	defDBUser          = ""
	defDBPass          = ""
	defDBName          = "auth"
	defTransactionsURL = "localhost:8081"
	defSecret          = "datapace"

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
	transactionsURL  string
	secret           string
}

func main() {
	cfg := loadConfig()

	logger := log.New(os.Stdout)

	ms := connectToDB(cfg, logger)
	defer ms.Close()

	conn := newGRPCConn(cfg.transactionsURL, logger)
	defer conn.Close()

	tc := transactionsapi.NewClient(conn)

	svc := newService(cfg, ms, tc, logger)

	errs := make(chan error, 2)

	go startHTTPServer(svc, cfg.httpPort, logger, errs)

	go startGRPCServer(svc, cfg.grpcPort, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err := <-errs
	logger.Error(fmt.Sprintf("Auth service terminated: %s", err))
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
		transactionsURL:  datapace.Env(envTransactionsURL, defTransactionsURL),
		secret:           datapace.Env(envSecret, defSecret),
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

func newGRPCConn(transactionsURL string, logger log.Logger) *grpc.ClientConn {
	conn, err := grpc.Dial(transactionsURL, grpc.WithInsecure())
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to transactions service: %s", err))
		os.Exit(1)
	}

	return conn
}

func newService(cfg config, ms *mgo.Session, tc datapace.TransactionsServiceClient, logger log.Logger) auth.Service {
	users := mongo.NewUserRepository(ms)
	hasher := bcrypt.New()
	idp := jwt.New(cfg.secret)
	ts := transactions.NewService(tc)

	svc := auth.New(users, hasher, idp, ts)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "auth",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "auth",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)

	return svc
}

func startHTTPServer(svc auth.Service, port string, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("Users HTTP service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, httpapi.MakeHandler(svc))
}

func startGRPCServer(svc auth.Service, port string, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", p)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to listen on port %s: %s", port, err))
	}

	server := grpc.NewServer()
	datapace.RegisterAuthServiceServer(server, grpcapi.NewServer(svc))
	logger.Info(fmt.Sprintf("Auth gRPC service started, exposed port %s", port))
	errs <- server.Serve(listener)
}
