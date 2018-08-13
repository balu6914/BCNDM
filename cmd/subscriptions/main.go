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
	streamsapi "monetasa/streams/api/grpc"
	"monetasa/subscriptions"
	"monetasa/subscriptions/api"
	"monetasa/subscriptions/mongo"
	"monetasa/subscriptions/streams"
	"monetasa/subscriptions/transactions"
	transactionsapi "monetasa/transactions/api/grpc"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	mgo "gopkg.in/mgo.v2"

	"google.golang.org/grpc"
)

const (
	envSubsPort        = "MONETASA_SUBSCRIPTIONS_PORT"
	envAuthURL         = "MONETASA_AUTH_URL"
	envTransactionsURL = "MONETASA_TRANSACTIONS_URL"
	envStreamsURL      = "MONETASA_STREAMS_URL"
	envMongoURL        = "MONETASA_SUBSCRIPTIONS_DB_URL"
	envMongoUser       = "MONETASA_SUBSCRIPTIONS_DB_USER"
	envMongoPass       = "MONETASA_SUBSCRIPTIONS_DB_PASS"
	envMongoDatabase   = "MONETASA_SUBSCRIPTIONS_DB_NAME"

	defSubsPort        = "8080"
	defAuthURL         = "localhost:8081"
	defTransactionsURL = "localhost:8081"
	defStreamsURL      = "localhost:8081"
	defSubsURL         = "0.0.0.0"
	defMongoURL        = "0.0.0.0"
	defMongoUser       = ""
	defMongoPass       = ""
	defMongoDatabase   = "subscriptions"

	defMongoConnectTimeout = 5000
	defMongoSocketTimeout  = 5000
)

type config struct {
	SubsPort            string
	AuthURL             string
	TransactionsURL     string
	StreamsURL          string
	MongoURL            string
	MongoUser           string
	MongoPass           string
	MongoDatabase       string
	MongoConnectTimeout int
	MongoSocketTimeout  int
}

func main() {
	cfg := loadConfig()

	logger := log.New(os.Stdout)

	ms := connectToDB(cfg, logger)
	defer ms.Close()

	authConn := newGRPCConn(cfg.AuthURL, logger)
	defer authConn.Close()
	ac := authapi.NewClient(authConn)

	transactionsConn := newGRPCConn(cfg.TransactionsURL, logger)
	defer transactionsConn.Close()
	tc := transactionsapi.NewClient(transactionsConn)

	streamsConn := newGRPCConn(cfg.StreamsURL, logger)
	defer streamsConn.Close()
	sc := streamsapi.NewClient(streamsConn)

	svc := newService(ms, sc, tc, logger)

	errs := make(chan error, 2)

	go serveHTTPServer(svc, ac, cfg.SubsPort, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err := <-errs
	logger.Error(fmt.Sprintf("Subscriptions service terminated: %s", err))
}

func loadConfig() config {
	return config{
		SubsPort:            monetasa.Env(envSubsPort, defSubsPort),
		StreamsURL:          monetasa.Env(envStreamsURL, defStreamsURL),
		TransactionsURL:     monetasa.Env(envTransactionsURL, defTransactionsURL),
		MongoURL:            monetasa.Env(envMongoURL, defMongoURL),
		MongoUser:           monetasa.Env(envMongoUser, defMongoUser),
		MongoPass:           monetasa.Env(envMongoPass, defMongoPass),
		MongoDatabase:       monetasa.Env(envMongoDatabase, defMongoDatabase),
		MongoConnectTimeout: defMongoConnectTimeout,
		MongoSocketTimeout:  defMongoSocketTimeout,
		AuthURL:             monetasa.Env(envAuthURL, defAuthURL),
	}
}

func connectToDB(cfg config, logger log.Logger) *mgo.Session {
	conn, err := mongo.Connect(
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

	return conn
}

func newGRPCConn(url string, logger log.Logger) *grpc.ClientConn {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to grpc service: %s", err))
		os.Exit(1)
	}

	return conn
}

func newService(ms *mgo.Session, sc monetasa.StreamsServiceClient, tc monetasa.TransactionsServiceClient, logger log.Logger) subscriptions.Service {
	ss := streams.NewService(sc)
	ts := transactions.NewService(tc)

	repo := mongo.NewRepository(ms)
	svc := subscriptions.New(repo, ss, ts)

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

	return svc
}

func serveHTTPServer(svc subscriptions.Service, ac monetasa.AuthServiceClient, port string, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("Subscriptions service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, api.MakeHandler(svc, ac))
}
