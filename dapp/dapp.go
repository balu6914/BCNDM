package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	api "./api"
	mongo "./mongo"
)

const (
	port       int    = 8081
	defDappURL string = "0.0.0.0"

	defMongoURL            string = "0.0.0.0"
	defMongoUser           string = ""
	defMongoPass           string = ""
	defMongoDatabase       string = "api"
	defMongoPort           int    = 27017
	defMongoConnectTimeout int    = 5000
	defMongoSocketTimeout  int    = 5000

	envMongoURL string = "MONETASA_DAPP_MONGO_URL"
	envDappURL  string = "MONETASA_DAPP_URL"
)

type config struct {
	Port    int
	dappURL string

	MongoURL            string
	MongoUser           string
	MongoPass           string
	MongoDatabase       string
	MongoPort           int
	MongoConnectTimeout int
	MongoSocketTimeout  int
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
		dappURL:             getenv(envDappURL, defDappURL),
		MongoURL:            getenv(envMongoURL, defMongoURL),
		MongoUser:           defMongoUser,
		MongoPass:           defMongoPass,
		MongoDatabase:       defMongoDatabase,
		MongoPort:           defMongoPort,
		MongoConnectTimeout: defMongoConnectTimeout,
		MongoSocketTimeout:  defMongoSocketTimeout,
	}

	ms, err := mongo.Connect(cfg.MongoURL, cfg.MongoConnectTimeout, cfg.MongoSocketTimeout,
		cfg.MongoDatabase, cfg.MongoUser, cfg.MongoPass)
	if err != nil {
		// logger.Error("Failed to connect to Mongo.", zap.Error(err))
		os.Exit(1)
	}
	defer ms.Close()

	sr := mongo.NewStreamRepository(ms)

	errs := make(chan error, 2)

	go func() {
		p := fmt.Sprintf(":%d", cfg.Port)
		errs <- http.ListenAndServe(p, api.MakeHandler(sr))
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
}
