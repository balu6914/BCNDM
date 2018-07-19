package mocks

import (
	"monetasa/transactions"
	"sync"
)

var _ transactions.BlockchainNetwork = (*mockNetwork)(nil)

type mockNetwork struct {
	users map[string]uint64
	mutex *sync.Mutex
}

// NewBlockchainNetwork returns mock instance of blockchain network.
func NewBlockchainNetwork(users map[string]uint64) transactions.BlockchainNetwork {
	return mockNetwork{users: users, mutex: &sync.Mutex{}}
}

func (mn mockNetwork) CreateUser(id, secret string) error {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	if _, ok := mn.users[id]; ok {
		return transactions.ErrFailedUserCreation
	}

	mn.users[id] = 0
	return nil
}

func (mn mockNetwork) Balance(name string) (uint64, error) {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	balance, ok := mn.users[name]
	if !ok {
		return 0, transactions.ErrNotFound
	}

	return balance, nil
}

func (mn mockNetwork) Transfer(from, to string, value uint64) error {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	balance, ok := mn.users[from]
	if !ok {
		return transactions.ErrFailedTransfer
	}

	if balance < value {
		return transactions.ErrFailedTransfer
	}

	mn.users[to] = mn.users[to] + value
	mn.users[from] = mn.users[from] - value

	return nil
}
