package api

import (
	"time"

	"monetasa/subscriptions"

	"github.com/go-kit/kit/metrics"
)

var _ subscriptions.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     subscriptions.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc subscriptions.Service, counter metrics.Counter, latency metrics.Histogram) subscriptions.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) CreateSubscription(token string, subs subscriptions.Subscription) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "create_subscription").Add(1)
		ms.latency.With("method", "create_subscription").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CreateSubscription(token, subs)
}

func (ms *metricsMiddleware) ReadSubscriptions(token string) ([]subscriptions.Subscription, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "read_subscriptions").Add(1)
		ms.latency.With("method", "read_subscriptions").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ReadSubscriptions(token)
}
