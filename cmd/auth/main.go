package main

import (
	"fmt"
	"monetasa"
	"monetasa/auth"
	"monetasa/auth/api"
	grpcapi "monetasa/auth/api/grpc"
	httpapi "monetasa/auth/api/http"
	"monetasa/auth/bcrypt"
	"monetasa/auth/jwt"
	"monetasa/auth/mongo"
	"monetasa/auth/transactions"
	log "monetasa/logger"
	transactionsapi "monetasa/transactions/api/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	mgo "gopkg.in/mgo.v2"
)

const (
	envHTTPPort        = "MONETASA_AUTH_HTTP_PORT"
	envGRPCPort        = "MONETASA_AUTH_GRPC_PORT"
	envMongoURL        = "MONETASA_AUTH_MONGO_URL"
	envMongoUser       = "MONETASA_AUTH_MONGO_USER"
	envMongoPass       = "MONETASA_AUTH_MONGO_PASS"
	envMongoDatabase   = "MONETASA_AUTH_MONGO_DB"
	envTransactionsURL = "MONETASA_TRANSACTIONS_URL"
	envSecret          = "MONETASA_AUTH_SECRET"

	defHTTPPort        = "8080"
	defGRPCPort        = "8081"
	defMongoURL        = "0.0.0.0"
	defMongoUser       = ""
	defMongoPass       = ""
	defMongoDatabase   = "auth"
	defTransactionsURL = "localhost:8081"
	defSecret          = "monetasa"

	mongoConnectTimeout = 5000
	mongoSocketTimeout  = 5000
)

type config struct {
	httpPort            string
	grpcPort            string
	mongoURL            string
	mongoUser           string
	mongoPass           string
	mongoDatabase       string
	mongoConnectTimeout int
	mongoSocketTimeout  int
	transactionsURL     string
	secret              string
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
		httpPort:            monetasa.Env(envHTTPPort, defHTTPPort),
		grpcPort:            monetasa.Env(envGRPCPort, defGRPCPort),
		mongoURL:            monetasa.Env(envMongoURL, defMongoURL),
		mongoUser:           monetasa.Env(envMongoUser, defMongoUser),
		mongoPass:           monetasa.Env(envMongoPass, defMongoPass),
		mongoDatabase:       monetasa.Env(envMongoDatabase, defMongoDatabase),
		mongoConnectTimeout: mongoConnectTimeout,
		mongoSocketTimeout:  mongoSocketTimeout,
		transactionsURL:     monetasa.Env(envTransactionsURL, defTransactionsURL),
		secret:              monetasa.Env(envSecret, defSecret),
	}
}

func connectToDB(cfg config, logger log.Logger) *mgo.Session {
	ms, err := mongo.Connect(
		cfg.mongoURL,
		cfg.mongoConnectTimeout,
		cfg.mongoSocketTimeout,
		cfg.mongoDatabase,
		cfg.mongoUser,
		cfg.mongoPass,
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

func newService(cfg config, ms *mgo.Session, tc monetasa.TransactionsServiceClient, logger log.Logger) auth.Service {
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
	monetasa.RegisterAuthServiceServer(server, grpcapi.NewServer(svc))
	logger.Info(fmt.Sprintf("Auth gRPC service started, exposed port %s", port))
	errs <- server.Serve(listener)
}
