package mocks

import (
	"sync"

	"monetasa/subscriptions"
)

var _ subscriptions.SubscriptionRepository = (*subscriptionRepositoryMock)(nil)

type subscriptionRepositoryMock struct {
	mu            sync.Mutex
	subscriptions map[string]subscriptions.Subscription
}

func equal(queryID, id string) bool {
	if queryID == "" {
		return true
	}
	return queryID == id
}

// NewSubscriptionsRepository creates in-memory stream repository.
func NewSubscriptionsRepository() subscriptions.SubscriptionRepository {
	return &subscriptionRepositoryMock{
		subscriptions: make(map[string]subscriptions.Subscription),
	}
}

func (srm *subscriptionRepositoryMock) Save(sub subscriptions.Subscription) (string, error) {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	if _, ok := srm.subscriptions[sub.ID.Hex()]; ok {
		return "", subscriptions.ErrConflict
	}

	srm.subscriptions[sub.ID.Hex()] = sub

	return sub.ID.Hex(), nil
}

func (srm *subscriptionRepositoryMock) One(subID string) (subscriptions.Subscription, error) {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	if sub, ok := srm.subscriptions[subID]; ok {
		return sub, nil
	}
	return subscriptions.Subscription{}, subscriptions.ErrNotFound
}

func (srm *subscriptionRepositoryMock) Search(query subscriptions.Query) (subscriptions.Page, error) {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	ret := []subscriptions.Subscription{}
	for _, v := range srm.subscriptions {
		if equal(query.StreamID, v.StreamID) &&
			equal(query.UserID, v.UserID) &&
			equal(query.StreamOwner, v.StreamOwner) {
			ret = append(ret, v)
		}
	}

	start := query.Page * query.Limit
	end := start + query.Limit
	page := subscriptions.Page{
		Total:   uint64(len(ret)),
		Limit:   query.Limit,
		Page:    query.Page,
		Content: []subscriptions.Subscription{},
	}

	n := uint64(len(ret))
	if start >= n {
		return page, nil
	}
	if end >= n {
		end = n
	}
	ret = ret[start:end]
	page.Content = ret

	return page, nil
}
