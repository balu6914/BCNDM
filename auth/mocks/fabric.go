package mocks

import (
	"monetasa/auth"
	"sync"
)

type fabricMock struct {
	mu    sync.Mutex
	users map[string]string
}

// NewUserRepository creates in-memory user repository.
func NewFabricNetwork() auth.FabricNetwork {
	return &fabricMock{
		users: make(map[string]string),
	}
}

func (fm *fabricMock) Initialize() error {
	return nil
}

func (fm *fabricMock) CreateUser(user *auth.User) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if _, ok := fm.users[user.ID.Hex()]; ok {
		return auth.ErrConflict
	}

	fm.users[user.ID.Hex()] = user.Password
	return nil
}

func (fm *fabricMock) Balance(id string) (uint64, error) {
	return 0, nil
}
