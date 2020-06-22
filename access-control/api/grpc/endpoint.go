package grpc

import (
	"context"

	access "github.com/datapace/access-control"

	"github.com/go-kit/kit/endpoint"
)

func partnersEndpoint(svc access.Service) endpoint.Endpoint {
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

func potentialEndpoint(svc access.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(partnersReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		partners, err := svc.ListPotentialPartners(req.id)
		if err != nil {
			return partnersRes{}, err
		}

		res := partnersRes{
			partners: partners,
		}

		return res, nil
	}
}
