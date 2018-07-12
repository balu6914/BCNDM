// Package transactions provides a transactions service client implementation.
package transactions

import (
	"context"
	"monetasa"
	"monetasa/auth"
	"monetasa/transactions"
	"time"
)

const timeout = time.Second

var _ auth.TransactionsService = (*transactionsService)(nil)

type transactionsService struct {
	tc monetasa.TransactionsServiceClient
}

// NewService returns transactions service implementation.
func NewService(tc monetasa.TransactionsServiceClient) auth.TransactionsService {
	return transactionsService{tc: tc}
}

func (ts transactionsService) CreateUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if _, err := ts.tc.CreateUser(ctx, &monetasa.ID{Value: id}); err != nil {
		return transactions.ErrFailedUserCreation
	}

	return nil
}