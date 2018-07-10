package api

import (
	"fmt"
	log "monetasa/logger"
	"monetasa/transactions"
	"time"
)

var _ transactions.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    transactions.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc transactions.Service, logger log.Logger) transactions.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) CreateUser(id string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method create_user for user %s took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.CreateUser(id)
}

func (lm *loggingMiddleware) Balance(userID, chanID string) (balance uint64, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method balance for user %s and channel %s took %s to complete", userID, chanID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.Balance(userID, chanID)
}
