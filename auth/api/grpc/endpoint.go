package grpc

import (
	"context"

	"github.com/datapace/datapace/auth"

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
		req := request.(identityReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		u, err := svc.ViewEmail(req.token)
		if err != nil {
			return emailRes{}, err
		}
		return emailRes{u.Email, u.ContactEmail, nil}, nil
	}
}

func userByEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(byIdReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		u, err := svc.ViewUserById(req.id)
		if err != nil {
			return userResp{err: err}, err
		}
		return userResp{
			id:           u.ID,
			email:        u.Email,
			contactEmail: u.ContactEmail,
			firstName:    u.FirstName,
			lastName:     u.LastName,
			company:      u.Company,
			address:      u.Address,
			country:      u.Country,
			mobile:       u.Mobile,
			phone:        u.Phone,
			role:         u.Role,
			err:          nil,
		}, nil
	}
}

func existsEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(existsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		if err := svc.Exists(req.id); err != nil {
			return existsRes{err: err}, err
		}

		return existsRes{}, nil
	}
}

func authorizeEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(authReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		id, err := svc.Authorize(req.token, auth.Action(req.action), req)
		if err != nil {
			return identityRes{err: err}, err
		}

		return identityRes{id: id}, nil
	}
}
