package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"monetasa/auth"
	"monetasa/auth/api"
	"monetasa/auth/bcrypt"
	"monetasa/auth/fabric"
	"monetasa/auth/jwt"
	"monetasa/auth/mongo"
	log "monetasa/logger"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	port                   int    = 8080
	defMongoURL            string = "0.0.0.0"
	defMongoUser           string = ""
	defMongoPass           string = ""
	defMongoDatabase       string = "auth"
	defMongoPort           int    = 27017
	defMongoConnectTimeout int    = 5000
	defMongoSocketTimeout  int    = 5000
	defAuthURL             string = "0.0.0.0"
	defSecret              string = "monetasa"

	envMongoURL string = "MONETASA_AUTH_MONGO_URL"
	envAuthURL  string = "MONETASA_AUTH_URL"
	envSecret   string = "MONETASA_AUTH_SECRET"

	defFabricOrgAdmin    string = "admin"
	defFabricOrgName     string = "org1"
	defFabricConfFile    string = "/src/monetasa/config/fabric/config.yaml"
	defFabricChannelID   string = "myc"
	defFabricChaincodeID string = "token"
)

type config struct {
	Port                int
	AuthURL             string
	PostgresURL         string
	MongoURL            string
	MongoUser           string
	MongoPass           string
	MongoDatabase       string
	MongoPort           int
	MongoConnectTimeout int
	MongoSocketTimeout  int
	Secret              string
	FabricOrgAdmin      string
	FabricOrgName       string
	FabricConfFile      string
	FabricChannelID     string
	FabricChaincodeID   string
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func main() {
	cfg := config{
		Port:                port,
		AuthURL:             getenv(envAuthURL, defAuthURL),
		MongoURL:            getenv(envMongoURL, defMongoURL),
		MongoUser:           defMongoUser,
		MongoPass:           defMongoPass,
		MongoDatabase:       defMongoDatabase,
		MongoPort:           defMongoPort,
		MongoConnectTimeout: defMongoConnectTimeout,
		MongoSocketTimeout:  defMongoSocketTimeout,
		Secret:              getenv(envSecret, defSecret),
		FabricOrgAdmin:      defFabricOrgAdmin,
		FabricOrgName:       defFabricOrgName,
		FabricConfFile:      os.Getenv("GOPATH") + defFabricConfFile,
		FabricChannelID:     defFabricChannelID,
		FabricChaincodeID:   defFabricChaincodeID,
	}

	logger := log.New(os.Stdout)

	ms, err := mongo.Connect(cfg.MongoURL, cfg.MongoConnectTimeout, cfg.MongoSocketTimeout,
		cfg.MongoDatabase, cfg.MongoUser, cfg.MongoPass)
	if err != nil {
		logger.Error("Failed to connect to Mongo.")
		os.Exit(1)
	}
	defer ms.Close()

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

	fields := []string{"method"}
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "auth",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fields),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "auth",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fields),
	)

	errs := make(chan error, 2)

	go func() {
		p := fmt.Sprintf(":%d", cfg.Port)
		logger.Info(fmt.Sprintf("Auth service started, exposed port %d", cfg.Port))
		errs <- http.ListenAndServe(p, api.MakeHandler(svc))
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("Auth service terminated: %s", err))
}
