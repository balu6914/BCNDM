package mocks

import (
	"sync"

	"gopkg.in/mgo.v2/bson"

	"github.com/datapace/datapace/auth"
)

var _ auth.UserRepository = (*userRepositoryMock)(nil)

type userRepositoryMock struct {
	mu       sync.Mutex
	users    map[string]auth.User
	policyMu *sync.Mutex
	policies map[string]auth.Policy
}

// NewUserRepository creates in-memory user repository.
func NewUserRepository(hasher auth.Hasher, usr auth.User, policies map[string]auth.Policy, mu *sync.Mutex) auth.UserRepository {
	users := make(map[string]auth.User)
	hash, _ := hasher.Hash(usr.Password)
	usr.Password = hash
	urm := &userRepositoryMock{users: users, policies: policies, policyMu: mu}
	urm.Save(usr)
	return urm
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
	if user.ID == "" {
		user.ID = bson.NewObjectId().Hex()
	}
	urm.users[user.Email], urm.users[user.ID] = user, user
	return user.ID, nil
}

func (urm *userRepositoryMock) OneByID(id string) (auth.User, error) {
	urm.mu.Lock()
	defer urm.mu.Unlock()

	u, ok := urm.users[id]
	if !ok {
		return auth.User{}, auth.ErrNotFound
	}
	urm.policyMu.Lock()
	defer urm.policyMu.Unlock()

	for _, p := range urm.policies {
		if p.Owner == u.ID {
			u.Policies = append(u.Policies, p)
		}
	}

	return u, nil
}

func (urm *userRepositoryMock) OneByEmail(email string) (auth.User, error) {
	urm.mu.Lock()
	defer urm.mu.Unlock()

	u, ok := urm.users[email]
	if !ok {
		return auth.User{}, auth.ErrNotFound
	}

	urm.policyMu.Lock()
	defer urm.policyMu.Unlock()

	for _, p := range urm.policies {
		if p.Owner == u.ID {
			u.Policies = append(u.Policies, p)
		}
	}

	return u, nil
}

func (urm *userRepositoryMock) Update(user auth.User) error {
	urm.mu.Lock()
	defer urm.mu.Unlock()

	u, ok := urm.users[user.ID]

	if !ok {
		return auth.ErrNotFound
	}
	if user.Email != "" {
		u.Email = user.Email
	}
	if user.ContactEmail != "" {
		u.ContactEmail = user.ContactEmail
	}
	if user.Password != "" {
		u.Password = user.Password
	}
	if user.FirstName != "" {
		u.FirstName = user.FirstName
	}
	if user.LastName != "" {
		u.LastName = user.LastName
	}
	if user.Company != "" {
		u.Company = user.Company
	}
	if user.Address != "" {
		u.Address = user.Address
	}
	if user.Country != "" {
		u.Country = user.Country
	}
	if user.Mobile != "" {
		u.Mobile = user.Mobile
	}
	if user.Phone != "" {
		u.Phone = user.Phone
	}
	u.Disabled = user.Disabled
	if len(user.Policies) != 0 {
		u.Policies = user.Policies
	}

	urm.users[user.Email] = u
	urm.users[urm.users[user.Email].ID] = u

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

func (urm *userRepositoryMock) RetrieveAll(auth.AdminFilters) ([]auth.User, error) {
	return nil, nil
}
