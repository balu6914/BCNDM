package api

import (
	"time"

	"github.com/datapace/datapace/subscriptions"

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

func (ms *metricsMiddleware) AddSubscription(id, token string, subs subscriptions.Subscription) (subscriptions.Subscription, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "add_subscription").Add(1)
		ms.latency.With("method", "add_subscription").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.AddSubscription(id, token, subs)
}

func (ms *metricsMiddleware) ViewSubscription(userID, subID string) (subscriptions.Subscription, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_subscription").Add(1)
		ms.latency.With("method", "view_subscription").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewSubscription(userID, subID)
}

func (ms *metricsMiddleware) SearchSubscriptions(query subscriptions.Query) (subscriptions.Page, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "search_subscriptions").Add(1)
		ms.latency.With("method", "search_subscriptions").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.SearchSubscriptions(query)
}

func (ms *metricsMiddleware) ViewSubByUserAndStream(userID, streamID string) (subscriptions.Subscription, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_sub_by_user_and_stream").Add(1)
		ms.latency.With("method", "view_sub_by_user_and_stream").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewSubByUserAndStream(userID, streamID)
}

func (ms *metricsMiddleware) Report(query subscriptions.Query, owner string) ([]byte, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "report").Add(1)
		ms.latency.With("method", "report").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.Report(query, owner)
}
