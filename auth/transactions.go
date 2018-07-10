package auth

// TransactionsService defined interface for communication with blockchain service.
type TransactionsService interface {
	// CreateUser creates new user on transaction network (blockchain).
	CreateUser(string) error
}
