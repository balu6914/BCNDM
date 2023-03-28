package api

import (
	"time"

	access "github.com/datapace/datapace/access-control"

	"github.com/go-kit/kit/metrics"
)

var _ access.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     access.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc access.Service, counter metrics.Counter, latency metrics.Histogram) access.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) RequestAccess(key, partner string) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "request_access").Add(1)
		ms.latency.With("method", "request_access").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RequestAccess(key, partner)
}

func (ms *metricsMiddleware) ListSentAccessRequests(key string, state access.State) ([]access.Request, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_sent_access_requests").Add(1)
		ms.latency.With("method", "list_sent_access_requests").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListSentAccessRequests(key, state)
}

func (ms *metricsMiddleware) ListReceivedAccessRequests(key string, state access.State) ([]access.Request, error) {
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

func (ms *metricsMiddleware) ListPotentialPartners(id string) ([]string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_potential_partners").Add(1)
		ms.latency.With("method", "list_potential_partners").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListPotentialPartners(id)
}

func (ms *metricsMiddleware) ApproveAccessRequest(key, id string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "approve_access_request").Add(1)
		ms.latency.With("method", "approve_access_request").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ApproveAccessRequest(key, id)
}

func (ms *metricsMiddleware) RevokeAccessRequest(key, id string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "revoke_access_request").Add(1)
		ms.latency.With("method", "revoke_access_request").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RevokeAccessRequest(key, id)
}

func (ms *metricsMiddleware) GrantAccess(key string, dstUserId string) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "grant_access").Add(1)
		ms.latency.With("method", "grant_access").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.GrantAccess(key, dstUserId)
}
