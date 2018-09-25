package grpc

import (
	"context"
	"monetasa/auth"

	"github.com/go-kit/kit/endpoint"
)

func identifyEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(identityReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		id, err := svc.Identify(req.token)
		if err != nil {
			return identityRes{}, err
		}
		return identityRes{id, nil}, nil
	}
}

func emailEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(emailReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		u, err := svc.View(req.id)
		if err != nil {
			return emailRes{}, err
		}
		return emailRes{u.Email, nil}, nil
	}
}
