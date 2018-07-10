package transactions

import "errors"

var (
	// ErrConflict indicates usage of the existing id during account
	// registration.
	ErrConflict = errors.New("entity already exists")

	// ErrMalformedEntity indicates malformed entity specification (e.g.
	// invalid username or password).
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrFailedUserCreation indicates that user creation was unsuccessful.
	ErrFailedUserCreation = errors.New("failed to create user")

	// ErrFailedBalanceFetch indicates that fetching of user's balance from
	// blockchain failed.
	ErrFailedBalanceFetch = errors.New("failed to fetch user's balance")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")
)

// Service specifies an API that must be fulfilled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// CreateUser receives unique user id, creates user on
	// blockchain and stores its cert.
	CreateUser(string) error

	// Balance receives user unique identifier and channel id and returns its
	// balance read from blockchain.
	Balance(string, string) (uint64, error)

	// TODO: add transfer method
}
