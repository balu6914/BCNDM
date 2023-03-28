package executions

import (
	"context"

	"github.com/datapace/datapace/errors"
	executionsproto "github.com/datapace/datapace/proto/executions"
	"github.com/datapace/datapace/streams"
)

var _ streams.AIService = (*executionsService)(nil)

// Creation errors
var (
	ErrCrateAlgorithm = errors.New("error creating algorithm")
	ErrCrateDataset   = errors.New("error creating dataset")
)

type executionsService struct {
	client executionsproto.ExecutionsServiceClient
}

// New returns ai service interface instance.
func New(client executionsproto.ExecutionsServiceClient) streams.AIService {
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

	algo := executionsproto.Algorithm{
		Id:       s.ID,
		Name:     s.Name,
		Metadata: metadata,
	}

	if _, err := es.client.CreateAlgorithm(context.Background(), &algo); err != nil {
		return errors.Wrap(ErrCrateAlgorithm, err)
	}
	return nil
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

	dataset := executionsproto.Dataset{
		Id:       s.ID,
		Metadata: metadata,
	}

	if _, err := es.client.CreateDataset(context.Background(), &dataset); err != nil {
		return errors.Wrap(ErrCrateDataset, err)
	}
	return nil
}
