package api

import (
	"github.com/datapace/datapace/terms"
	"time"

	"github.com/go-kit/kit/metrics"
)

var _ terms.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     terms.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc terms.Service, counter metrics.Counter, latency metrics.Histogram) terms.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (mm metricsMiddleware) CreateTerms(terms terms.Terms) (string, error) {
	defer func(begin time.Time) {
		mm.counter.With("method", "create_terms").Add(1)
		mm.latency.With("method", "create_terms").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.CreateTerms(terms)
}

func (mm metricsMiddleware) ValidateTerms(terms terms.Terms) (bool, error) {
	defer func(begin time.Time) {
		mm.counter.With("method", "validate_terms").Add(1)
		mm.latency.With("method", "validate_terms").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.ValidateTerms(terms)
}
