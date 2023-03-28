package grpc

import (
	"context"

	"github.com/datapace/datapace/dproxy"
	"github.com/go-kit/kit/endpoint"
)

func listEndpoint(svc dproxy.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listRequest)
		if err := req.validate(); err != nil {
			return nil, err
		}

		page, err := svc.List(req.query)
		if err != nil {
			return nil, err
		}

		res := listResponse{
			events: page,
		}

		return res, nil
	}
}
