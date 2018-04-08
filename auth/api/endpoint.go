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

func updateEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.Update(req.key, req.id, req.user); err != nil {
			return nil, err
		}

		return userRes{id: req.id, created: false}, nil
	}
}

func viewEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(viewReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		user, err := svc.View(req.key, req.id)
		if err != nil {
			return nil, err
		}

		return viewRes{user}, nil
	}
}

func listEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(listReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		users, err := svc.List(req.key)
		if err != nil {
			return nil, err
		}

		return listRes{users, len(users)}, nil
	}
}

func deleteEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(viewReq)

		err := req.validate()
		if err == auth.ErrNotFound {
			return removeRes{}, nil
		}

		if err != nil {
			return nil, err
		}

		if err = svc.Delete(req.key, req.id); err != nil {
			return nil, err
		}

		return removeRes{}, nil
	}
}

func identityEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(identityReq)

		if err := req.validate(); err != nil {
			return nil, auth.ErrUnauthorizedAccess
		}

		id, err := svc.Identity(req.key)
		if err != nil {
			return nil, err
		}

		return identityRes{id: id}, nil
	}
}

func canAccessEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(viewReq)

		if err := req.validate(); err != nil {
			return nil, auth.ErrUnauthorizedAccess
		}

		id, err := svc.CanAccess(req.key, req.id)
		if err != nil {
			return nil, err
		}

		return identityRes{id: id}, nil
	}
}
