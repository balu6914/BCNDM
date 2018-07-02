package transactions

var _ Service = (*transactionService)(nil)

type transactionService struct {
	bn BlockchainNetwork
}

// New instantiaces domain service implementation.
func New(bn BlockchainNetwork) Service {
	return transactionService{bn: bn}
}

func (ts transactionService) CreateUser(id, secret string) ([]byte, error) {
	key, err := ts.bn.CreateUser(id, secret)
	if err != nil {
		return []byte{}, ErrFailedUserCreation
	}

	return key, nil
}

func (ts transactionService) Balance(userID, chanID string) (uint64, error) {
	balance, err := ts.bn.Balance(userID, chanID)
	if err != nil {
		return 0, ErrFailedBalanceFetch
	}

	return balance, nil
}
