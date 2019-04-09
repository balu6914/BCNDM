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
	iname, ok := s.Metadata["name"]
	if !ok {
		return streams.ErrMalformedData
	}

	name, ok := iname.(string)
	if !ok {
		return streams.ErrMalformedData
	}

	ipath, ok := s.Metadata["path"]
	if !ok {
		return streams.ErrMalformedData
	}

	path, ok := ipath.(string)
	if !ok {
		return streams.ErrMalformedData
	}

	modelToken := ""
	if imodelToken, ok := s.Metadata["model_token"]; ok {
		if mt, ok := imodelToken.(string); ok {
			modelToken = mt
		}
	}

	modelName := ""
	if imodelName, ok := s.Metadata["model_name"]; ok {
		if mn, ok := imodelName.(string); ok {
			modelName = mn
		}
	}

	algo := datapace.Algorithm{
		Id:         s.ID.Hex(),
		Name:       name,
		Path:       path,
		ModelToken: modelToken,
		ModelName:  modelName,
	}

	_, err := es.client.CreateAlgorithm(context.Background(), &algo)
	return err
}

func (es executionsService) CreateDataset(s streams.Stream) error {
	ipath, ok := s.Metadata["path"]
	if !ok {
		return streams.ErrMalformedData
	}

	path, ok := ipath.(string)
	if !ok {
		return streams.ErrMalformedData
	}

	dataset := datapace.Dataset{
		Id:   s.ID.Hex(),
		Path: path,
	}

	_, err := es.client.CreateDataset(context.Background(), &dataset)
	return err
}
