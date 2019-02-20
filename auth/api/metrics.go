package api

import (
	"time"

	"datapace/auth"

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

func (ms *metricsMiddleware) UpdatePassword(key string, old string, user auth.User) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_password").Add(1)
		ms.latency.With("method", "update_password").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdatePassword(key, old, user)
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

func (ms *metricsMiddleware) RequestAccess(key, partner string) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "request_access").Add(1)
		ms.latency.With("method", "request_access").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RequestAccess(key, partner)
}

func (ms *metricsMiddleware) ListSentAccessRequests(key string, state auth.State) ([]auth.AccessRequest, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_sent_access_requests").Add(1)
		ms.latency.With("method", "list_sent_access_requests").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListSentAccessRequests(key, state)
}

func (ms *metricsMiddleware) ListReceivedAccessRequests(key string, state auth.State) ([]auth.AccessRequest, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_received_access_requests").Add(1)
		ms.latency.With("method", "list_received_access_requests").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListReceivedAccessRequests(key, state)
}

func (ms *metricsMiddleware) ListPartners(id string) ([]string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_partners").Add(1)
		ms.latency.With("method", "list_partners").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListPartners(id)
}

func (ms *metricsMiddleware) ApproveAccessRequest(key, id string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "approve_access_request").Add(1)
		ms.latency.With("method", "approve_access_request").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ApproveAccessRequest(key, id)
}

func (ms *metricsMiddleware) RejectAccessRequest(key, id string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "reject_access_request").Add(1)
		ms.latency.With("method", "reject_access_request").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RejectAccessRequest(key, id)
}

func (ms *metricsMiddleware) Identify(key string) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "identify").Add(1)
		ms.latency.With("method", "identify").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Identify(key)
}
