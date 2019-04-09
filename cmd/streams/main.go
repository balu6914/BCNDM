package main

import (
	"datapace"
	accessapi "datapace/access-control/api/grpc"
	authapi "datapace/auth/api/grpc"
	executionsapi "datapace/executions/api/grpc"
	log "datapace/logger"
	"datapace/streams"
	"datapace/streams/api"
	grpcapi "datapace/streams/api/grpc"
	httpapi "datapace/streams/api/http"
	"datapace/streams/executions"
	"datapace/streams/mongo"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"datapace/streams/access"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	mgo "gopkg.in/mgo.v2"
)

const (
	envHTTPPort      = "DATAPACE_STREAMS_HTTP_PORT"
	envGRPCPort      = "DATAPACE_STREAMS_GRPC_PORT"
	envDBURL         = "DATAPACE_STREAMS_DB_URL"
	envDBName        = "DATAPACE_STREAMS_DB_NAME"
	envDBUser        = "DATAPACE_STREAMS_DB_USER"
	envDBPass        = "DATAPACE_STREAMS_DB_PASS"
	envAuthURL       = "DATAPACE_AUTH_URL"
	envAccessURL     = "DATAPACE_ACCESS_CONTROL_URL"
	envExecutionsURL = "DATAPACE_EXECUTIONS_URL"

	defHTTPPort      = "8080"
	defGRPCPort      = "8081"
	defDBURL         = "0.0.0.0"
	defDBName        = "streams"
	defDBUser        = ""
	defDBPass        = ""
	defAuthURL       = "localhost:8081"
	defAccessURL     = "localhost:8081"
	defExecutionsURL = "localhost:8081"

	dbConnectTimeout = 5000
	dbSocketTimeout  = 5000
)

type config struct {
	HTTPPort         string
	GRPCPort         string
	DBURL            string
	DBName           string
	DBUser           string
	DBPass           string
	DBConnectTimeout int
	DBSocketTimeout  int
	AuthURL          string
	AccessURL        string
	ExecutionsURL    string
}

func main() {
	cfg := loadConfig()
	logger := log.New(os.Stdout)
	ms := connectToDB(cfg, logger)
	defer ms.Close()

	authConn := connectToGRPCService(cfg.AuthURL, logger)
	accessConn := connectToGRPCService(cfg.AccessURL, logger)
	execConn := connectToGRPCService(cfg.ExecutionsURL, logger)

	svc, auth := newServices(ms, authConn, accessConn, execConn, logger)

	errs := make(chan error, 2)
	go startHTTPServer(svc, auth, cfg.HTTPPort, logger, errs)

	go startGRPCServer(svc, cfg.GRPCPort, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err := <-errs
	logger.Error(fmt.Sprintf("Streams service terminated: %s", err))
}

func loadConfig() config {
	return config{
		HTTPPort:         datapace.Env(envHTTPPort, defHTTPPort),
		GRPCPort:         datapace.Env(envGRPCPort, defGRPCPort),
		DBURL:            datapace.Env(envDBURL, defDBURL),
		DBName:           datapace.Env(envDBName, defDBName),
		DBUser:           datapace.Env(envDBUser, defDBUser),
		DBPass:           datapace.Env(envDBPass, defDBPass),
		DBConnectTimeout: dbConnectTimeout,
		DBSocketTimeout:  dbSocketTimeout,
		AuthURL:          datapace.Env(envAuthURL, defAuthURL),
		AccessURL:        datapace.Env(envAccessURL, defAccessURL),
		ExecutionsURL:    datapace.Env(envExecutionsURL, defExecutionsURL),
	}
}

func connectToDB(cfg config, logger log.Logger) *mgo.Session {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{cfg.DBURL},
		Timeout:  time.Duration(cfg.DBSocketTimeout) * time.Millisecond,
		Database: cfg.DBName,
		Username: cfg.DBUser,
		Password: cfg.DBPass,
	}

	ms, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to Mongo: %s", err))
		os.Exit(1)
	}

	ms.SetSocketTimeout(time.Duration(cfg.DBSocketTimeout) * time.Millisecond)
	ms.SetMode(mgo.Monotonic, true)

	return ms
}

func connectToGRPCService(addr string, logger log.Logger) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to gRPC service on address %s: %s", addr, err))
		os.Exit(1)
	}
	return conn
}

func newServices(ms *mgo.Session, authConn *grpc.ClientConn, accessConn *grpc.ClientConn, execConn *grpc.ClientConn, logger log.Logger) (streams.Service, streams.Authorization) {
	repo := mongo.New(ms)
	acc := accessapi.NewClient(accessConn)
	accessControl := access.New(acc)

	ec := executionsapi.NewClient(execConn)
	ai := executions.New(ec)

	svc := streams.NewService(repo, accessControl, ai)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "streams",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "streams",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)
	auc := authapi.NewClient(authConn)
	auth := streams.NewAuthorization(auc, logger)

	return svc, auth
}

func startHTTPServer(svc streams.Service, auth streams.Authorization, port string, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("Streams service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, httpapi.MakeHandler(svc, auth))
}

func startGRPCServer(svc streams.Service, port string, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", p)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to listen on port %s: %s", port, err))
	}

	server := grpc.NewServer()
	datapace.RegisterStreamsServiceServer(server, grpcapi.NewServer(svc))
	logger.Info(fmt.Sprintf("Streams gRPC service started, exposed port %s", port))
	errs <- server.Serve(listener)
}
