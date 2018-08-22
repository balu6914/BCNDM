package subscriptions

import (
	"errors"
	"time"
)

var (
	// ErrConflict indicates usage of the existing stream id for the new stream.
	ErrConflict = errors.New("Subscription ID already taken")

	// ErrMalformedEntity indicates malformed entity specification (e.g.
	// invalid username or password).
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrFailedTransfer indicates that token transfer failed.
	ErrFailedTransfer = errors.New("failed to transfer tokens")

	// ErrFailedCreateSub indicates that creation of subscription failed.
	ErrFailedCreateSub = errors.New("failed to create subscription")

	// ErrNotEnoughTokens indicates that spender doesn't have enough tokens.
	ErrNotEnoughTokens = errors.New("not enough tokens")
)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// AddSubscription subscribes to a stream making a new Subscription.
	AddSubscription(string, Subscription) (string, error)

	// SearchSubscriptions searches subscriptions by the query.
	SearchSubscriptions(Query) (Page, error)

	// ViewSubscription retrieves subscription by the given ID.
	ViewSubscription(string, string) (Subscription, error)
}

var _ Service = (*subscriptionsService)(nil)

type subscriptionsService struct {
	subscriptions SubscriptionRepository
	streams       StreamsService
	proxy         Proxy
	transactions  TransactionsService
}

// New instantiates the domain service implementation.
func New(subs SubscriptionRepository, streams StreamsService, proxy Proxy, transactions TransactionsService) Service {
	return &subscriptionsService{
		subscriptions: subs,
		streams:       streams,
		proxy:         proxy,
		transactions:  transactions,
	}
}

func (ss subscriptionsService) AddSubscription(userID string, sub Subscription) (string, error) {
	sub.UserID = userID
	sub.StartDate = time.Now()
	sub.EndDate = time.Now().Add(time.Hour * time.Duration(sub.Hours))

	stream, err := ss.streams.One(sub.StreamID)
	if err != nil {
		return "", ErrNotFound
	}
	sub.StreamOwner = stream.Owner

	hash, err := ss.proxy.Register(sub.Hours, stream.URL)
	if err != nil {
		return "", err
	}

	sub.StreamURL = hash

	id, err := ss.subscriptions.Save(sub)
	if err != nil {
		if err == ErrConflict {
			return "", ErrConflict
		}

		return "", ErrFailedCreateSub
	}

	if err := ss.transactions.Transfer(userID, stream.Owner, stream.Price*sub.Hours); err != nil {
		if err == ErrNotEnoughTokens {
			return "", err
		}
		return "", ErrFailedTransfer
	}

	return id, nil
}

func (ss subscriptionsService) SearchSubscriptions(q Query) (Page, error) {
	return ss.subscriptions.Search(q)
}

func (ss subscriptionsService) ViewSubscription(userID, subID string) (Subscription, error) {
	sub, err := ss.subscriptions.One(subID)
	if err != nil {
		return Subscription{}, err
	}

	if userID != sub.StreamOwner && userID != sub.UserID {
		return Subscription{}, ErrNotFound
	}

	return sub, nil
}
