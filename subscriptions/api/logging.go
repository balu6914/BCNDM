package api

import (
	"fmt"
	"time"

	log "github.com/datapace/logger"
	"github.com/datapace/subscriptions"
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

func (lm *loggingMiddleware) AddSubscription(userID, userToken string, sub subscriptions.Subscription) (id string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method add_subscription for user %s and stream %s took %s to complete", userID, sub.StreamID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.AddSubscription(userID, userToken, sub)
}

func (lm *loggingMiddleware) ViewSubscription(userID, subID string) (sub subscriptions.Subscription, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method view_subscription took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.ViewSubscription(userID, subID)
}

func (lm *loggingMiddleware) SearchSubscriptions(query subscriptions.Query) (page subscriptions.Page, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method search_subscriptions took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.SearchSubscriptions(query)
}

func (lm *loggingMiddleware) ViewSubByUserAndStream(userID, streamID string) (sub subscriptions.Subscription, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method view_sub_by_user_and_stream took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}

		lm.logger.Warn(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ViewSubByUserAndStream(userID, streamID)
}
