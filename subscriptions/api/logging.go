package api

import (
	"fmt"
	"time"

	log "monetasa/logger"
	"monetasa/subscriptions"
)

var _ subscriptions.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    subscriptions.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc subscriptions.Service, logger log.Logger) subscriptions.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) CreateSubscription(token string, sub subscriptions.Subscription) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method create_subscription for user %s took %s to complete", sub.UserID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.CreateSubscription(token, sub)
}

func (lm *loggingMiddleware) ReadSubscriptions(token string) (subs []subscriptions.Subscription, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method read_subscriptions for user %s took %s to complete", token, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.ReadSubscriptions(token)
}
