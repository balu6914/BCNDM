package mocks

import (
	"sync"

	"github.com/datapace/datapace/auth"
)

var _ auth.UserRepository = (*userRepositoryMock)(nil)

type userRepositoryMock struct {
	mu    sync.Mutex
	users map[string]auth.User
}

// NewUserRepository creates in-memory user repository.
func NewUserRepository(hasher auth.Hasher, usr auth.User) auth.UserRepository {
	users := make(map[string]auth.User)
	hash, _ := hasher.Hash(usr.Password)
	users[usr.Email] = auth.User{
		ID:       usr.ID,
		Email:    usr.Email,
		Password: hash,
		Roles:    usr.Roles,
	}
	return &userRepositoryMock{
		users: users,
	}
}

func (urm *userRepositoryMock) Save(user auth.User) (string, error) {
	urm.mu.Lock()
	defer urm.mu.Unlock()

	// Save user with two keys. This will allow to simulate
	// mongo queries by email or _id (login and register or other endpoints)
	if _, ok := urm.users[user.Email]; ok {
		return "", auth.ErrConflict
	}
	if _, ok := urm.users[user.ID]; ok {
		return "", auth.ErrConflict
	}

	urm.users[user.Email], urm.users[user.ID] = user, user
	return user.ID, nil
}

func (urm *userRepositoryMock) OneByID(id string) (auth.User, error) {
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

	if _, ok := urm.users[user.ID]; !ok {
		return auth.ErrNotFound
	}

	urm.users[user.Email] = user
	urm.users[urm.users[user.Email].ID] = user

	return nil
}

func (urm *userRepositoryMock) Remove(id string) error {
	urm.mu.Lock()
	defer urm.mu.Unlock()

	if _, ok := urm.users[id]; !ok {
		return auth.ErrNotFound
	}

	delete(urm.users, id)
	return nil
}

func (urm *userRepositoryMock) AllExcept([]string) ([]auth.User, error) {
	return nil, nil
}
