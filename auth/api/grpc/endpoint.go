package grpc

import (
	"context"
	"datapace/auth"

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
		u, err := svc.View(req.token)
		if err != nil {
			return emailRes{}, err
		}
		return emailRes{u.Email, u.ContactEmail, nil}, nil
	}
}

func partnersEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(partnersReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		partners, err := svc.ListPartners(req.id)
		if err != nil {
			return partnersRes{}, err
		}

		res := partnersRes{
			partners: partners,
		}

		return res, nil
	}
}
