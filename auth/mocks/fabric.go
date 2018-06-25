package mocks

import (
	"monetasa/auth"
	"sync"
)

type fabricMock struct {
	mu    sync.Mutex
	users map[string]string
}

// NewFabricNetwork creates in-memory user repository that simulates blockchain.
func NewFabricNetwork() auth.FabricNetwork {
	return &fabricMock{
		users: make(map[string]string),
	}
}

func (fm *fabricMock) Initialize() error {
	return nil
}

func (fm *fabricMock) CreateUser(id, secret string) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if _, ok := fm.users[id]; ok {
		return auth.ErrConflict
	}

	fm.users[id] = secret
	return nil
}

func (fm *fabricMock) Balance(id string) (uint64, error) {
	return 0, nil
}
