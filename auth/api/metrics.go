package api

import (
	"time"

	"github.com/datapace/datapace/auth"

	"github.com/go-kit/kit/metrics"
)

var _ auth.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     auth.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc auth.Service, counter metrics.Counter, latency metrics.Histogram) auth.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) Register(key string, user auth.User) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "register").Add(1)
		ms.latency.With("method", "register").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Register(key, user)
}

func (ms *metricsMiddleware) InitAdmin(user auth.User, policies map[string]auth.Policy) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "init_admin").Add(1)
		ms.latency.With("method", "init_admin").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.InitAdmin(user, policies)
}

func (ms *metricsMiddleware) Login(user auth.User) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "login").Add(1)
		ms.latency.With("method", "login").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Login(user)
}

func (ms *metricsMiddleware) RecoverPassword(email string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "recover_password").Add(1)
		ms.latency.With("method", "recover_password").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RecoverPassword(email)
}

func (ms *metricsMiddleware) ValidateRecoveryToken(token string, id string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "validate_recovery_token").Add(1)
		ms.latency.With("method", "validate_recovery_token").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ValidateRecoveryToken(token, id)
}

func (ms *metricsMiddleware) UpdatePassword(token string, id string, password string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_password").Add(1)
		ms.latency.With("method", "update_password").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdatePassword(token, id, password)
}

func (ms *metricsMiddleware) UpdateUser(key string, user auth.User) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_user").Add(1)
		ms.latency.With("method", "update_user").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateUser(key, user)
}

func (ms *metricsMiddleware) ViewUser(key, ID string) (auth.User, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_user").Add(1)
		ms.latency.With("method", "view_user").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewUser(key, ID)
}

func (ms *metricsMiddleware) ViewEmail(key string) (auth.User, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_email").Add(1)
		ms.latency.With("method", "view_email").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewEmail(key)
}

func (ms *metricsMiddleware) ViewUserById(id string) (auth.User, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_user_by_id").Add(1)
		ms.latency.With("method", "view_user_by_id").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewUserById(id)
}

func (ms *metricsMiddleware) ListUsers(key string) ([]auth.User, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_users").Add(1)
		ms.latency.With("method", "list_users").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListUsers(key)
}

func (ms *metricsMiddleware) ListNonPartners(key string) ([]auth.User, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_non_partners").Add(1)
		ms.latency.With("method", "list_non_partners").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListNonPartners(key)
}

func (ms *metricsMiddleware) Identify(key string) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "identify").Add(1)
		ms.latency.With("method", "identify").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Identify(key)
}

func (ms *metricsMiddleware) Authorize(key string, action auth.Action, resource auth.Resource) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "authorize").Add(1)
		ms.latency.With("method", "authorize").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Authorize(key, action, resource)
}

func (ms *metricsMiddleware) Exists(id string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "exists").Add(1)
		ms.latency.With("method", "exists").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Exists(id)
}

func (ms *metricsMiddleware) AddPolicy(key string, policy auth.Policy) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "add_policy").Add(1)
		ms.latency.With("method", "add_policy").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.AddPolicy(key, policy)
}

func (ms *metricsMiddleware) ViewPolicy(key, id string) (auth.Policy, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_policy").Add(1)
		ms.latency.With("method", "view_policy").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewPolicy(key, id)
}

func (ms *metricsMiddleware) ListPolicies(key string) ([]auth.Policy, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_policies").Add(1)
		ms.latency.With("method", "list_policies").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListPolicies(key)
}

func (ms *metricsMiddleware) RemovePolicy(key, id string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "remove_policy").Add(1)
		ms.latency.With("method", "remove_policy").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RemovePolicy(key, id)
}

func (ms *metricsMiddleware) AttachPolicy(key, userID, policyID string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "attach_policy").Add(1)
		ms.latency.With("method", "attach_policy").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.AttachPolicy(key, userID, policyID)
}

func (ms *metricsMiddleware) DetachPolicy(key, userID, policyID string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "detach_policy").Add(1)
		ms.latency.With("method", "detach_policy").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.DetachPolicy(key, userID, policyID)
}
