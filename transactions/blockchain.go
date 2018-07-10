package transactions

// BlockchainNetwork contains blockchain specific API definition.
type BlockchainNetwork interface {
	// CreateUser creates user on blockchain and returns this private key.
	CreateUser(string, string) error

	// Balance returns users account balance for given channel.
	Balance(string, string) (uint64, error)

	// TODO: add transfer method
}
