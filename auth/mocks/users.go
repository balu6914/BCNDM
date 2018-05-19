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

	// Save user with two keys. This will allow to simulate
	// mongo queries by email or _id (login and register or other endpoints)
	if _, ok := urm.users[user.Email]; ok {
		return auth.ErrConflict
	}
	if _, ok := urm.users[user.ID.Hex()]; ok {
		return auth.ErrConflict
	}

	urm.users[user.Email], urm.users[user.ID.Hex()] = user, user
	return nil
}

func (urm *userRepositoryMock) OneById(id string) (auth.User, error) {
	urm.mu.Lock()
	defer urm.mu.Unlock()

	if u, ok := urm.users[id]; ok {
		return u, nil
	}

	return auth.User{}, auth.ErrNotFound
}

func (urm *userRepositoryMock) OneByEmail(email string) (auth.User, error) {
	urm.mu.Lock()
	defer urm.mu.Unlock()

	if u, ok := urm.users[email]; ok {
		return u, nil
	}

	return auth.User{}, auth.ErrNotFound
}

func (urm *userRepositoryMock) Update(user auth.User) error {
	urm.mu.Lock()
	defer urm.mu.Unlock()

	urm.users[user.Email] = user
	urm.users[urm.users[user.Email].ID.Hex()] = user

	return nil
}

func (urm *userRepositoryMock) Remove(id string) error {
	urm.mu.Lock()
	defer urm.mu.Unlock()

	if _, ok := urm.users[id]; !ok {
		return auth.ErrNotFound
	}
	if _, ok := urm.users[urm.users[id].Email]; !ok {
		return auth.ErrNotFound
	}

	urm.users[urm.users[id].Email] = auth.User{}
	urm.users[id] = auth.User{}
	return nil
}
