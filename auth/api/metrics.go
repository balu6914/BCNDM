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

func (ms *metricsMiddleware) InitAdmin(user auth.User) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "init_admin").Add(1)
		ms.latency.With("method", "init_admin").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.InitAdmin(user)
}

func (ms *metricsMiddleware) Login(user auth.User) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "login").Add(1)
		ms.latency.With("method", "login").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Login(user)
}

func (ms *metricsMiddleware) Update(key string, user auth.User) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "update").Add(1)
		ms.latency.With("method", "update").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Update(key, user)
}

func (ms *metricsMiddleware) UpdatePassword(key string, old string, user auth.User) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_password").Add(1)
		ms.latency.With("method", "update_password").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdatePassword(key, old, user)
}

func (ms *metricsMiddleware) View(key, ID string) (auth.User, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view").Add(1)
		ms.latency.With("method", "view").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.View(key, ID)
}

func (ms *metricsMiddleware) ViewEmail(key string) (auth.User, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_email").Add(1)
		ms.latency.With("method", "view_email").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewEmail(key)
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

func (ms *metricsMiddleware) Exists(id string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "exists").Add(1)
		ms.latency.With("method", "exists").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Exists(id)
}
