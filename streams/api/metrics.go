package api

import (
	"time"

	"github.com/datapace/datapace/streams"

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

func (ms *metricsMiddleware) AddStream(stream streams.Stream) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "add_stream").Add(1)
		ms.latency.With("method", "add_stream").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.AddStream(stream)
}

func (ms *metricsMiddleware) AddBulkStreams(streams []streams.Stream) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "add_bulk_stream").Add(1)
		ms.latency.With("method", "add_bulk_stream").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.AddBulkStreams(streams)
}

func (ms *metricsMiddleware) UpdateStream(stream streams.Stream) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_stream").Add(1)
		ms.latency.With("method", "update_stream").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateStream(stream)
}

func (ms *metricsMiddleware) ViewFullStream(id string) (streams.Stream, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_full_stream").Add(1)
		ms.latency.With("method", "view_full_stream").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewFullStream(id)
}

func (ms *metricsMiddleware) ViewStream(id, owner string) (streams.Stream, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_stream").Add(1)
		ms.latency.With("method", "view_stream").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewStream(id, owner)
}

func (ms *metricsMiddleware) SearchStreams(owner string, query streams.Query) (streams.Page, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "search_streams").Add(1)
		ms.latency.With("method", "search_streams").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.SearchStreams(owner, query)
}

func (ms *metricsMiddleware) RemoveStream(key, id string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "remove_stream").Add(1)
		ms.latency.With("method", "remove_stream").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RemoveStream(key, id)
}

func (ms *metricsMiddleware) ExportStreams(owner string) ([]streams.Stream, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "export_streams").Add(1)
		ms.latency.With("method", "export_streams").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.ExportStreams(owner)
}
