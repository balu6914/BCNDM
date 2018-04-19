
package api

import (
	"fmt"
	"time"

	log "monetasa/logger"
	"monetasa/auth"
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

func (lm *loggingMiddleware) List(key string) (useres []auth.User, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list for key %s took %s to complete", key, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.List(key)
}

func (lm *loggingMiddleware) Delete(key string, id string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method remove for key %s and user %s took %s to complete", key, id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Delete(key, id)
}


func (lm *loggingMiddleware) Identity(key string) (id string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method identity for user %s took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Identity(key)
}

func (lm *loggingMiddleware) CanAccess(key string, id string) (pub string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method can_access for key %s, channel %s and publisher %s took %s to complete", key, id, pub, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.CanAccess(key, id)
}
