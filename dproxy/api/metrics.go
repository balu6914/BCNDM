package api

import (
	"datapace/dproxy"
	"github.com/go-kit/kit/metrics"
	"time"
)

var _ dproxy.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     dproxy.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc dproxy.Service, counter metrics.Counter, latency metrics.Histogram) dproxy.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) CreateToken(url string, ttl, maxCalls int) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "create_token").Add(1)
		ms.latency.With("method", "create_token").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CreateToken(url, ttl, maxCalls)
}

func (ms *metricsMiddleware) GetTargetURL(url string) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "get_target_url").Add(1)
		ms.latency.With("method", "get_target_url").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.GetTargetURL(url)
}
