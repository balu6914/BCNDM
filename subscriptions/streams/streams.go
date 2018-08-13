package streams

import (
	"context"
	"monetasa"
	"monetasa/subscriptions"
	"time"
)

var _ subscriptions.StreamsService = (*streamsService)(nil)

type streamsService struct {
	client monetasa.StreamsServiceClient
}

// NewService returns instance of streams service client.
func NewService(client monetasa.StreamsServiceClient) subscriptions.StreamsService {
	return streamsService{client: client}
}

func (ss streamsService) One(id string) (subscriptions.Stream, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	s, err := ss.client.One(ctx, &monetasa.ID{Value: id})
	if err != nil {
		return subscriptions.Stream{}, err
	}

	stream := subscriptions.Stream{
		ID:    s.GetId(),
		Owner: s.GetOwner(),
		Price: s.GetPrice(),
	}

	return stream, nil
}
