package streams

import (
	"context"
	"time"

	commonproto "github.com/datapace/datapace/proto/common"
	streamsproto "github.com/datapace/datapace/proto/streams"
	"github.com/datapace/datapace/subscriptions"
)

var _ subscriptions.StreamsService = (*streamsService)(nil)

type streamsService struct {
	client streamsproto.StreamsServiceClient
}

// NewService returns instance of streams service client.
func NewService(client streamsproto.StreamsServiceClient) subscriptions.StreamsService {
	return streamsService{client: client}
}

func (ss streamsService) One(id string) (subscriptions.Stream, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	s, err := ss.client.One(ctx, &commonproto.ID{Value: id})
	if err != nil {
		return subscriptions.Stream{}, err
	}

	parsedEndDate, _ := time.Parse(time.RFC3339, s.GetEndDate())
	stream := subscriptions.Stream{
		ID:         s.GetId(),
		Name:       s.GetName(),
		Owner:      s.GetOwner(),
		URL:        s.GetUrl(),
		Price:      s.GetPrice(),
		External:   s.GetExternal(),
		Offer:      s.GetOffer(),
		Project:    s.GetProject(),
		Dataset:    s.GetDataset(),
		Table:      s.GetTable(),
		Fields:     s.GetFields(),
		Visibility: s.GetVisibility(),
		AccessType: s.GetAccessType(),
		MaxCalls:   s.GetMaxCalls(),
		MaxUnit:    s.GetMaxUnit(),
		EndDate:    &parsedEndDate,
	}

	return stream, nil
}
