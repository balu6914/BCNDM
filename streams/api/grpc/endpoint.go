package grpc

import (
	"context"
	"datapace/streams"

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
			id:       stream.ID.Hex(),
			name:     stream.Name,
			owner:    stream.Owner,
			url:      stream.URL,
			price:    stream.Price,
			external: stream.External,
			project:  stream.BQ.Project,
			dataset:  stream.BQ.Dataset,
			table:    stream.BQ.Table,
			fields:   stream.BQ.Fields,
			terms:    stream.Terms,
		}

		return res, nil
	}
}
