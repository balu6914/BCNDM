package main

import (
	"datapace"
	"datapace/dproxy"
	"datapace/dproxy/api"
	httpapi "datapace/dproxy/api/http"
	"datapace/dproxy/jwt"
	"datapace/dproxy/persistence/postgres"
	"datapace/logger"
	log "datapace/logger"
	"fmt"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/jmoiron/sqlx"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	defHTTPPort      = "9090"
	defJWTSecret     = "examplesecret"
	defLocalFsRoot   = "/tmp/test"
	defDBHost        = "localhost"
	defDBPort        = "5432"
	defDBUser        = "dproxy"
	defDBPass        = "dproxy"
	defDBName        = "dproxy"
	defDBSSLMode     = "disable"
	defDBSSLCert     = ""
	defDBSSLKey      = ""
	defDBSSLRootCert = ""

	envHTTPPort      = "DATAPACE_PROXY_HTTP_PORT"
	envJWTSecret     = "DATAPACE_JWT_SECRET"
	envLocalFsRoot   = "DATAPACE_LOCAL_FS_ROOT"
	envDBHost        = "DATAPACE_DPROXY_DB_HOST"
	envDBPort        = "DATAPACE_DPROXY_DB_PORT"
	envDBUser        = "DATAPACE_DPROXY_DB_USER"
	envDBPass        = "DATAPACE_DPROXY_DB_PASS"
	envDBName        = "DATAPACE_DPROXY_DB"
	envDBSSLMode     = "DATAPACE_DPROXY_DB_SSL_MODE"
	envDBSSLCert     = "DATAPACE_DPROXY_DB_SSL_CERT"
	envDBSSLKey      = "DATAPACE_DPROXY_DB_SSL_KEY"
	envDBSSLRootCert = "DATAPACE_DPROXY_DB_SSL_ROOT_CERT"
)

type config struct {
	httpPort    string
	jwtSecret   string
	localFsRoot string
	dbConfig    postgres.Config
}

func main() {
	cfg := loadConfig()
	logger := logger.New(os.Stdout)
	errs := make(chan error, 2)
	svc := newService(cfg.jwtSecret, cfg.dbConfig, logger)
	r := httpapi.NewReverseProxy(svc, logger)
	f := httpapi.NewFsProxy(svc, cfg.localFsRoot, logger)
	go startHTTPServer(svc, r, f, cfg.httpPort, logger, errs)
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
		httpPort:    datapace.Env(envHTTPPort, defHTTPPort),
		jwtSecret:   datapace.Env(envJWTSecret, defJWTSecret),
		localFsRoot: datapace.Env(envLocalFsRoot, defLocalFsRoot),
		dbConfig:    dbConfig,
	}
}

func startHTTPServer(svc dproxy.Service, rp *httpapi.ReverseProxy, fs *httpapi.FsProxy, port string, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("Proxy HTTP service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, httpapi.MakeHandler(svc, rp, fs))
}
