package api

import (
	"fmt"
	log "github.com/datapace/datapace/logger"
	"github.com/datapace/datapace/terms"
	"time"
)

var _ terms.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    terms.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc terms.Service, logger log.Logger) terms.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm loggingMiddleware) CreateTerms(terms terms.Terms) (id string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method create_terms %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())
	return lm.svc.CreateTerms(terms)
}

func (lm loggingMiddleware) ValidateTerms(terms terms.Terms) (isValid bool, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method validate_terms %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())
	return lm.svc.ValidateTerms(terms)
}
