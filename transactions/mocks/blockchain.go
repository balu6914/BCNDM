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

func (mn mockNetwork) CreateUser(id, secret string) ([]byte, error) {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	if _, ok := mn.users[id]; ok {
		return []byte{}, transactions.ErrFailedUserCreation
	}

	mn.users[id] = 0
	return []byte(secret), nil
}

func (mn mockNetwork) Balance(name, _ string) (uint64, error) {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	balance, ok := mn.users[name]
	if !ok {
		return 0, transactions.ErrFailedBalanceFetch
	}

	return balance, nil
}