package main

import (
	"fmt"
	accessProtoV2 "github.com/datapace/datapace/proto/accessv2"
	"github.com/datapace/datapace/subscriptions/accessv2"
	"github.com/datapace/datapace/subscriptions/pub"
	"github.com/datapace/datapace/subscriptions/sharing"
	"github.com/datapace/events/pubsub"
	sharingApi "github.com/datapace/sharing/api/grpc"
	sharingProto "github.com/datapace/sharing/proto"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/datapace/datapace"

	authapi "github.com/datapace/datapace/auth/api/grpc"
	log "github.com/datapace/datapace/logger"
	authproto "github.com/datapace/datapace/proto/auth"
	streamsproto "github.com/datapace/datapace/proto/streams"
	transactionsproto "github.com/datapace/datapace/proto/transactions"
	streamsapi "github.com/datapace/datapace/streams/api/grpc"
	"github.com/datapace/datapace/subscriptions"
	"github.com/datapace/datapace/subscriptions/api"
	"github.com/datapace/datapace/subscriptions/mongo"
	"github.com/datapace/datapace/subscriptions/proxy"
	"github.com/datapace/datapace/subscriptions/streams"
	"github.com/datapace/datapace/subscriptions/transactions"
	transactionsapi "github.com/datapace/datapace/transactions/api/grpc"

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
	envSharingUrl      = "DATAPACE_SHARING_URL"
	envMsgBusUrl       = "DATAPACE_MSG_BUS_URL"
	envSubjFmtCreate   = "DATAPACE_SUBSCRIPTIONS_SUBJ_FMT_CREATE"
	envAccessV2Url     = "DATAPACE_ACCESS_V2_URL"

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
	defSharingUrl      = "localhost:8081"
	defMsgBusUrl       = "nats://localhost:4222"
	defSubjFmtCreate   = "subscriptions.stream.owner.%s"
	defProxyURL        = "http://localhost:8080"
	defSubsURL         = "0.0.0.0"
	defMongoURL        = "0.0.0.0"
	defMongoUser       = ""
	defMongoPass       = ""
	defMongoDatabase   = "subscriptions"
	defAccessV2Url     = "localhost:8081"

	defMongoConnectTimeout = 5000
	defMongoSocketTimeout  = 5000
)

type config struct {
	SubsPort            string
	AuthURL             string
	TransactionsURL     string
	StreamsURL          string
	SharingUrl          string
	ProxyURL            string
	MongoURL            string
	MongoUser           string
	MongoPass           string
	MongoDatabase       string
	MongoConnectTimeout int
	MongoSocketTimeout  int
	MsgBusUrl           string
	SubjFmt             pub.SubjectFormat
	AccessV2Url         string
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

	sharingConn := newOptionalGrpcConn(cfg.SharingUrl, logger)
	var sharingClient sharingProto.SharingServiceClient = nil
	if sharingConn != nil {
		defer sharingConn.Close()
		sharingClient = sharingApi.NewClient(sharingConn)
	}

	accessV2Conn := newOptionalGrpcConn(cfg.AccessV2Url, logger)
	var accessV2Client accessProtoV2.ServiceClient = nil
	if accessV2Conn != nil {
		defer accessV2Conn.Close()
		accessV2Client = accessProtoV2.NewServiceClient(accessV2Conn)
	}

	svc := newService(ac, ms, sc, tc, sharingClient, accessV2Client, cfg.ProxyURL, logger)

	pubSubSvc, msgBusConnErr := pubsub.NewService(cfg.MsgBusUrl)
	if msgBusConnErr != nil {
		logger.Warn(fmt.Sprintf("Failed to connect message bus @ %s, events publishing won't be available", cfg.MsgBusUrl))
	}
	pubSvc := pub.NewService(pubSubSvc, cfg.SubjFmt)
	pubSvc = pub.NewLoggingMiddleware(pubSvc, logger)
	svc = subscriptions.NewPubMiddleware(svc, pubSvc)

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
		SharingUrl:          datapace.Env(envSharingUrl, defSharingUrl),
		ProxyURL:            datapace.Env(envProxyURL, defProxyURL),
		MongoURL:            datapace.Env(envMongoURL, defMongoURL),
		MongoUser:           datapace.Env(envMongoUser, defMongoUser),
		MongoPass:           datapace.Env(envMongoPass, defMongoPass),
		MongoDatabase:       datapace.Env(envMongoDatabase, defMongoDatabase),
		MongoConnectTimeout: defMongoConnectTimeout,
		MongoSocketTimeout:  defMongoSocketTimeout,
		AuthURL:             datapace.Env(envAuthURL, defAuthURL),
		MsgBusUrl:           datapace.Env(envMsgBusUrl, defMsgBusUrl),
		SubjFmt: pub.SubjectFormat{
			SubscriptionCreate: datapace.Env(envSubjFmtCreate, defSubjFmtCreate),
		},
		AccessV2Url: datapace.Env(envAccessV2Url, defAccessV2Url),
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

// tries to start the optional connection (no fatal failure, if it can not connect, returns nil)
func newOptionalGrpcConn(url string, logger log.Logger) *grpc.ClientConn {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to connect to grpc service: %s", err))
		conn = nil
	}
	return conn
}

func newService(
	ac authproto.AuthServiceClient, ms *mgo.Session, sc streamsproto.StreamsServiceClient,
	tc transactionsproto.TransactionsServiceClient, sharingClient sharingProto.SharingServiceClient,
	accessV2Client accessProtoV2.ServiceClient, proxyURL string, logger log.Logger,
) subscriptions.Service {
	ss := streams.NewService(sc)
	ts := transactions.NewService(tc)
	var sharingSvc sharing.Service = nil
	if sharingClient != nil {
		sharingSvc = sharing.NewService(sharingClient)
		sharingSvc = sharing.NewLoggingMiddleware(sharingSvc, logger)
	}
	var accessV2Svc accessv2.Service = nil
	if accessV2Client != nil {
		accessV2Svc = accessv2.NewService(accessV2Client)
		accessV2Svc = accessv2.NewLoggingMiddleware(accessV2Svc, logger)
	}
	ps := proxy.New(proxyURL)

	repo := mongo.NewSubscriptionRepository(ms)
	svc := subscriptions.New(ac, repo, ss, ps, ts, sharingSvc, accessV2Svc)

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

func serveHTTPServer(svc subscriptions.Service, ac authproto.AuthServiceClient, port string, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("Subscriptions service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, api.MakeHandler(svc, ac))
}
