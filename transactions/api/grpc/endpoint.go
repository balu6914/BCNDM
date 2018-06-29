package grpc

import (
	"context"
	"monetasa/transactions"

	"github.com/go-kit/kit/endpoint"
)

func createUserEndpoint(svc transactions.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createUserReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		privateKey, err := svc.CreateUser(req.id, req.secret)
		if err != nil {
			return createUserRes{}, err
		}
		return createUserRes{privateKey, nil}, nil
	}
}
