package grpc

import (
	"context"

	"github.com/datapace/datapace/streams"

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
			id:         stream.ID,
			name:       stream.Name,
			owner:      stream.Owner,
			url:        stream.URL,
			price:      stream.Price,
			external:   stream.External,
			project:    stream.BQ.Project,
			dataset:    stream.BQ.Dataset,
			table:      stream.BQ.Table,
			fields:     stream.BQ.Fields,
			terms:      stream.Terms,
			visibility: string(stream.Visibility),
			accessType: string(stream.AccessType),
			maxCalls:   stream.MaxCalls,
			maxUnit:    string(stream.MaxUnit),
		}

		return res, nil
	}
}
