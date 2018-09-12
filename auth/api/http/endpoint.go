package http

import (
	"context"
	"monetasa/auth"

	"github.com/go-kit/kit/endpoint"
)

func registrationEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(registerReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		user := auth.User{
			Email:        req.Email,
			ContactEmail: req.Email,
			Password:     req.Password,
			FirstName:    req.FirstName,
			LastName:     req.LastName,
		}
		err := svc.Register(user)
		return createRes{}, err
	}
}

func loginEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(credentialsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		user := auth.User{
			Email:    req.Email,
			Password: req.Password,
		}
		token, err := svc.Login(user)
		if err != nil {
			return nil, err
		}

		return tokenRes{Token: token}, nil
	}
}

func updateEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		user := auth.User{
			Email:        req.Email,
			Password:     req.Password,
			ContactEmail: req.ContactEmail,
			FirstName:    req.FirstName,
			LastName:     req.LastName,
		}
		if err := svc.Update(req.key, user); err != nil {
			return nil, err
		}

		return updateRes{}, nil
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

		res := viewRes{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}

		return res, nil
	}
}
