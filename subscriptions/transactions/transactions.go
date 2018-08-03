package transactions

import (
	"context"
	"monetasa"
	"monetasa/subscriptions"
	"time"
)

var _ subscriptions.TransactionsService = (*transactionsService)(nil)

type transactionsService struct {
	client monetasa.TransactionsServiceClient
}

// NewService returns new transactions service instance.
func NewService(client monetasa.TransactionsServiceClient) subscriptions.TransactionsService {
	return transactionsService{client: client}
}

func (ts transactionsService) Transfer(from, to string, value uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	td := &monetasa.TransferData{
		From:  from,
		To:    to,
		Value: value,
	}
	_, err := ts.client.Transfer(ctx, td)

	return err
}
