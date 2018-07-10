package main

import (
	"fmt"
	"monetasa"
	authapi "monetasa/auth/api/grpc"
	log "monetasa/logger"
	"monetasa/transactions"
	"monetasa/transactions/api"
	grpcapi "monetasa/transactions/api/grpc"
	httpapi "monetasa/transactions/api/http"
	"monetasa/transactions/fabric"
	"monetasa/transactions/mongo"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	mgo "gopkg.in/mgo.v2"
)

const (
	envHTTPPort          = "MONETASA_TRANSACTIONS_HTTP_PORT"
	envGRPCPort          = "MONETASA_TRANSACTIONS_GRPC_PORT"
	envMongoURL          = "MONETASA_TRANSACTIONS_MONGO_URL"
	envMongoUser         = "MONETASA_TRANSACTIONS_MONGO_USER"
	envMongoPass         = "MONETASA_TRANSACTIONS_MONGO_PASS"
	envMongoDatabase     = "MONETASA_TRANSACTIONS_MONGO_DB"
	envFabricOrgAdmin    = "MONETASA_TRANSACTIONS_FABRIC_ADMIN"
	envFabricOrgName     = "MONETASA_TRANSACTIONS_FABRIC_NAME"
	envFabricConfFile    = "MONETASA_TRANSACTIONS_FABRIC_CONF"
	envFabricChaincodeID = "MONETASA_TRANSACTIONS_FABRIC_CHAINCODE"
	envAuthURL           = "MONETASA_AUTH_URL"

	defHTTPPort          = "8080"
	defGRPCPort          = "8081"
	defMongoURL          = "0.0.0.0"
	defMongoUser         = ""
	defMongoPass         = ""
	defMongoDatabase     = "transactions"
	defFabricOrgAdmin    = "admin"
	defFabricOrgName     = "org1"
	defFabricConfFile    = "/src/monetasa/config/fabric/config.yaml"
	defFabricChaincodeID = "token"
	defAuthURL           = "localhost:8081"

	mongoConnectTimeout = 5000
	mongoSocketTimeout  = 5000
)

type conf struct {
	httpPort            string
	grpcPort            string
	mongoURL            string
	mongoUser           string
	mongoPass           string
	mongoDatabase       string
	mongoConnectTimeout int
	mongoSocketTimeout  int
	fabricOrgAdmin      string
	fabricOrgName       string
	fabricConfFile      string
	fabricChaincodeID   string
	authURL             string
}

func main() {
	cfg := loadConfig()

	logger := log.New(os.Stdout)

	sdk := newFabricSDK(cfg.fabricConfFile, logger)
	defer sdk.Close()

	ms := connectToDB(cfg, logger)
	defer ms.Close()

	conn := newGRPCConn(cfg.authURL, logger)
	defer conn.Close()

	ac := authapi.NewClient(conn)

	svc := newService(cfg, sdk, ms, logger)

	errs := make(chan error, 2)

	go startHTTPServer(svc, ac, cfg.httpPort, logger, errs)

	go startGRPCServer(svc, cfg.grpcPort, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err := <-errs
	logger.Error(fmt.Sprintf("Auth service terminated: %s", err))
}

func loadConfig() conf {
	confPath := monetasa.Env(envFabricConfFile, defFabricConfFile)
	fullConfPath := fmt.Sprintf("%s%s", os.Getenv("GOPATH"), confPath)
	return conf{
		httpPort:            monetasa.Env(envHTTPPort, defHTTPPort),
		grpcPort:            monetasa.Env(envGRPCPort, defGRPCPort),
		mongoURL:            monetasa.Env(envMongoURL, defMongoURL),
		mongoUser:           monetasa.Env(envMongoUser, defMongoUser),
		mongoPass:           monetasa.Env(envMongoPass, defMongoPass),
		mongoDatabase:       monetasa.Env(envMongoDatabase, defMongoDatabase),
		mongoConnectTimeout: mongoConnectTimeout,
		mongoSocketTimeout:  mongoSocketTimeout,
		fabricOrgAdmin:      monetasa.Env(envFabricOrgAdmin, defFabricOrgAdmin),
		fabricOrgName:       monetasa.Env(envFabricOrgName, defFabricOrgName),
		fabricConfFile:      fullConfPath,
		fabricChaincodeID:   monetasa.Env(envFabricChaincodeID, defFabricChaincodeID),
		authURL:             monetasa.Env(envAuthURL, defAuthURL),
	}
}

func newFabricSDK(confPath string, logger log.Logger) *fabsdk.FabricSDK {
	sdk, err := fabsdk.New(config.FromFile(confPath))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize fabric SDK: %s", err))
		os.Exit(1)
	}

	return sdk
}

func connectToDB(cfg conf, logger log.Logger) *mgo.Session {
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

func newService(cfg conf, sdk *fabsdk.FabricSDK, ms *mgo.Session, logger log.Logger) transactions.Service {
	bn := fabric.NewNetwork(sdk, cfg.fabricOrgAdmin, cfg.fabricOrgName, cfg.fabricChaincodeID, logger)
	users := mongo.NewUserRepository(ms)

	svc := transactions.New(users, bn)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "transactions",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "transactions",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)

	return svc
}

func newGRPCConn(authURL string, logger log.Logger) *grpc.ClientConn {
	conn, err := grpc.Dial(authURL, grpc.WithInsecure())
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to auth service: %s", err))
		os.Exit(1)
	}

	return conn
}

func startHTTPServer(svc transactions.Service, ac monetasa.AuthServiceClient, port string, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("Users HTTP service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, httpapi.MakeHandler(svc, ac))
}

func startGRPCServer(svc transactions.Service, port string, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", p)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to listen on port %s: %s", port, err))
	}

	server := grpc.NewServer()
	monetasa.RegisterTransactionsServiceServer(server, grpcapi.NewServer(svc))
	logger.Info(fmt.Sprintf("Transactions gRPC service started, exposed port %s", port))
	errs <- server.Serve(listener)
}
