package mocks

import (
	"sync"

	"monetasa/dapp"
)

var _ dapp.StreamRepository = (*streamRepositoryMock)(nil)

const subrName = "sub_name"

type subscriptionsRepositoryMock struct {
	mu            sync.Mutex
	subscriptions map[string]dapp.Subscription
}

// NewStreamRepository creates in-memory stream repository.
func NewSubscriptionsRepository() dapp.SubscriptionsRepository {
	return &subscriptionsRepositoryMock{
		subscriptions: make(map[string]dapp.Subscription),
	}
}

func (srm *subscriptionsRepositoryMock) Create(sub dapp.Subscription) error {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	srm.subscriptions[sub.UserID] = sub

	return nil
}

func (srm *subscriptionsRepositoryMock) Update(id string, stream dapp.Subscription) error {
	// TODO: implement a geolocation search mock

	return nil
}

func (srm *subscriptionsRepositoryMock) Read(id string) ([]dapp.Subscription, error) {
	// TODO: implement a geolocation search mock

	return []dapp.Subscription{}, dapp.ErrNotFound
}
