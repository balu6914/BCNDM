package main

import (
	"datapace"
	"datapace/dproxy"
	"datapace/dproxy/api"
	httpapi "datapace/dproxy/api/http"
	"datapace/dproxy/jwt"
	"datapace/logger"
	log "datapace/logger"
	"fmt"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	envHTTPPort  = "DATAPACE_PROXY_HTTP_PORT"
	defHTTPPort  = "9090"
	envJWTSecret = "DATAPACE_JWT_SECRET"
	defJWTSecret = "examplesecret"
)

type config struct {
	httpPort  string
	jwtSecret string
}

func main() {
	cfg := loadConfig()
	logger := logger.New(os.Stdout)
	errs := make(chan error, 2)
	svc := newService(cfg.jwtSecret, logger)
	go startHTTPServer(svc, cfg.httpPort, logger, errs)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	err := <-errs
	logger.Error(fmt.Sprintf("Proxy service terminated: %s", err))
}

func newService(jwtSecret string, logger log.Logger) dproxy.Service {
	svc := jwt.NewService(jwtSecret)
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

func loadConfig() config {
	return config{
		httpPort:  datapace.Env(envHTTPPort, defHTTPPort),
		jwtSecret: datapace.Env(envJWTSecret, defJWTSecret),
	}
}

func startHTTPServer(svc dproxy.Service, port string, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("Proxy HTTP service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, httpapi.MakeHandler(svc))
}
