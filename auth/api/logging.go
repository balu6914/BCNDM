package api

import (
	"fmt"
	"time"

	"github.com/datapace/auth"
	log "github.com/datapace/logger"
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

func (lm *loggingMiddleware) Register(key string, user auth.User) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method register for user %s took %s to complete", user.Email, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.Register(key, user)
}

func (lm *loggingMiddleware) InitAdmin(user auth.User) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method init_admin for admin %s took %s to complete", user.Email, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.InitAdmin(user)
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

func (lm *loggingMiddleware) ListUsers(key string) (list []auth.User, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_users for key %s and took %s to complete", key, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ListUsers(key)
}

func (lm *loggingMiddleware) ListNonPartners(key string) (list []auth.User, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_non_partners for key %s and took %s to complete", key, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ListNonPartners(key)
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

func (lm *loggingMiddleware) Exists(id string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method exists for id %s took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Exists(id)
}
