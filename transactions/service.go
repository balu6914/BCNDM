package transactions

import "errors"

const (
	// Owner role enumeration.
	Owner Role = 1 << iota
	// Partner role enumeration.
	Partner
	// AllRoles enumeration.
	AllRoles
)

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

	// ErrFailedTransfer indicates that token transfer failed.
	ErrFailedTransfer = errors.New("failed to transfer tokens")

	// ErrNotEnoughTokens indicates that spender doesn't have enough tokens.
	ErrNotEnoughTokens = errors.New("not enough tokens")

	ErrFailedTxHistoryFetch = errors.New("failed to fetch tx history")

	ErrFailedSerialization = errors.New("failed while object serialization")
)

// Role represents enumeration for user roles.
type Role int

// Service specifies an API that must be fulfilled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// CreateUser receives unique user id, creates user on
	// blockchain and stores its cert.
	CreateUser(string) error

	// Balance receives user unique identifier and returns its balance read from
	// blockchain.
	Balance(string) (uint64, error)

	// Transfer receives from and to ids and amount of tokens that it should
	// transfer. It returns error only if transfer failed.
	Transfer(string, string, string, uint64) error

	// BuyTokens transfers tokens to user's account.
	BuyTokens(string, uint64) error

	// WithdrawTokens exchanges tokens for real money.
	WithdrawTokens(string, uint64) error

	// TxHistory gives transaction history for a user.
	TxHistory(string, string, string, string) (TokenTxHistory, error)

	// CreateContracts creates multiple contracts at once.
	CreateContracts(...Contract) error

	// SignContract signs existing contract.
	SignContract(Contract) error

	// ListContracts finds and returns page of contracts by contract
	// owner or partner, depending on passed role.
	ListContracts(string, uint64, uint64, Role) ContractPage
}
