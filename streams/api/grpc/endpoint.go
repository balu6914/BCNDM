package grpc

import (
	"context"
	"monetasa/streams"

	"github.com/go-kit/kit/endpoint"
)

func oneEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(oneReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		stream, err := svc.ViewFullStream(req.id)
		if err != nil {
			return nil, err
		}

		res := oneRes{
			id:    stream.ID.Hex(),
			owner: stream.Owner,
			url:   stream.URL,
			price: stream.Price,
		}

		return res, nil
	}
}
