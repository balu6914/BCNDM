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
	// Subscribe to a stream
	CreateSubscription(string, Subscription) error

	// Get all user subscriptions
	ReadSubscriptions(string) ([]Subscription, error)
}

var _ Service = (*subscriptionsService)(nil)

type subscriptionsService struct {
	subscriptions SubscriptionsRepository
	streams       StreamsService
	proxy         Proxy
	transactions  TransactionsService
}

// New instantiates the domain service implementation.
func New(subs SubscriptionsRepository, streams StreamsService, proxy Proxy, transactions TransactionsService) Service {
	return &subscriptionsService{
		subscriptions: subs,
		streams:       streams,
		proxy:         proxy,
		transactions:  transactions,
	}
}

func (ss *subscriptionsService) CreateSubscription(userID string, sub Subscription) error {
	sub.UserID = userID
	sub.StartDate = time.Now()
	sub.EndDate = time.Now().Add(time.Hour * time.Duration(sub.Hours))

	stream, err := ss.streams.One(sub.StreamID)
	if err != nil {
		return ErrNotFound
	}

	hash, err := ss.proxy.Register(sub.Hours, stream.URL)
	if err != nil {
		return err
	}

	sub.StreamURL = hash

	if err := ss.subscriptions.Create(sub); err != nil {
		if err == ErrConflict {
			return ErrConflict
		}

		return ErrFailedCreateSub
	}

	if err := ss.transactions.Transfer(userID, stream.Owner, stream.Price*sub.Hours); err != nil {
		if err == ErrNotEnoughTokens {
			return err
		}
		return ErrFailedTransfer
	}

	return nil
}

func (ss *subscriptionsService) ReadSubscriptions(userID string) ([]Subscription, error) {
	return ss.subscriptions.Read(userID)
}
