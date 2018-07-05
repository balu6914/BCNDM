package subscriptions

import (
	"errors"
	"time"
)

var (
	// ErrConflict indicates usage of the existing stream id for the new stream.
	ErrConflict error = errors.New("Subscription ID already taken")

	// ErrMalformedEntity indicates malformed entity specification (e.g.
	// invalid username or password).
	ErrMalformedEntity error = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess error = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound error = errors.New("non-existent entity")

	ErrUnsupportedContentType error = errors.New("unsupported content type")
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
}

// New instantiates the domain service implementation.
func New(subs SubscriptionsRepository) Service {
	return &subscriptionsService{
		subscriptions: subs,
	}
}

func (ss *subscriptionsService) CreateSubscription(userID string, sub Subscription) error {
	sub.UserID = userID
	sub.StartDate = time.Now()
	sub.EndDate = time.Now().Add(time.Hour * time.Duration(sub.Hours))

	// TODO: Verify User Balance and transfer tokens

	return ss.subscriptions.Create(sub)
}

func (ss *subscriptionsService) ReadSubscriptions(userID string) ([]Subscription, error) {
	return ss.subscriptions.Read(userID)
}
