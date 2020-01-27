package grpc

import (
	"context"
	"datapace/executions"

	"github.com/go-kit/kit/endpoint"
)

func createAlgoEndpoint(svc executions.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(algoReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		algo := executions.Algorithm{
			ID:       req.id,
			Name:     req.name,
			Metadata: req.metadata,
		}

		if err := svc.CreateAlgorithm(algo); err != nil {
			return nil, err
		}

		return createRes{}, nil
	}
}

func createDataEndpoint(svc executions.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dataReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		data := executions.Dataset{
			ID:       req.id,
			Metadata: req.metadata,
		}

		if err := svc.CreateDataset(data); err != nil {
			return nil, err
		}

		return createRes{}, nil
	}
}
