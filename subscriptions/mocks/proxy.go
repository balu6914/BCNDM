package mocks

import "github.com/datapace/datapace/subscriptions"

var _ subscriptions.Proxy = (*mockProxy)(nil)

type mockProxy struct{}

// NewProxy returns mock proxy instance.
func NewProxy() subscriptions.Proxy {
	return mockProxy{}
}

func (mp mockProxy) Register(ttl uint64, url string, maxCalls uint64, maxUnit, subscriptionID string) (string, error) {
	return "", nil
}
