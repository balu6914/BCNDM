package api

import (
	access "datapace/access-control"
	log "datapace/logger"
	"fmt"
	"time"
)

var _ access.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    access.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc access.Service, logger log.Logger) access.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) RequestAccess(key, partner string) (id string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method request_access for key %s and partner %s and took %s to complete", key, partner, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RequestAccess(key, partner)
}

func (lm *loggingMiddleware) ListSentAccessRequests(key string, state access.State) (list []access.Request, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_sent_access_requests for key %s and state %s and took %s to complete", key, state, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ListSentAccessRequests(key, state)
}

func (lm *loggingMiddleware) ListReceivedAccessRequests(key string, state access.State) (list []access.Request, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_received_access_requests for key %s and state %s and took %s to complete", key, state, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ListReceivedAccessRequests(key, state)
}

func (lm *loggingMiddleware) ListPartners(id string) (list []string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_partners for user %s and took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ListPartners(id)
}

func (lm *loggingMiddleware) ListPotentialPartners(id string) (list []string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_potential_partners for user %s and took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ListPotentialPartners(id)
}

func (lm *loggingMiddleware) ApproveAccessRequest(key, id string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method approve_access_request for key %s and access request %s and took %s to complete", key, id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ApproveAccessRequest(key, id)
}

func (lm *loggingMiddleware) RevokeAccessRequest(key, id string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method revoke_access_request for key %s and access request %s and took %s to complete", key, id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RevokeAccessRequest(key, id)
}
