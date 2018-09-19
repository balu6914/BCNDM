package mocks

import (
	"monetasa/transactions"
	"sync"
)

var _ transactions.TokenLedger = (*mockNetwork)(nil)

type mockNetwork struct {
	users     map[string]uint64
	remaining uint64
	mutex     *sync.Mutex
}

// NewTokenLedger returns mock instance of blockchain network.
func NewTokenLedger(users map[string]uint64, remaining uint64) transactions.TokenLedger {
	return &mockNetwork{
		users:     users,
		remaining: remaining,
		mutex:     &sync.Mutex{},
	}
}

func (mn *mockNetwork) CreateUser(id, secret string) error {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	if _, ok := mn.users[id]; ok {
		return transactions.ErrFailedUserCreation
	}

	mn.users[id] = 0
	return nil
}

func (mn *mockNetwork) Balance(name string) (uint64, error) {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	balance, ok := mn.users[name]
	if !ok {
		return 0, transactions.ErrNotFound
	}

	return balance, nil
}

func (mn *mockNetwork) Transfer(stream, from, to string, value uint64) error {
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

func (mn *mockNetwork) BuyTokens(account string, value uint64) error {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	if value > mn.remaining {
		return transactions.ErrFailedTransfer
	}

	mn.users[account] += value
	mn.remaining -= value

	return nil
}

func (mn *mockNetwork) WithdrawTokens(account string, value uint64) error {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	if value > mn.users[account] {
		return transactions.ErrFailedTransfer
	}

	mn.users[account] -= value
	mn.remaining += value

	return nil
}
