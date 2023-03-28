package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/datapace/datapace"
	authapi "github.com/datapace/datapace/auth/api/grpc"
	log "github.com/datapace/datapace/logger"
	authproto "github.com/datapace/datapace/proto/auth"
	termsproto "github.com/datapace/datapace/proto/terms"
	"github.com/datapace/datapace/terms"
	"github.com/datapace/datapace/terms/api"
	grpcapi "github.com/datapace/datapace/terms/api/grpc"
	httpapi "github.com/datapace/datapace/terms/api/http"
	"github.com/datapace/datapace/terms/fabric"
	"github.com/datapace/datapace/terms/mongo"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	fabricConfig "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"gopkg.in/mgo.v2"
)

const (
	envHTTPPort       = "DATAPACE_TERMS_HTTP_PORT"
	envGRPCPort       = "DATAPACE_TERMS_GRPC_PORT"
	envDBURL          = "DATAPACE_TERMS_DB_URL"
	envDBUser         = "DATAPACE_TERMS_DB_USER"
	envDBPass         = "DATAPACE_TERMS_DB_PASS"
	envDBName         = "DATAPACE_TERMS_DB_NAME"
	envAuthURL        = "DATAPACE_AUTH_URL"
	envFabricOrgAdmin = "DATAPACE_TERMS_FABRIC_ADMIN"
	envFabricOrgName  = "DATAPACE_TERMS_FABRIC_NAME"
	envDatapaceConfig = "DATAPACE_CONFIG"
	envChaincodeID    = "DATAPACE_TERMS_CHAINCODE"

	defHTTPPort       = "8080"
	defGRPCPort       = "8081"
	defDBURL          = "0.0.0.0"
	defDBUser         = ""
	defDBPass         = ""
	defDBName         = "terms"
	defAuthURL        = "localhost:8081"
	defFabricOrgAdmin = "admin"
	defFabricOrgName  = "org1"
	defDatapaceConfig = "/src/github.com/datapace/datapace/config"
	defChaincodeID    = "terms"

	fabricConfigFile = "fabric/config.yaml"
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
	fabricOrgAdmin   string
	fabricOrgName    string
	fabricConfigFile string
	chaincodeID      string
}

func main() {
	cfg := loadConfig()

	logger := log.New(os.Stdout)

	ms := connectToDB(cfg, logger)
	defer ms.Close()
	aconn := newGRPCConn(cfg.authURL, logger)
	defer aconn.Close()

	auth := authapi.NewClient(aconn)
	sdk := newFabricSDK(cfg.fabricConfigFile, logger)
	defer sdk.Close()
	svc := newService(cfg, ms, sdk, logger)
	errs := make(chan error, 2)
	go startGRPCServer(svc, cfg.grpcPort, logger, errs)
	go startHTTPServer(svc, auth, cfg.httpPort, logger, errs)
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
	configDir := datapace.Env(envDatapaceConfig, defDatapaceConfig)
	configFile := fmt.Sprintf("%s/%s", configDir, fabricConfigFile)
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
		fabricOrgName:    datapace.Env(envFabricOrgName, defFabricOrgName),
		fabricOrgAdmin:   datapace.Env(envFabricOrgAdmin, defFabricOrgAdmin),
		fabricConfigFile: configFile,
		chaincodeID:      datapace.Env(envChaincodeID, defChaincodeID),
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

func newService(cfg config, ms *mgo.Session, sdk *fabsdk.FabricSDK, logger log.Logger) terms.Service {
	tl := fabric.NewTermsLedger(
		sdk,
		cfg.fabricOrgAdmin,
		cfg.fabricOrgName,
		cfg.chaincodeID,
		logger,
	)
	repo := mongo.NewTermsRepository(ms)

	svc := terms.New(repo, tl)

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

func startHTTPServer(svc terms.Service, auth authproto.AuthServiceClient, port string, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("Terms service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, httpapi.MakeHandler(svc, auth))
}

func newFabricSDK(configFile string, logger log.Logger) *fabsdk.FabricSDK {
	sdk, err := fabsdk.New(fabricConfig.FromFile(configFile))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize fabric SDK: %s", err))
		os.Exit(1)
	}

	return sdk
}
