// Package transactions provides a transactions service client implementation.
package transactions

import (
	"context"
	"time"

	"github.com/datapace/datapace/auth"
	commonproto "github.com/datapace/datapace/proto/common"
	transactionsproto "github.com/datapace/datapace/proto/transactions"
	"github.com/datapace/datapace/transactions"
)

const timeout = time.Second

var _ auth.TransactionsService = (*transactionsService)(nil)

type transactionsService struct {
	tc transactionsproto.TransactionsServiceClient
}

// NewService returns transactions service implementation.
func NewService(tc transactionsproto.TransactionsServiceClient) auth.TransactionsService {
	return transactionsService{tc: tc}
}

func (ts transactionsService) CreateUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if _, err := ts.tc.CreateUser(ctx, &commonproto.ID{Value: id}); err != nil {
		return transactions.ErrFailedUserCreation
	}

	return nil
}
