package mocks

import (
	"sync"

	"monetasa/subscriptions"
)

var _ subscriptions.SubscriptionsRepository = (*subscriptionsRepositoryMock)(nil)

const subrName = "sub_name"

type subscriptionsRepositoryMock struct {
	mu            sync.Mutex
	subscriptions map[string]subscriptions.Subscription
}

// NewStreamRepository creates in-memory stream repository.
func NewSubscriptionsRepository() subscriptions.SubscriptionsRepository {
	return &subscriptionsRepositoryMock{
		subscriptions: make(map[string]subscriptions.Subscription),
	}
}

func (srm *subscriptionsRepositoryMock) Create(sub subscriptions.Subscription) error {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	if _, ok := srm.subscriptions[sub.UserID]; ok {
		return subscriptions.ErrConflict
	}
	srm.subscriptions[sub.UserID] = sub

	return nil
}

func (srm *subscriptionsRepositoryMock) Read(userId string) ([]subscriptions.Subscription, error) {
	if s, ok := srm.subscriptions[userId]; ok {
		sl := []subscriptions.Subscription{
			s,
		}
		return sl, nil
	}

	return []subscriptions.Subscription{}, subscriptions.ErrNotFound
}
