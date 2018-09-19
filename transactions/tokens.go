package transactions

// TokenLedger contains token specific API definition.
type TokenLedger interface {
	// CreateUser creates user on blockchain and returns this private key.
	CreateUser(string, string) error

	// Balance returns users account balance.
	Balance(string) (uint64, error)

	// Transfers tokens from one account to another.
	Transfer(string, string, string, uint64) error

	// BuyTokens transfers tokens to given account.
	BuyTokens(string, uint64) error

	// WithdrawTokens transfers tokens from account to coin base.
	WithdrawTokens(string, uint64) error
}
