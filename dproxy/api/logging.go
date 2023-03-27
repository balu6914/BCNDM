package api

import (
	"fmt"
	"time"

	"github.com/datapace/datapace/dproxy"
	"github.com/datapace/datapace/dproxy/persistence"
	log "github.com/datapace/datapace/logger"
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

func (lm *loggingMiddleware) CreateToken(url string, ttl, maxCalls int, maxUnit, subID string) (token string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method create_token for url %s subscriptionID %s ttl %d, max_calls %d, max_unit %s and took %s to complete", url, subID, ttl, maxCalls, maxUnit, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.CreateToken(url, ttl, maxCalls, maxUnit, subID)
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

func (lm *loggingMiddleware) List(q persistence.Query) (page []persistence.Event, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.List(q)
}
