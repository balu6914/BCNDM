package api

import (
	"time"

	"monetasa/auth"

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

func (ms *metricsMiddleware) Register(user auth.User) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "register").Add(1)
		ms.latency.With("method", "register").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Register(user)
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

func (ms *metricsMiddleware) View(key string) (auth.User, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view").Add(1)
		ms.latency.With("method", "view").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.View(key)
}

func (ms *metricsMiddleware) List(key string) ([]auth.User, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list").Add(1)
		ms.latency.With("method", "list").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.List(key)
}

func (ms *metricsMiddleware) Delete(key string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "remove").Add(1)
		ms.latency.With("method", "remove").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Delete(key)
}

func (ms *metricsMiddleware) Identity(key string) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "identity").Add(1)
		ms.latency.With("method", "identity").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Identity(key)
}

func (ms *metricsMiddleware) CanAccess(key, id string) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "can_access").Add(1)
		ms.latency.With("method", "can_access").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CanAccess(key, id)
}
