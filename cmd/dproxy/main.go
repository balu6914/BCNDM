package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/datapace/datapace"
	"github.com/datapace/datapace/dproxy"
	"github.com/datapace/datapace/dproxy/api"
	httpapi "github.com/datapace/datapace/dproxy/api/http"
	"github.com/datapace/datapace/dproxy/jwt"
	"github.com/datapace/datapace/dproxy/persistence"
	"github.com/datapace/datapace/dproxy/persistence/mongo"
	"github.com/datapace/datapace/dproxy/persistence/postgres"
	"github.com/datapace/datapace/errors"
	"github.com/datapace/datapace/logger"
	log "github.com/datapace/datapace/logger"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	defHTTPProto         = "http"
	defHTTPHost          = "localhost"
	defHTTPPort          = "8087"
	defJWTSecret         = "examplesecret"
	defLocalFsRoot       = "/tmp/test"
	defDBType            = "mongo"
	defDBHost            = "0.0.0.0"
	defDBPort            = "27017"
	defDBUser            = ""
	defDBPass            = ""
	defDBName            = "dproxy"
	defDBSSLMode         = "disable"
	defDBSSLCert         = ""
	defDBSSLKey          = ""
	defDBSSLRootCert     = ""
	defFSPathPrefix      = "/fs"
	defHTTPPathPrefix    = "/http"
	defHttpTlsSkipVerify = "true"
	defEncKey            = ""
	defStandalone        = "true"

	envHTTPProto         = "DATAPACE_PROXY_HTTP_PROTO"
	envHTTPHost          = "DATAPACE_PROXY_HTTP_HOST"
	envHTTPPort          = "DATAPACE_PROXY_HTTP_PORT"
	envJWTSecret         = "DATAPACE_JWT_SECRET"
	envLocalFsRoot       = "DATAPACE_LOCAL_FS_ROOT"
	envDBType            = "DATAPACE_DPROXY_DB_TYPE"
	envDBHost            = "DATAPACE_DPROXY_DB_HOST"
	envDBPort            = "DATAPACE_DPROXY_DB_PORT"
	envDBUser            = "DATAPACE_DPROXY_DB_USER"
	envDBPass            = "DATAPACE_DPROXY_DB_PASS"
	envDBName            = "DATAPACE_DPROXY_DB_NAME"
	envDBSSLMode         = "DATAPACE_DPROXY_DB_SSL_MODE"
	envDBSSLCert         = "DATAPACE_DPROXY_DB_SSL_CERT"
	envDBSSLKey          = "DATAPACE_DPROXY_DB_SSL_KEY"
	envDBSSLRootCert     = "DATAPACE_DPROXY_DB_SSL_ROOT_CERT"
	envFSPathPrefix      = "DATAPACE_DPROXY_FS_PATH_PREFIX"
	envHTTPPathPrefix    = "DATAPACE_DPROXY_HTTP_PATH_PREFIX"
	envHttpTlsSkipVerify = "DATAPACE_DPROXY_HTTP_TLS_SKIP_VERIFY"
	envEncKey            = "DATAPACE_DPROXY_ENCRYPTION_KEY"
	envStandalone        = "DATAPACE_DPROXY_STANDALONE"
)

type config struct {
	httpProto         string
	httpHost          string
	httpPort          string
	jwtSecret         string
	localFsRoot       string
	dbConfig          postgres.Config
	fsPathPrefix      string
	httpPathPrefix    string
	httpTlsSkipVerify bool
	dbType            string
	encKey            string
	standalone        bool
}

func main() {
	logger := logger.New(os.Stdout)
	cfg := loadConfig(logger)
	errs := make(chan error, 2)
	eventsRepository, err := connectToEventsRepository()
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while connecting to events repository: %s. Exiting.", err))
		os.Exit(1)
	}
	key, err := base64.StdEncoding.DecodeString(cfg.encKey)
	if err != nil {
		logger.Error(fmt.Sprintf("Error reading AES key: %s", err))
		os.Exit(1)
	}
	svc := newService(cfg.jwtSecret, eventsRepository, key, logger)
	r := httpapi.NewReverseProxy(svc, cfg.httpPathPrefix, cfg.httpTlsSkipVerify, logger)
	f := httpapi.NewFsProxy(svc, cfg.localFsRoot, cfg.fsPathPrefix, logger)
	url := fmt.Sprintf("%s://%s:%s/dproxy", cfg.httpProto, cfg.httpHost, cfg.httpPort)
	if cfg.standalone {
		url = fmt.Sprintf("%s://%s:%s", cfg.httpProto, cfg.httpHost, cfg.httpPort)
	}

	go startHTTPServer(svc, r, f, cfg.httpPort, url, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	err = <-errs
	logger.Error(fmt.Sprintf("Proxy service terminated: %s", err))
}

func newService(jwtSecret string, eventsRepository persistence.EventRepository, key []byte, logger log.Logger) dproxy.Service {
	tokenService := jwt.NewService(jwtSecret)
	svc := dproxy.NewService(tokenService, eventsRepository, key)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "dproxy",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "dproxy",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)
	return svc
}

func connectToEventsRepository() (persistence.EventRepository, error) {
	dbType := datapace.Env(envDBType, defDBType)
	var eventsRepo persistence.EventRepository
	switch dbType {
	case "postgres":
		dbConfig := postgres.Config{
			Host:        datapace.Env(envDBHost, defDBHost),
			Port:        datapace.Env(envDBPort, defDBPort),
			User:        datapace.Env(envDBUser, defDBUser),
			Pass:        datapace.Env(envDBPass, defDBPass),
			Name:        datapace.Env(envDBName, defDBName),
			SSLMode:     datapace.Env(envDBSSLMode, defDBSSLMode),
			SSLCert:     datapace.Env(envDBSSLCert, defDBSSLCert),
			SSLKey:      datapace.Env(envDBSSLKey, defDBSSLKey),
			SSLRootCert: datapace.Env(envDBSSLRootCert, defDBSSLRootCert),
		}
		db, err := postgres.Connect(dbConfig)
		if err != nil {
			return nil, err
		}
		eventsRepo = postgres.NewEventsRepository(db)
	case "mongo":
		ms, err := mongo.Connect(
			datapace.Env(envDBHost, defDBHost)+":"+datapace.Env(envDBPort, defDBPort),
			5000,
			5000,
			datapace.Env(envDBName, defDBName),
			datapace.Env(envDBUser, defDBUser),
			datapace.Env(envDBPass, defDBPass),
		)
		if err != nil {
			return nil, err
		}
		eventsRepo = mongo.NewEventsRepository(ms)
	default:
		return nil, errors.New("unknown database type")
	}
	return eventsRepo, nil
}

func loadConfig(logger log.Logger) config {

	standalone, err := strconv.ParseBool(datapace.Env(envStandalone, defStandalone))
	if err != nil {
		logger.Error(fmt.Sprintf("Invalid %s value: %s", envStandalone, err.Error()))
	}
	httpTlsSkipVerify, err := strconv.ParseBool(datapace.Env(envHttpTlsSkipVerify, defHttpTlsSkipVerify))
	if err != nil {
		logger.Error(fmt.Sprintf("Invalid %s value: %s", envHttpTlsSkipVerify, err.Error()))
	}
	return config{
		httpProto:         datapace.Env(envHTTPProto, defHTTPProto),
		httpHost:          datapace.Env(envHTTPHost, defHTTPHost),
		httpPort:          datapace.Env(envHTTPPort, defHTTPPort),
		jwtSecret:         datapace.Env(envJWTSecret, defJWTSecret),
		localFsRoot:       datapace.Env(envLocalFsRoot, defLocalFsRoot),
		fsPathPrefix:      datapace.Env(envFSPathPrefix, defFSPathPrefix),
		httpPathPrefix:    datapace.Env(envHTTPPathPrefix, defHTTPPathPrefix),
		httpTlsSkipVerify: httpTlsSkipVerify,
		encKey:            datapace.Env(envEncKey, defEncKey),
		standalone:        standalone,
	}
}

func startHTTPServer(svc dproxy.Service, rp *httpapi.ReverseProxy, fs *httpapi.FsProxy, port, dProxyRootUrl string, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("Proxy HTTP service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, httpapi.MakeHandler(svc, rp, fs, dProxyRootUrl))
}
