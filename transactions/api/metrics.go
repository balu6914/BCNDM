package api

import (
	"monetasa/transactions"
	"time"

	"github.com/go-kit/kit/metrics"
)

var _ transactions.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     transactions.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc transactions.Service, counter metrics.Counter, latency metrics.Histogram) transactions.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (mm *metricsMiddleware) CreateUser(id, password string) ([]byte, error) {
	defer func(begin time.Time) {
		mm.counter.With("method", "create_user").Add(1)
		mm.latency.With("method", "create_user").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.CreateUser(id, password)
}

func (mm *metricsMiddleware) Balance(userID, chanID string) (uint64, error) {
	defer func(begin time.Time) {
		mm.counter.With("method", "balance").Add(1)
		mm.latency.With("method", "balance").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.Balance(userID, chanID)
}
