package subscriptions_test

import (
	"fmt"
	"github.com/datapace/datapace/subscriptions/sharing"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"

	"github.com/datapace/datapace/subscriptions"
	"github.com/datapace/datapace/subscriptions/mocks"
)

const (
	user1ID   = "user1ID"
	user2ID   = "user2ID"
	noUser    = "noUser"
	stream1ID = "myStream1ID"
	stream2ID = "myStream2ID"
	wrong     = "wrong"
	streamURL = "myUrl"
	token     = "token"
	email     = "user@example.com"
	balance   = 100
)

var subscription = subscriptions.Subscription{
	ID:          bson.NewObjectId(),
	UserID:      user1ID,
	StreamOwner: user2ID,
	StreamID:    stream1ID,
	StreamURL:   streamURL,
	Hours:       5,
}

func newService(tokens map[string]string) subscriptions.Service {
	subs := mocks.NewSubscriptionsRepository()
	streams := mocks.NewStreamsService(map[string]subscriptions.Stream{
		stream1ID: {Price: 10, Owner: user2ID},
		stream2ID: {Price: 100, Owner: user1ID},
	})
	proxy := mocks.NewProxy()
	transactions := mocks.NewTransactionsService(balance)
	auth := mocks.NewAuthClient(tokens, nil)
	sharingSvc := sharing.NewServiceMock()
	return subscriptions.New(auth, subs, streams, proxy, transactions, sharingSvc)
}

func TestAddSubscription(t *testing.T) {
	svc := newService(map[string]string{token: user1ID})

	cases := []struct {
		desc string
		sub  subscriptions.Subscription
		err  error
	}{
		{
			desc: "create a new subscription",
			sub:  subscription,
			err:  nil,
		},
		{
			desc: "create an existing subscription",
			sub:  subscription,
			err:  subscriptions.ErrConflict,
		},
		{
			desc: "create a subscription with non-existent stream",
			sub: subscriptions.Subscription{
				ID:        bson.NewObjectId(),
				UserID:    user1ID,
				StreamID:  wrong,
				StreamURL: streamURL,
			},
			err: subscriptions.ErrNotFound,
		},
		{
			desc: "create a subscription with too big price stream",
			sub: subscriptions.Subscription{
				ID:        bson.NewObjectId(),
				UserID:    user2ID,
				StreamID:  stream1ID,
				StreamURL: streamURL,
				Hours:     100,
			},
			err: subscriptions.ErrFailedTransfer,
		},
	}

	for _, tc := range cases {
		_, err := svc.AddSubscription(tc.sub.UserID, "", tc.sub)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestSearchSubscriptions(t *testing.T) {
	svc := newService(map[string]string{token: user1ID})

	s := subscription
	defLimit := uint64(3)
	total := uint64(20)
	for i := 0; i < 20; i++ {
		s.ID = bson.NewObjectId()
		svc.AddSubscription(user1ID, "", s)
	}

	s.ID = bson.NewObjectId()
	s.StreamID = stream2ID
	svc.AddSubscription(user2ID, "", s)

	cases := []struct {
		desc  string
		query subscriptions.Query
		page  subscriptions.Page
		size  int
	}{
		{
			desc: "get subscription for the specific user",
			query: subscriptions.Query{
				UserID: user1ID,
				Limit:  defLimit,
			},
			page: subscriptions.Page{
				Total: total,
				Page:  0,
				Limit: defLimit,
			},
			size: 3,
		},
		{
			desc: "get subscription with user and page specified",
			query: subscriptions.Query{
				UserID: user1ID,
				Page:   3,
				Limit:  defLimit,
			},
			page: subscriptions.Page{
				Total: total,
				Page:  3,
				Limit: defLimit,
			},
			size: 3,
		},
		{
			desc: "get subscription for the user with no subscriptions",
			query: subscriptions.Query{
				StreamOwner: noUser,
				Limit:       3,
			},
			page: subscriptions.Page{
				Total: 0,
				Page:  0,
				Limit: defLimit,
			},
			size: 0,
		},
		{
			desc: "get subscription for the specific owner",
			query: subscriptions.Query{
				StreamOwner: user1ID,
				Limit:       defLimit,
			},
			page: subscriptions.Page{
				Total: 1,
				Page:  0,
				Limit: defLimit,
			},
			size: 1,
		},
		{
			desc: "get subscription for the owner with no subscriptions on his streams",
			query: subscriptions.Query{
				StreamOwner: noUser,
				Limit:       3,
			},
			page: subscriptions.Page{
				Total: 0,
				Page:  0,
				Limit: defLimit,
			},
			size: 0,
		},
	}

	for _, tc := range cases {
		ret, _ := svc.SearchSubscriptions(tc.query)
		assert.Equal(t, tc.page.Limit, ret.Limit, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.page.Limit, ret.Limit))
		assert.Equal(t, tc.page.Total, ret.Total, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.page.Total, ret.Total))
		assert.Equal(t, tc.page.Page, ret.Page, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.page.Page, ret.Page))
		assert.Equal(t, tc.size, len(ret.Content), fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, len(ret.Content)))
	}
}

func TestViewSubscription(t *testing.T) {
	svc := newService(map[string]string{token: user1ID})

	_, err := svc.AddSubscription(subscription.UserID, "", subscription)
	require.Nil(t, err, "Saving Subscription expected to succeed.")

	cases := []struct {
		desc           string
		subscriptionID string
		userID         string
		err            error
	}{
		{
			desc:           "get a subscription by the user",
			subscriptionID: subscription.ID.Hex(),
			userID:         subscription.UserID,
			err:            nil,
		},
		{
			desc:           "get a subscription by the owner",
			subscriptionID: subscription.ID.Hex(),
			userID:         user2ID,
			err:            nil,
		},
		{
			desc:           "get a wrong subscription",
			subscriptionID: subscription.ID.Hex(),
			userID:         noUser,
			err:            subscriptions.ErrNotFound,
		},
	}

	for _, tc := range cases {
		_, err := svc.ViewSubscription(tc.userID, tc.subscriptionID)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}
