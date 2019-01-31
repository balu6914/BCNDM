package api

import (
	"datapace/executions"
	"time"

	"github.com/go-kit/kit/metrics"
)

var _ executions.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     executions.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc executions.Service, counter metrics.Counter, latency metrics.Histogram) executions.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (mm metricsMiddleware) Start(owner, algo, data string, mode executions.JobMode) (string, error) {
	defer func(begin time.Time) {
		mm.counter.With("method", "start").Add(1)
		mm.latency.With("method", "start").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.Start(owner, algo, data, mode)
}

func (mm metricsMiddleware) Finish(id string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "finish").Add(1)
		mm.latency.With("method", "finish").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.Finish(id)
}

func (mm metricsMiddleware) Execution(owner, id string) (executions.Execution, error) {
	defer func(begin time.Time) {
		mm.counter.With("method", "execution").Add(1)
		mm.latency.With("method", "execution").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.Execution(owner, id)
}

func (mm metricsMiddleware) List(owner string) ([]executions.Execution, error) {
	defer func(begin time.Time) {
		mm.counter.With("method", "list").Add(1)
		mm.latency.With("method", "list").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.List(owner)
}
