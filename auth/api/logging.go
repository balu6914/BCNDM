package api

import (
	"fmt"
	"time"

	"github.com/datapace/datapace/auth"
	log "github.com/datapace/datapace/logger"
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

func (lm *loggingMiddleware) Register(key string, user auth.User) (ID string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method register for user %s with role %s took %s to complete", user.Email, user.Role, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.Register(key, user)
}

func (lm *loggingMiddleware) InitAdmin(user auth.User, policies map[string]auth.Policy) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method init_admin for admin %s took %s to complete", user.Email, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.InitAdmin(user, policies)
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

func (lm *loggingMiddleware) UpdateUser(key string, user auth.User) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method update_user for key %s and user %s took %s to complete", key, user.Email, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.UpdateUser(key, user)
}

func (lm *loggingMiddleware) ViewUser(key, ID string) (user auth.User, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method view_user for key %s and ID %s and took %s to complete", key, ID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ViewUser(key, ID)
}

func (lm *loggingMiddleware) ViewEmail(key string) (user auth.User, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method view for key %s  and took %s to complete", key, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ViewEmail(key)
}

func (lm *loggingMiddleware) ViewUserById(id string) (u auth.User, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method view_user_by_id for ID %s and took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ViewUserById(id)
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

func (lm *loggingMiddleware) Authorize(key string, action auth.Action, resource auth.Resource) (id string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method authorize for user %s for action %d over resource %s took %s to complete",
			id, action, resource.ResourceType(), time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Authorize(key, action, resource)
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

func (lm *loggingMiddleware) AddPolicy(key string, policy auth.Policy) (id string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method add_policy for policy %s took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.AddPolicy(key, policy)
}

func (lm *loggingMiddleware) ViewPolicy(key, id string) (policy auth.Policy, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method view_policy for policy %s took %s to complete", id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.ViewPolicy(key, id)
}

func (lm *loggingMiddleware) ListPolicies(key string) (policy []auth.Policy, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method list_policies for policies took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.ListPolicies(key)
}

func (lm *loggingMiddleware) RemovePolicy(key, policyID string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method remove_policy for policy %s took %s to complete", policyID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.RemovePolicy(key, policyID)
}

func (lm *loggingMiddleware) AttachPolicy(key, policyID, userID string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method attach_policy for user %s took %s to complete", userID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.AttachPolicy(key, policyID, userID)
}

func (lm *loggingMiddleware) DetachPolicy(key, policyID, userID string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method detach_policy for user %s took %s to complete", userID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.DetachPolicy(key, policyID, userID)
}
