package mocks

import "github.com/datapace/datapace/auth"

var _ auth.TransactionsService = (*mockTransactionsService)(nil)

type mockTransactionsService struct{}

// NewTransactionsService returns mock instance of transactions service.
func NewTransactionsService() auth.TransactionsService {
	return mockTransactionsService{}
}

func (ts mockTransactionsService) CreateUser(id string) error {
	if id == "" {
		return auth.ErrMalformedEntity
	}

	return nil
}
