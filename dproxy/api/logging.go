package api

import (
	"fmt"
	"time"

	"github.com/datapace/dproxy"
	log "github.com/datapace/logger"
)

var _ dproxy.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    dproxy.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc dproxy.Service, logger log.Logger) dproxy.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) CreateToken(url string, ttl, maxCalls int) (token string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method create_token for url %s ttl %d, max_calls %d and took %s to complete", url, ttl, maxCalls, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.CreateToken(url, ttl, maxCalls)
}

func (lm *loggingMiddleware) GetTargetURL(url string) (targeturl string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method get_target_url for url %s and took %s to complete", url, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.GetTargetURL(url)
}
