package main

import (
	"fmt"
	"datapace"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	authapi "datapace/auth/api/grpc"
	log "datapace/logger"
	streamsapi "datapace/streams/api/grpc"
	"datapace/subscriptions"
	"datapace/subscriptions/api"
	"datapace/subscriptions/mongo"
	"datapace/subscriptions/proxy"
	"datapace/subscriptions/streams"
	"datapace/subscriptions/transactions"
	transactionsapi "datapace/transactions/api/grpc"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	mgo "gopkg.in/mgo.v2"

	"google.golang.org/grpc"
)

const (
	envSubsPort        = "DATAPACE_SUBSCRIPTIONS_PORT"
	envAuthURL         = "DATAPACE_AUTH_URL"
	envTransactionsURL = "DATAPACE_TRANSACTIONS_URL"
	envStreamsURL      = "DATAPACE_STREAMS_URL"

	// HTTP prefixed, because all others are gRPC.
	envProxyURL      = "DATAPACE_PROXY_URL"
	envMongoURL      = "DATAPACE_SUBSCRIPTIONS_DB_URL"
	envMongoUser     = "DATAPACE_SUBSCRIPTIONS_DB_USER"
	envMongoPass     = "DATAPACE_SUBSCRIPTIONS_DB_PASS"
	envMongoDatabase = "DATAPACE_SUBSCRIPTIONS_DB_NAME"

	defSubsPort        = "8080"
	defAuthURL         = "localhost:8081"
	defTransactionsURL = "localhost:8081"
	defStreamsURL      = "localhost:8081"
	defProxyURL        = "http://localhost:8080"
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
	ProxyURL            string
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

	svc := newService(ac, ms, sc, tc, cfg.ProxyURL, logger)

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
		SubsPort:            datapace.Env(envSubsPort, defSubsPort),
		StreamsURL:          datapace.Env(envStreamsURL, defStreamsURL),
		TransactionsURL:     datapace.Env(envTransactionsURL, defTransactionsURL),
		ProxyURL:            datapace.Env(envProxyURL, defProxyURL),
		MongoURL:            datapace.Env(envMongoURL, defMongoURL),
		MongoUser:           datapace.Env(envMongoUser, defMongoUser),
		MongoPass:           datapace.Env(envMongoPass, defMongoPass),
		MongoDatabase:       datapace.Env(envMongoDatabase, defMongoDatabase),
		MongoConnectTimeout: defMongoConnectTimeout,
		MongoSocketTimeout:  defMongoSocketTimeout,
		AuthURL:             datapace.Env(envAuthURL, defAuthURL),
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

func newService(ac datapace.AuthServiceClient, ms *mgo.Session, sc datapace.StreamsServiceClient, tc datapace.TransactionsServiceClient, proxyURL string, logger log.Logger) subscriptions.Service {
	ss := streams.NewService(sc)
	ts := transactions.NewService(tc)
	ps := proxy.New(proxyURL)

	repo := mongo.NewSubscriptionRepository(ms)
	svc := subscriptions.New(ac, repo, ss, ps, ts)

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

func serveHTTPServer(svc subscriptions.Service, ac datapace.AuthServiceClient, port string, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("Subscriptions service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, api.MakeHandler(svc, ac))
}
