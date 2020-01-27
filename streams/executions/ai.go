package executions

import (
	"context"
	"datapace"
	"datapace/streams"
)

var _ streams.AIService = (*executionsService)(nil)

type executionsService struct {
	client datapace.ExecutionsServiceClient
}

// New returns ai service interface instance.
func New(client datapace.ExecutionsServiceClient) streams.AIService {
	return executionsService{
		client: client,
	}
}

func (es executionsService) CreateAlgorithm(s streams.Stream) error {
	metadata := map[string]string{}
	for k, v := range s.Metadata {
		val, ok := v.(string)
		if !ok {
			return streams.ErrMalformedData
		}
		metadata[k] = val
	}

	algo := datapace.Algorithm{
		Id:       s.ID.Hex(),
		Name:     s.Name,
		Metadata: metadata,
	}

	_, err := es.client.CreateAlgorithm(context.Background(), &algo)
	return err
}

func (es executionsService) CreateDataset(s streams.Stream) error {
	metadata := map[string]string{}
	for k, v := range s.Metadata {
		val, ok := v.(string)
		if !ok {
			return streams.ErrMalformedData
		}
		metadata[k] = val
	}

	dataset := datapace.Dataset{
		Id:       s.ID.Hex(),
		Metadata: metadata,
	}

	_, err := es.client.CreateDataset(context.Background(), &dataset)
	return err
}
