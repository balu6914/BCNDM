package http

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

func updateEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.Update(req.key, req.user); err != nil {
			return nil, err
		}

		return userRes{created: true}, nil
	}
}

func viewEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(identityReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		user, err := svc.View(req.key)
		if err != nil {
			return nil, err
		}

		return viewRes{user}, nil
	}
}

func deleteEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(identityReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.Delete(req.key); err != nil {
			return nil, err
		}

		return removeRes{}, nil
	}
}
