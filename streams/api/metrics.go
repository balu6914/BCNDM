package api

import (
	"monetasa/streams"
	"time"

	"github.com/go-kit/kit/metrics"
)

var _ streams.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     streams.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc streams.Service, counter metrics.Counter, latency metrics.Histogram) streams.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) AddStream(key string, stream streams.Stream) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "add_stream").Add(1)
		ms.latency.With("method", "add_stream").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.AddStream(key, stream)
}

func (ms *metricsMiddleware) AddBulkStream(key string, streams []streams.Stream) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "add_bulk_stream").Add(1)
		ms.latency.With("method", "add_bulk_stream").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.AddBulkStream(key, streams)
}

func (ms *metricsMiddleware) UpdateStream(key string, stream streams.Stream) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_stream").Add(1)
		ms.latency.With("method", "update_stream").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateStream(key, stream)
}

func (ms *metricsMiddleware) ViewStream(id string) (streams.Stream, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_stream").Add(1)
		ms.latency.With("method", "view_stream").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewStream(id)
}

func (ms *metricsMiddleware) SearchStreams(coords [][]float64) ([]streams.Stream, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "search_streams").Add(1)
		ms.latency.With("method", "search_streams").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.SearchStreams(coords)
}

func (ms *metricsMiddleware) RemoveStream(key, id string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "remove_stream").Add(1)
		ms.latency.With("method", "remove_stream").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RemoveStream(key, id)
}
