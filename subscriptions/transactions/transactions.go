package transactions

import (
	"context"
	"datapace"
	"datapace/subscriptions"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ subscriptions.TransactionsService = (*transactionsService)(nil)

type transactionsService struct {
	client datapace.TransactionsServiceClient
}

// NewService returns new transactions service instance.
func NewService(client datapace.TransactionsServiceClient) subscriptions.TransactionsService {
	return transactionsService{client: client}
}

func (ts transactionsService) Transfer(streamID, from, to string, value uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	td := &datapace.TransferData{
		StreamID: streamID,
		From:     from,
		To:       to,
		Value:    value,
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
