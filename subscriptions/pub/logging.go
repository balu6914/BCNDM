package pub

import (
	"fmt"
	log "github.com/datapace/datapace/logger"
	"time"
)

type (
	loggingMiddleware struct {
		svc    Service
		logger log.Logger
	}
)

// NewLoggingMiddleware adds logging facilities to the core service.
func NewLoggingMiddleware(svc Service, logger log.Logger) Service {
	return &loggingMiddleware{svc, logger}
}

func (lm loggingMiddleware) Publish(evt SubscriptionCreateEvent, toUserId string) (msgId uint64, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method pub.Publish(evt=%s, toUserId=%s) took %s to complete", evt, toUserId, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())
	return lm.Publish(evt, toUserId)
}
