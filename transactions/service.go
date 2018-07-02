package transactions

import "errors"

var (
	// ErrFailedUserCreation indicates that user creation was unsuccessful.
	ErrFailedUserCreation = errors.New("failed to create user")

	// ErrFailedBalanceFetch indicates that fetching of user's balance from
	// blockchain failed.
	ErrFailedBalanceFetch = errors.New("failed to fetch user's balance")
)

// Service specifies an API that must be fulfilled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// CreateUser receives unique user id and secret, creates user on
	// blockchain and returns generated private key.
	CreateUser(string, string) ([]byte, error)

	// Balance receives user unique identifier and channel id and returns its
	// balance read from blockchain.
	Balance(string, string) (uint64, error)

	// TODO: add transfer method
}
