package api

import (
	"datapace/transactions"
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

func (mm *metricsMiddleware) CreateUser(id string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "create_user").Add(1)
		mm.latency.With("method", "create_user").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.CreateUser(id)
}

func (mm *metricsMiddleware) Balance(userID string) (uint64, error) {
	defer func(begin time.Time) {
		mm.counter.With("method", "balance").Add(1)
		mm.latency.With("method", "balance").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.Balance(userID)
}

func (mm *metricsMiddleware) Transfer(stream, from, to string, value uint64) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "transfer").Add(1)
		mm.latency.With("method", "transfer").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.Transfer(stream, from, to, value)
}

func (mm *metricsMiddleware) BuyTokens(account string, value uint64) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "buy_tokens").Add(1)
		mm.latency.With("method", "buy_tokens").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.BuyTokens(account, value)
}

func (mm *metricsMiddleware) WithdrawTokens(account string, value uint64) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "withdraw_tokens").Add(1)
		mm.latency.With("method", "withdraw_tokens").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.WithdrawTokens(account, value)
}

func (mm *metricsMiddleware) CreateContracts(contracts ...transactions.Contract) (err error) {
	defer func(begin time.Time) {
		mm.counter.With("method", "create_contracts").Add(1)
		mm.latency.With("method", "create_contracts").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.CreateContracts(contracts...)
}

func (mm *metricsMiddleware) SignContract(contract transactions.Contract) (err error) {
	defer func(begin time.Time) {
		mm.counter.With("method", "sign_contract").Add(1)
		mm.latency.With("method", "sign_contract").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.SignContract(contract)
}

func (mm *metricsMiddleware) ListContracts(userID string, pageNo uint64, limit uint64, role transactions.Role) transactions.ContractPage {
	defer func(begin time.Time) {
		mm.counter.With("method", "list_contracts").Add(1)
		mm.latency.With("method", "list_contracts").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.ListContracts(userID, pageNo, limit, role)
}
