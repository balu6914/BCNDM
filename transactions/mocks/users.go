package mocks

import "monetasa/transactions"

var _ transactions.UserRepository = (*mockUserRepository)(nil)

type mockUserRepository struct {
	users map[string]string
}

// NewUserRepository returns mock implementation of user repository.
func NewUserRepository(users map[string]string) transactions.UserRepository {
	return mockUserRepository{users: users}
}

func (repo mockUserRepository) Save(user transactions.User) error {
	if _, ok := repo.users[user.ID]; ok {
		return transactions.ErrConflict
	}

	repo.users[user.ID] = user.Secret
	return nil
}

func (repo mockUserRepository) Remove(id string) error {
	if _, ok := repo.users[id]; !ok {
		return transactions.ErrNotFound
	}

	delete(repo.users, id)
	return nil
}
