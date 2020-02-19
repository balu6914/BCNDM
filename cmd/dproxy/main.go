package main

import (
	"datapace"
	"datapace/dproxy"
	"datapace/logger"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	envHTTPPort  = "DATAPACE_PROXY_HTTP_PORT"
	defHTTPPort  = "9090"
	envTargetURL = "DATAPACE_PROXY_TARGET_URL"
	defTargetURL = "http://localhost"
)

type config struct {
	httpPort  string
	targetURL string
}

func main() {
	cfg := loadConfig()
	logger := logger.New(os.Stdout)
	errs := make(chan error, 2)

	go startHTTPServer(cfg.targetURL, cfg.httpPort, logger, errs)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err := <-errs
	logger.Error(fmt.Sprintf("Proxy service terminated: %s", err))

}

func loadConfig() config {
	return config{
		httpPort:  datapace.Env(envHTTPPort, defHTTPPort),
		targetURL: datapace.Env(envTargetURL, defTargetURL),
	}
}

func startHTTPServer(target, port string, logger logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("Executions HTTP service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, http.HandlerFunc(dproxy.MakeReverseProxy(target)))
}
