package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/datapace/datapace"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	fabricConfig "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	mgo "gopkg.in/mgo.v2"

	access "github.com/datapace/datapace/access-control"
	"github.com/datapace/datapace/access-control/api"
	grpcapi "github.com/datapace/datapace/access-control/api/grpc"
	httpapi "github.com/datapace/datapace/access-control/api/http"
	"github.com/datapace/datapace/access-control/auth"
	"github.com/datapace/datapace/access-control/fabric"
	"github.com/datapace/datapace/access-control/mongo"
	authapi "github.com/datapace/datapace/auth/api/grpc"
	log "github.com/datapace/datapace/logger"
	accessproto "github.com/datapace/datapace/proto/access"
)

const (
	envHTTPPort       = "DATAPACE_ACCESS_CONTROL_HTTP_PORT"
	envGRPCPort       = "DATAPACE_ACCESS_CONTROL_GRPC_PORT"
	envDBURL          = "DATAPACE_ACCESS_CONTROL_DB_URL"
	envDBUser         = "DATAPACE_ACCESS_CONTROL_DB_USER"
	envDBPass         = "DATAPACE_ACCESS_CONTROL_DB_PASS"
	envDBName         = "DATAPACE_ACCESS_CONTROL_DB_NAME"
	envDbMigration    = "DATAPACE_ACCESS_CONTROL_DB_MIGRATION"
	envAuthURL        = "DATAPACE_AUTH_URL"
	envFabricOrgAdmin = "DATAPACE_ACCESS_CONTROL_FABRIC_ADMIN"
	envFabricOrgName  = "DATAPACE_ACCESS_CONTROL_FABRIC_NAME"
	envDatapaceConfig = "DATAPACE_CONFIG"
	envChaincodeID    = "DATAPACE_ACCESS_CONTROL_CHAINCODE"

	defHTTPPort       = "8080"
	defGRPCPort       = "8081"
	defDBURL          = "0.0.0.0"
	defDBUser         = ""
	defDBPass         = ""
	defDBName         = "access"
	defDbMigration    = "false"
	defAuthURL        = "localhost:8081"
	defFabricOrgAdmin = "admin"
	defFabricOrgName  = "org1"
	defDatapaceConfig = "/src/github.com/datapace/datapace/config"
	defChaincodeID    = "access"

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
	dbMigration      bool
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

	conn := newGRPCConn(cfg.authURL, logger)
	defer conn.Close()

	sdk := newFabricSDK(cfg.fabricConfigFile, logger)
	defer sdk.Close()

	svc := newService(cfg, sdk, ms, conn, logger)

	errs := make(chan error, 2)

	go startHTTPServer(svc, cfg.httpPort, logger, errs)

	go startGRPCServer(svc, cfg.grpcPort, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err := <-errs
	logger.Error(fmt.Sprintf("Access Control service terminated: %s", err))

}

func loadConfig() config {
	configDir := datapace.Env(envDatapaceConfig, defDatapaceConfig)
	configFile := fmt.Sprintf("%s/%s", configDir, fabricConfigFile)

	dbMigrationRaw := datapace.Env(envDbMigration, defDbMigration)
	dbMigration, err := strconv.ParseBool(dbMigrationRaw)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse the %s value: %s, should be \"true\"/\"false\"", envDbMigration, dbMigrationRaw))
	}

	return config{
		httpPort:         datapace.Env(envHTTPPort, defHTTPPort),
		grpcPort:         datapace.Env(envGRPCPort, defGRPCPort),
		dbURL:            datapace.Env(envDBURL, defDBURL),
		dbUser:           datapace.Env(envDBUser, defDBUser),
		dbPass:           datapace.Env(envDBPass, defDBPass),
		dbName:           datapace.Env(envDBName, defDBName),
		dbConnectTimeout: dbConnectTimeout,
		dbSocketTimeout:  dbSocketTimeout,
		dbMigration:      dbMigration,
		authURL:          datapace.Env(envAuthURL, defAuthURL),
		fabricOrgName:    datapace.Env(envFabricOrgName, defFabricOrgName),
		fabricOrgAdmin:   datapace.Env(envFabricOrgAdmin, defFabricOrgAdmin),
		fabricConfigFile: configFile,
		chaincodeID:      datapace.Env(envChaincodeID, defChaincodeID),
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

func newGRPCConn(url string, logger log.Logger) *grpc.ClientConn {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to grpc service on address %s: %s", url, err))
		os.Exit(1)
	}

	return conn
}

func newFabricSDK(configFile string, logger log.Logger) *fabsdk.FabricSDK {
	sdk, err := fabsdk.New(fabricConfig.FromFile(configFile))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize fabric SDK: %s", err))
		os.Exit(1)
	}

	return sdk
}

func newService(cfg config, sdk *fabsdk.FabricSDK, ms *mgo.Session, conn *grpc.ClientConn, logger log.Logger) access.Service {
	repo := mongo.NewAccessRequestRepository(ms)
	if cfg.dbMigration {
		logger.Info("Starting the data migration...")
		err := repo.MigrateSameSenderAndReceiver()
		if err == nil {
			logger.Info("Data migration finished successfully")
		} else {
			panic(fmt.Sprintf("Data migration failed: %s", err))
		}
	}
	if err := mongo.EnsureIndices(ms); err != nil {
		panic(err)
	}
	al := fabric.NewRequestLedger(
		sdk,
		cfg.fabricOrgAdmin,
		cfg.fabricOrgName,
		cfg.chaincodeID,
		logger,
	)
	asc := authapi.NewClient(conn)
	as := auth.New(asc)
	svc := access.New(as, repo, al)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "access_control",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "access_control",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)
	return svc
}

func startHTTPServer(svc access.Service, port string, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("Access Control HTTP service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, httpapi.MakeHandler(svc))
}

func startGRPCServer(svc access.Service, port string, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", p)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to listen on port %s: %s", port, err))
	}

	server := grpc.NewServer()
	accessproto.RegisterAccessServiceServer(server, grpcapi.NewServer(svc))
	logger.Info(fmt.Sprintf("Access Control gRPC service started, exposed port %s", port))
	errs <- server.Serve(listener)
}
