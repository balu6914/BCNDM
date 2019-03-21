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

func (mm metricsMiddleware) Start(exec executions.Execution) (string, error) {
	defer func(begin time.Time) {
		mm.counter.With("method", "start").Add(1)
		mm.latency.With("method", "start").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.Start(exec)
}

func (mm metricsMiddleware) Result(owner, id string) (map[string]interface{}, error) {
	defer func(begin time.Time) {
		mm.counter.With("method", "result").Add(1)
		mm.latency.With("method", "result").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.Result(owner, id)
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

func (mm metricsMiddleware) CreateAlgorithm(algo executions.Algorithm) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "create_algorithm").Add(1)
		mm.latency.With("method", "create_algorithm").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.CreateAlgorithm(algo)
}

func (mm metricsMiddleware) CreateDataset(data executions.Dataset) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "create_dataset").Add(1)
		mm.latency.With("method", "create_dataset").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.CreateDataset(data)
}

func (mm metricsMiddleware) ProcessEvents() error {
	defer func(begin time.Time) {
		mm.counter.With("method", "process_events").Add(1)
		mm.latency.With("method", "process_events").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.ProcessEvents()
}
