package mocks

import (
	"monetasa/auth"
	"sync"
)

type userRepositoryMock struct {
	mu    sync.Mutex
	users map[string]auth.User
}

// NewUserRepository creates in-memory user repository.
func NewUserRepository() auth.UserRepository {
	return &userRepositoryMock{
		users: make(map[string]auth.User),
	}
}

func (urm *userRepositoryMock) Save(user auth.User) error {
	urm.mu.Lock()
	defer urm.mu.Unlock()

	if _, ok := urm.users[user.Email]; ok {
		return auth.ErrConflict
	}

	urm.users[user.Email] = user
	return nil
}

func (urm *userRepositoryMock) One(email string) (auth.User, error) {
	urm.mu.Lock()
	defer urm.mu.Unlock()

	if val, ok := urm.users[email]; ok {
		return val, nil
	}

	return auth.User{}, auth.ErrUnauthorizedAccess
}

func (ur *userRepositoryMock) Update(user auth.User) error {
	return nil
}

func (ur *userRepositoryMock) Remove(email string) error {
	return nil
}
