package main

import (
	"fmt"
	"monetasa"
	authapi "monetasa/auth/api/grpc"
	log "monetasa/logger"
	"monetasa/streams"
	"monetasa/streams/api"
	grpcapi "monetasa/streams/api/grpc"
	httpapi "monetasa/streams/api/http"
	"monetasa/streams/mongo"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	mgo "gopkg.in/mgo.v2"
)

const (
	envHTTPPort = "MONETASA_STREAMS_HTTP_PORT"
	envGRPCPort = "MONETASA_STREAMS_GRPC_PORT"
	envDBURL    = "MONETASA_STREAMS_DB_URL"
	envDBName   = "MONETASA_STREAMS_DB_NAME"
	envDBUser   = "MONETASA_STREAMS_DB_USER"
	envDBPass   = "MONETASA_STREAMS_DB_PASS"
	envAuthURL  = "MONETASA_AUTH_URL"

	defHTTPPort = "8080"
	defGRPCPort = "8081"
	defDBURL    = "0.0.0.0"
	defDBName   = "streams"
	defDBUser   = ""
	defDBPass   = ""
	defAuthURL  = "localhost:8081"

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
}

func main() {
	cfg := loadConfig()
	logger := log.New(os.Stdout)
	ms := connectToDB(cfg, logger)
	defer ms.Close()

	conn := connectToAuthService(cfg.AuthURL, logger)
	svc, auth := newServices(ms, conn, logger)

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
		HTTPPort:         monetasa.Env(envHTTPPort, defHTTPPort),
		GRPCPort:         monetasa.Env(envGRPCPort, defGRPCPort),
		DBURL:            monetasa.Env(envDBURL, defDBURL),
		DBName:           monetasa.Env(envDBName, defDBName),
		DBUser:           monetasa.Env(envDBUser, defDBUser),
		DBPass:           monetasa.Env(envDBPass, defDBPass),
		DBConnectTimeout: dbConnectTimeout,
		DBSocketTimeout:  dbSocketTimeout,
		AuthURL:          monetasa.Env(envAuthURL, defAuthURL),
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

func connectToAuthService(authAddr string, logger log.Logger) *grpc.ClientConn {
	conn, err := grpc.Dial(authAddr, grpc.WithInsecure())
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to auth service: %s", err))
		os.Exit(1)
	}
	return conn
}

func newServices(ms *mgo.Session, conn *grpc.ClientConn, logger log.Logger) (streams.Service, streams.Authorization) {
	repo := mongo.New(ms)
	svc := streams.NewService(repo)
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
	a := authapi.NewClient(conn)
	auth := streams.NewAuthorization(a, logger)

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
	monetasa.RegisterStreamsServiceServer(server, grpcapi.NewServer(svc))
	logger.Info(fmt.Sprintf("Streams gRPC service started, exposed port %s", port))
	errs <- server.Serve(listener)
}
