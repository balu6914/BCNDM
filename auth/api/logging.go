package api

import (
	"fmt"
	"time"

	"datapace/auth"
	log "datapace/logger"
)

var _ auth.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    auth.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc auth.Service, logger log.Logger) auth.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) Register(user auth.User) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method register for user %s took %s to complete", user.Email, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.Register(user)
}

func (lm *loggingMiddleware) Login(user auth.User) (token string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method login for user %s took %s to complete", user.Email, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Login(user)
}

func (lm *loggingMiddleware) Update(key string, user auth.User) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method update for key %s and user %s took %s to complete", key, user.Email, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Update(key, user)
}

func (lm *loggingMiddleware) UpdatePassword(key string, old string, user auth.User) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method update_password for key %s and user %s took %s to complete", key, user.Email, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.UpdatePassword(key, old, user)
}

func (lm *loggingMiddleware) View(key string) (user auth.User, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method view for key %s and took %s to complete", key, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.View(key)
}

func (lm *loggingMiddleware) List(key string) (list []auth.User, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list for key %s and took %s to complete", key, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.List(key)
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

func (lm *loggingMiddleware) ListSentAccessRequests(key string, state auth.State) (list []auth.AccessRequest, err error) {
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

func (lm *loggingMiddleware) ListReceivedAccessRequests(key string, state auth.State) (list []auth.AccessRequest, err error) {
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

func (lm *loggingMiddleware) RejectAccessRequest(key, id string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method reject_access_request for key %s and access request %s and took %s to complete", key, id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RejectAccessRequest(key, id)
}

func (lm *loggingMiddleware) Identify(key string) (id string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method identify for user %s took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Identify(key)
}
