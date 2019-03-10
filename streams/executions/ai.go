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
	algo := datapace.Algorithm{
		Id:   s.ID.Hex(),
		Name: s.Metadata["name"].(string),
		Path: s.Metadata["path"].(string),
	}

	if modelToken, ok := s.Metadata["model_token"]; ok {
		if mt, ok := modelToken.(string); ok {
			algo.ModelToken = mt
		}
	}

	if modelName, ok := s.Metadata["model_name"]; ok {
		if mn, ok := modelName.(string); ok {
			algo.ModelName = mn
		}
	}

	_, err := es.client.CreateAlgorithm(context.Background(), &algo)
	return err
}

func (es executionsService) CreateDataset(s streams.Stream) error {
	dataset := datapace.Dataset{
		Id:   s.ID.Hex(),
		Path: s.Metadata["path"].(string),
	}
	_, err := es.client.CreateDataset(context.Background(), &dataset)
	return err
}
