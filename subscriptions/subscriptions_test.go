package subscriptions_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"monetasa/subscriptions"
	"monetasa/subscriptions/mocks"
)

const (
	userID    = "myUserID"
	userID1   = "myUserID1"
	streamID  = "myStreamID"
	wrong     = "wrong"
	streamURL = "myUrl"
	token     = "token"
	email     = "user@example.com"
	balance   = 100
)

var subscription = subscriptions.Subscription{
	UserID:    userID,
	StreamID:  streamID,
	StreamURL: streamURL,
	Hours:     5,
}

func newService(tokens map[string]string) subscriptions.Service {
	subs := mocks.NewSubscriptionsRepository()
	streams := mocks.NewStreamsService(map[string]subscriptions.Stream{
		streamID: subscriptions.Stream{Price: 10},
	})
	proxy := mocks.NewProxy()
	transactions := mocks.NewTransactionsService(balance)

	return subscriptions.New(subs, streams, proxy, transactions)
}

func TestCreateSubscription(t *testing.T) {
	svc := newService(map[string]string{token: userID})

	cases := []struct {
		desc string
		sub  subscriptions.Subscription
		err  error
	}{
		{
			desc: "create new subscription",
			sub:  subscription,
			err:  nil,
		},
		{
			desc: "create existing subscription",
			sub:  subscription,
			err:  subscriptions.ErrConflict,
		},
		{
			desc: "create subscription with non-existent stream",
			sub: subscriptions.Subscription{
				UserID:    userID,
				StreamID:  wrong,
				StreamURL: streamURL,
			},
			err: subscriptions.ErrNotFound,
		},
		{
			desc: "create subscription with too big price stream",
			sub: subscriptions.Subscription{
				UserID:    userID1,
				StreamID:  streamID,
				StreamURL: streamURL,
				Hours:     100,
			},
			err: subscriptions.ErrFailedTransfer,
		},
	}

	for _, tc := range cases {
		err := svc.CreateSubscription(tc.sub.UserID, tc.sub)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestGetSubscriptions(t *testing.T) {
	svc := newService(map[string]string{token: userID})

	svc.CreateSubscription(userID, subscription)

	cases := map[string]struct {
		userID string
		err    error
	}{
		"get subscription with valid token": {
			userID,
			nil,
		},
		"read subscriptions with non-existent entity": {
			wrong,
			subscriptions.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		_, err := svc.ReadSubscriptions(tc.userID)
		assert.Equal(t, tc.err, err,
			fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}
