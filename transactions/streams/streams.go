package streams

import (
	"context"
	"monetasa"
	"monetasa/transactions"
	"time"
)

var _ transactions.StreamsService = (*streamsService)(nil)

type streamsService struct {
	client monetasa.StreamsServiceClient
}

// NewService returns instance of streams service client.
func NewService(client monetasa.StreamsServiceClient) transactions.StreamsService {
	return streamsService{client: client}
}

func (ss streamsService) One(id string) (transactions.Stream, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	s, err := ss.client.One(ctx, &monetasa.ID{Value: id})
	if err != nil {
		return transactions.Stream{}, err
	}

	stream := transactions.Stream{
		ID:    s.GetId(),
		Name:  s.GetName(),
		Owner: s.GetOwner(),
	}

	return stream, nil
}
