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

func (lm *loggingMiddleware) Balance(userID string) (balance uint64, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method balance for user %s took %s to complete", userID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Balance(userID)
}

func (lm *loggingMiddleware) Transfer(from, to string, value uint64) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method transfer from %s to %s with amount %d took %s to complete", from, to, value, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Transfer(from, to, value)
}

func (lm *loggingMiddleware) BuyTokens(account string, value uint64) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method buy_tokens for account %s and amount %d took %s to complete", account, value, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.BuyTokens(account, value)
}
