package main

import (
	"fmt"
	"net"
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
	"github.com/datapace/datapace/transactions"
	"github.com/datapace/datapace/transactions/api"
	grpcapi "github.com/datapace/datapace/transactions/api/grpc"
	httpapi "github.com/datapace/datapace/transactions/api/http"
	"github.com/datapace/datapace/transactions/fabric"
	"github.com/datapace/datapace/transactions/mongo"
	"github.com/datapace/datapace/transactions/streams"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	fabricConfig "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	mgo "gopkg.in/mgo.v2"
)

const (
	envHTTPPort             = "DATAPACE_TRANSACTIONS_HTTP_PORT"
	envGRPCPort             = "DATAPACE_TRANSACTIONS_GRPC_PORT"
	envDBURL                = "DATAPACE_TRANSACTIONS_DB_URL"
	envDBUser               = "DATAPACE_TRANSACTIONS_DB_USER"
	envDBPass               = "DATAPACE_TRANSACTIONS_DB_PASS"
	envDBName               = "DATAPACE_TRANSACTIONS_DB_NAME"
	envFabricOrgAdmin       = "DATAPACE_TRANSACTIONS_FABRIC_ADMIN"
	envFabricOrgName        = "DATAPACE_TRANSACTIONS_FABRIC_NAME"
	envDatapaceConfig       = "DATAPACE_CONFIG"
	envTokenChaincodeID     = "DATAPACE_TRANSACTIONS_TOKEN_CHAINCODE"
	envContractsChaincodeID = "DATAPACE_TRANSACTIONS_CONTRACTS_CHAINCODE"
	envAuthURL              = "DATAPACE_AUTH_URL"
	envStreamsURL           = "DATAPACE_STREAMS_URL"

	defHTTPPort             = "8080"
	defGRPCPort             = "8081"
	defDBURL                = "0.0.0.0"
	defDBUser               = ""
	defDBPass               = ""
	defDBName               = "transactions"
	defFabricOrgAdmin       = "admin"
	defFabricOrgName        = "org1"
	defDatapaceConfig       = "/src/github.com/datapace/datapace/config"
	defTokenChaincodeID     = "token"
	defContractsChaincodeID = "contracts"
	defAuthURL              = "localhost:8081"
	defStreamsURL           = "localhost:8081"

	fabricConfigFile = "fabric/config.yaml"
	dbConnectTimeout = 5000
	dbSocketTimeout  = 5000
)

type config struct {
	httpPort             string
	grpcPort             string
	dbURL                string
	dbUser               string
	dbPass               string
	dbName               string
	dbConnectTimeout     int
	dbSocketTimeout      int
	fabricOrgAdmin       string
	fabricOrgName        string
	fabricConfigFile     string
	tokenChaincodeID     string
	feeChaincodeID       string
	contractsChaincodeID string
	authURL              string
	streamsURL           string
}

func main() {
	cfg := loadConfig()

	logger := log.New(os.Stdout)

	sdk := newFabricSDK(cfg.fabricConfigFile, logger)
	defer sdk.Close()

	ms := connectToDB(cfg, logger)
	defer ms.Close()

	authConn := newGRPCConn(cfg.authURL, logger)
	defer authConn.Close()

	ac := authapi.NewClient(authConn)

	streamsConn := newGRPCConn(cfg.streamsURL, logger)
	defer streamsConn.Close()
	sc := streamsapi.NewClient(streamsConn)

	svc := newService(cfg, sdk, ms, sc, logger)

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

func loadConfig() config {
	configDir := datapace.Env(envDatapaceConfig, defDatapaceConfig)
	configFile := fmt.Sprintf("%s/%s", configDir, fabricConfigFile)

	return config{
		httpPort:             datapace.Env(envHTTPPort, defHTTPPort),
		grpcPort:             datapace.Env(envGRPCPort, defGRPCPort),
		dbURL:                datapace.Env(envDBURL, defDBURL),
		dbUser:               datapace.Env(envDBUser, defDBUser),
		dbPass:               datapace.Env(envDBPass, defDBPass),
		dbName:               datapace.Env(envDBName, defDBName),
		dbConnectTimeout:     dbConnectTimeout,
		dbSocketTimeout:      dbSocketTimeout,
		fabricOrgAdmin:       datapace.Env(envFabricOrgAdmin, defFabricOrgAdmin),
		fabricOrgName:        datapace.Env(envFabricOrgName, defFabricOrgName),
		fabricConfigFile:     configFile,
		tokenChaincodeID:     datapace.Env(envTokenChaincodeID, defTokenChaincodeID),
		contractsChaincodeID: datapace.Env(envContractsChaincodeID, defContractsChaincodeID),
		authURL:              datapace.Env(envAuthURL, defAuthURL),
		streamsURL:           datapace.Env(envStreamsURL, defStreamsURL),
	}
}

func newFabricSDK(configFile string, logger log.Logger) *fabsdk.FabricSDK {
	sdk, err := fabsdk.New(fabricConfig.FromFile(configFile))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize fabric SDK: %s", err))
		os.Exit(1)
	}

	return sdk
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

func newService(cfg config, sdk *fabsdk.FabricSDK, ms *mgo.Session, streamsClient streamsproto.StreamsServiceClient, logger log.Logger) transactions.Service {
	tl := fabric.NewTokenLedger(
		sdk,
		cfg.fabricOrgAdmin,
		cfg.fabricOrgName,
		cfg.tokenChaincodeID,
		cfg.contractsChaincodeID,
		logger,
	)
	cl := fabric.NewContractLedger(
		sdk,
		cfg.fabricOrgAdmin,
		cfg.fabricOrgName,
		cfg.contractsChaincodeID,
		logger,
	)
	users := mongo.NewUserRepository(ms)
	contracts := mongo.NewContractRepository(ms)
	sc := streams.NewService(streamsClient)

	svc := transactions.New(users, tl, cl, contracts, sc)
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

func startHTTPServer(svc transactions.Service, ac authproto.AuthServiceClient, port string, logger log.Logger, errs chan error) {
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
	transactionsproto.RegisterTransactionsServiceServer(server, grpcapi.NewServer(svc))
	logger.Info(fmt.Sprintf("Transactions gRPC service started, exposed port %s", port))
	errs <- server.Serve(listener)
}
