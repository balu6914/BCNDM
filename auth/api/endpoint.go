package api

import (
	"context"

	"monetasa/auth"
	"github.com/go-kit/kit/endpoint"
)

func registrationEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(userReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.Register(req.user)
		return tokenRes{}, err
	}
}

func loginEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(userReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		token, err := svc.Login(req.user)
		if err != nil {
			return nil, err
		}

		return tokenRes{token}, nil
	}
}

func addClientEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(addClientReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		id, err := svc.AddClient(req.key, req.client)
		if err != nil {
			return nil, err
		}

		return clientRes{id: id, created: true}, nil
	}
}
