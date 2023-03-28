package subscriptions

// TransactionsService contains API necessary to transfer tokens.
type TransactionsService interface {
	// Transfers specified amount of tokens from one account to another.
	Transfer(string, string, string, uint64) error
}
