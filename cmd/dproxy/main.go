package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/datapace/datapace"

	"github.com/datapace/datapace/dproxy"
	"github.com/datapace/datapace/dproxy/api"
	httpapi "github.com/datapace/datapace/dproxy/api/http"
	"github.com/datapace/datapace/dproxy/jwt"
	"github.com/datapace/datapace/dproxy/persistence/postgres"
	"github.com/datapace/datapace/logger"
	log "github.com/datapace/datapace/logger"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/jmoiron/sqlx"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	defHTTPProto      = "http"
	defHTTPHost       = "localhost"
	defHTTPPort       = "8087"
	defJWTSecret      = "examplesecret"
	defLocalFsRoot    = "/tmp/test"
	defDBHost         = "localhost"
	defDBPort         = "5432"
	defDBUser         = "dproxy"
	defDBPass         = "dproxy"
	defDBName         = "dproxy"
	defDBSSLMode      = "disable"
	defDBSSLCert      = ""
	defDBSSLKey       = ""
	defDBSSLRootCert  = ""
	defFSPathPrefix   = "/fs"
	defHTTPPathPrefix = "/http"

	envHTTPProto      = "DATAPACE_PROXY_HTTP_PROTO"
	envHTTPHost       = "DATAPACE_PROXY_HTTP_HOST"
	envHTTPPort       = "DATAPACE_PROXY_HTTP_PORT"
	envJWTSecret      = "DATAPACE_JWT_SECRET"
	envLocalFsRoot    = "DATAPACE_LOCAL_FS_ROOT"
	envDBHost         = "DATAPACE_DPROXY_DB_HOST"
	envDBPort         = "DATAPACE_DPROXY_DB_PORT"
	envDBUser         = "DATAPACE_DPROXY_DB_USER"
	envDBPass         = "DATAPACE_DPROXY_DB_PASS"
	envDBName         = "DATAPACE_DPROXY_DB"
	envDBSSLMode      = "DATAPACE_DPROXY_DB_SSL_MODE"
	envDBSSLCert      = "DATAPACE_DPROXY_DB_SSL_CERT"
	envDBSSLKey       = "DATAPACE_DPROXY_DB_SSL_KEY"
	envDBSSLRootCert  = "DATAPACE_DPROXY_DB_SSL_ROOT_CERT"
	envFSPathPrefix   = "DATAPACE_DPROXY_FS_PATH_PREFIX"
	envHTTPPathPrefix = "DATAPACE_DPROXY_HTTP_PATH_PREFIX"
)

type config struct {
	httpProto      string
	httpHost       string
	httpPort       string
	jwtSecret      string
	localFsRoot    string
	dbConfig       postgres.Config
	fsPathPrefix   string
	httpPathPrefix string
}

func main() {
	cfg := loadConfig()
	logger := logger.New(os.Stdout)
	errs := make(chan error, 2)
	svc := newService(cfg.jwtSecret, cfg.dbConfig, logger)
	r := httpapi.NewReverseProxy(svc, cfg.httpPathPrefix, logger)
	f := httpapi.NewFsProxy(svc, cfg.localFsRoot, cfg.fsPathPrefix, logger)
	go startHTTPServer(svc, r, f, cfg.httpPort, fmt.Sprintf("%s://%s:%s", cfg.httpProto, cfg.httpHost, cfg.httpPort), logger, errs)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	err := <-errs
	logger.Error(fmt.Sprintf("Proxy service terminated: %s", err))
}

func newService(jwtSecret string, dbConfig postgres.Config, logger log.Logger) dproxy.Service {
	tokenService := jwt.NewService(jwtSecret)
	db := connectToDB(dbConfig, logger)
	eventsRepo := postgres.NewEventsRepository(db)
	svc := dproxy.NewService(tokenService, eventsRepo)
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

func connectToDB(dbConfig postgres.Config, logger logger.Logger) *sqlx.DB {
	db, err := postgres.Connect(dbConfig)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to postgres: %s", err))
		os.Exit(1)
	}
	return db
}

func loadConfig() config {
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
	return config{
		httpProto:      datapace.Env(envHTTPProto, defHTTPProto),
		httpHost:       datapace.Env(envHTTPHost, defHTTPHost),
		httpPort:       datapace.Env(envHTTPPort, defHTTPPort),
		jwtSecret:      datapace.Env(envJWTSecret, defJWTSecret),
		localFsRoot:    datapace.Env(envLocalFsRoot, defLocalFsRoot),
		dbConfig:       dbConfig,
		fsPathPrefix:   datapace.Env(envFSPathPrefix, defFSPathPrefix),
		httpPathPrefix: datapace.Env(envHTTPPathPrefix, defHTTPPathPrefix),
	}
}

func startHTTPServer(svc dproxy.Service, rp *httpapi.ReverseProxy, fs *httpapi.FsProxy, port, dProxyRootUrl string, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("Proxy HTTP service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, httpapi.MakeHandler(svc, rp, fs, dProxyRootUrl))
}
