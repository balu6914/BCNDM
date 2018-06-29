package main

import (
	"fmt"
	"monetasa"
	"monetasa/auth"
	"monetasa/auth/api"
	grpcapi "monetasa/auth/api/grpc"
	httpapi "monetasa/auth/api/http"
	"monetasa/auth/bcrypt"
	"monetasa/auth/fabric"
	"monetasa/auth/jwt"
	"monetasa/auth/mongo"
	log "monetasa/logger"
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
	envHTTPPort          = "MONETASA_AUTH_HTTP_PORT"
	envGRPCPort          = "MONETASA_AUTH_GRPC_PORT"
	envMongoURL          = "MONETASA_AUTH_MONGO_URL"
	envMongoUser         = "MONETASA_AUTH_MONGO_USER"
	envMongoPass         = "MONETASA_AUTH_MONGO_PASS"
	envMongoDatabase     = "MONETASA_AUTH_MONGO_DB"
	envSecret            = "MONETASA_AUTH_SECRET"
	envFabricOrgAdmin    = "MONETASA_AUTH_FABRIC_ADMIN"
	envFabricOrgName     = "MONETASA_AUTH_FABRIC_NAME"
	envFabricConfFile    = "MONETASA_AUTH_FABRIC_CONF"
	envFabricChannelID   = "MONETASA_AUTH_FABRIC_CHANNEL"
	envFabricChaincodeID = "MONETASA_AUTH_FABRIC_CHAINCODE"

	defHTTPPort          = "8080"
	defGRPCPort          = "8081"
	defMongoURL          = "0.0.0.0"
	defMongoUser         = ""
	defMongoPass         = ""
	defMongoDatabase     = "auth"
	defSecret            = "monetasa"
	defFabricOrgAdmin    = "admin"
	defFabricOrgName     = "org1"
	defFabricConfFile    = "/src/monetasa/config/fabric/config.yaml"
	defFabricChannelID   = "myc"
	defFabricChaincodeID = "token"

	mongoConnectTimeout = 5000
	mongoSocketTimeout  = 5000
)

type config struct {
	HTTPPort            string
	GRPCPort            string
	PostgresURL         string
	MongoURL            string
	MongoUser           string
	MongoPass           string
	MongoDatabase       string
	MongoConnectTimeout int
	MongoSocketTimeout  int
	Secret              string
	FabricOrgAdmin      string
	FabricOrgName       string
	FabricConfFile      string
	FabricChannelID     string
	FabricChaincodeID   string
}

func main() {
	fullConfPath := fmt.Sprintf("%s%s", os.Getenv("GOPATH"), defFabricConfFile)
	cfg := config{
		HTTPPort:            monetasa.Env(envHTTPPort, defHTTPPort),
		GRPCPort:            monetasa.Env(envGRPCPort, defGRPCPort),
		MongoURL:            monetasa.Env(envMongoURL, defMongoURL),
		MongoUser:           monetasa.Env(envMongoUser, defMongoUser),
		MongoPass:           monetasa.Env(envMongoPass, defMongoPass),
		MongoDatabase:       monetasa.Env(envMongoDatabase, defMongoDatabase),
		MongoConnectTimeout: mongoConnectTimeout,
		MongoSocketTimeout:  mongoSocketTimeout,
		Secret:              monetasa.Env(envSecret, defSecret),
		FabricOrgAdmin:      monetasa.Env(envFabricOrgAdmin, defFabricOrgAdmin),
		FabricOrgName:       monetasa.Env(envFabricOrgName, defFabricOrgName),
		FabricConfFile:      monetasa.Env(envFabricConfFile, fullConfPath),
		FabricChannelID:     monetasa.Env(envFabricChannelID, defFabricChannelID),
		FabricChaincodeID:   monetasa.Env(envFabricChaincodeID, defFabricChaincodeID),
	}

	logger := log.New(os.Stdout)

	ms := connectToDB(cfg, logger)
	defer ms.Close()

	svc := newService(cfg, ms, logger)

	errs := make(chan error, 2)

	go startHTTPServer(svc, cfg.HTTPPort, logger, errs)

	go startGRPCServer(svc, cfg.GRPCPort, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err := <-errs
	logger.Error(fmt.Sprintf("Auth service terminated: %s", err))
}

func connectToDB(cfg config, logger log.Logger) *mgo.Session {
	ms, err := mongo.Connect(
		cfg.MongoURL,
		cfg.MongoConnectTimeout,
		cfg.MongoSocketTimeout,
		cfg.MongoDatabase,
		cfg.MongoUser,
		cfg.MongoPass,
	)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to Mongo: %s", err))
		os.Exit(1)
	}

	return ms
}

func newService(cfg config, ms *mgo.Session, logger log.Logger) auth.Service {
	// Initialization of the Fabric SDK
	fs := auth.FabricSetup{
		OrgAdmin:    cfg.FabricOrgAdmin,
		OrgName:     cfg.FabricOrgName,
		ConfigFile:  cfg.FabricConfFile,
		ChannelID:   cfg.FabricChannelID,
		ChaincodeID: cfg.FabricChaincodeID,
	}

	users := mongo.NewUserRepository(ms)
	hasher := bcrypt.New()
	idp := jwt.New(cfg.Secret)
	fn := fabric.NewFabricNetwork(&fs)
	if err := fn.Initialize(); err != nil {
		logger.Error(fmt.Sprintf("Unable to initialize the Fabric SDK: %v\n", err))
		os.Exit(1)
	}

	svc := auth.New(users, hasher, idp, fn)
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
