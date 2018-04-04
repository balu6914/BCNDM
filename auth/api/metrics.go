package api

import (
	"time"

	"github.com/go-kit/kit/metrics"
	"monetasa/monetasa/auth"
	"gitlab.com/drasko/monetasa/auth"
)

var _ auth.Service = (*metricService)(nil)

type metricService struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc auth.Service
}

// NewMetricService instruments adapter by tracking request count and latency.
func NewMetricService(counter metrics.Counter, latency metrics.Histogram, s auth.Service) auth.Service {
	return &metricService{
		counter: counter,
		latency: latency,
		Service: s,
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

func (ms *metricsMiddleware) Update(user manager.User) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "update").Add(1)
		ms.latency.With("method", "update").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateClient(key, user)
}

func (ms *metricsMiddleware) View(id string) (manager.User, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_client").Add(1)
		ms.latency.With("method", "view_client").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewClient(id)
}

func (ms *metricsMiddleware) List() ([]manager.User, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_clients").Add(1)
		ms.latency.With("method", "list_clients").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListClients(key)
}

func (ms *metricsMiddleware) Delete(id string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "remove").Add(1)
		ms.latency.With("method", "remove").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Delete(key, id)
}

func (ms *metricsMiddleware) CanAccess(key string, id string) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "can_access").Add(1)
		ms.latency.With("method", "can_access").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CanAccess(key, id)
