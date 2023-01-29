package transactions

import (
	"context"
	"time"

	transactionsproto "github.com/datapace/datapace/proto/transactions"
	"github.com/datapace/datapace/subscriptions"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	dateTimeFormat = "02-01-2006 15:04:05"
)

var _ subscriptions.TransactionsService = (*transactionsService)(nil)

type transactionsService struct {
	client transactionsproto.TransactionsServiceClient
}

// NewService returns new transactions service instance.
func NewService(client transactionsproto.TransactionsServiceClient) subscriptions.TransactionsService {
	return transactionsService{client: client}
}

func (ts transactionsService) Transfer(streamID, from, to string, value uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	td := &transactionsproto.TransferData{
		StreamID: streamID,
		From:     from,
		To:       to,
		Value:    value,
		DateTime: time.Now().Format(dateTimeFormat),
	}
	if _, err := ts.client.Transfer(ctx, td); err != nil {
		e, ok := status.FromError(err)
		if ok && e.Code() == codes.FailedPrecondition {
			return subscriptions.ErrNotEnoughTokens
		}

		return err
	}

	return nil
}
