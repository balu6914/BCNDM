package grpc

import (
	"context"

	"github.com/datapace/datapace/transactions"

	"github.com/go-kit/kit/endpoint"
)

func transferEndpoint(svc transactions.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transferReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.Transfer(req.streamID, req.from, req.to, req.value); err != nil {
			return nil, err
		}

		return transferRes{}, nil
	}
}

func createUserEndpoint(svc transactions.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createUserReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.CreateUser(req.id); err != nil {
			return nil, err
		}

		return createUserRes{}, nil
	}
}
