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
	streamID  = "myStreamID"
	wrong     = "wrong"
	streamURL = "myUrl"
	token     = "token"
	email     = "user@example.com"
)

var subscription = subscriptions.Subscription{
	UserID:    userID,
	StreamID:  streamID,
	StreamURL: streamURL,
}

func newService(tokens map[string]string) subscriptions.Service {
	subs := mocks.NewSubscriptionsRepository()
	return subscriptions.New(subs)
}

func TestCreateSubscription(t *testing.T) {
	svc := newService(map[string]string{token: userID})

	cases := map[string]struct {
		sub subscriptions.Subscription
		err error
	}{
		"create new subsription": {
			subscription,
			nil,
		},
	}

	for desc, tc := range cases {
		err := svc.CreateSubscription(userID, tc.sub)
		assert.Equal(t, tc.err, err,
			fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}

func TestGetSubscriptions(t *testing.T) {
	svc := newService(map[string]string{token: userID})

	svc.CreateSubscription(userID, subscription)

	cases := map[string]struct {
		userID string
		err    error
	}{
		"read subsription with valid token": {
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
