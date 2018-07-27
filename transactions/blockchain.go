package transactions

// BlockchainNetwork contains blockchain specific API definition.
type BlockchainNetwork interface {
	// CreateUser creates user on blockchain and returns this private key.
	CreateUser(string, string) error

	// Balance returns users account balance.
	Balance(string) (uint64, error)

	// Transfers tokens from one account to another.
	Transfer(string, string, uint64) error

	// BuyTokens transfers tokens to given account.
	BuyTokens(string, uint64) error
}
