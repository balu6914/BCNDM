package http

import (
	"context"
	"monetasa/streams"

	"github.com/go-kit/kit/endpoint"
	"gopkg.in/mgo.v2/bson"
)

func createStreamEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(createStreamReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		id, err := svc.AddStream(req.owner, req.stream)
		if err != nil {
			return nil, err
		}

		res := createStreamRes{
			ID: id,
		}
		return res, nil
	}
}

func createBulkStreamEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(createBulkStreamReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.AddBulkStream(req.owner, req.Streams); err != nil {
			return nil, err
		}

		return createBulkStreamRes{}, nil
	}
}

func updateStreamEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateStreamReq)

		if req.stream.Owner == "" {
			req.stream.Owner = req.owner
		}

		// Need to set owner before the validation because
		// stream.Validate() won't pass otherwise.
		if err := req.validate(); err != nil {
			return nil, err
		}

		if req.stream.ID == "" {
			req.stream.ID = bson.ObjectIdHex(req.id)
		}

		if err := svc.UpdateStream(req.owner, req.stream); err != nil {
			return nil, err
		}

		return editStreamRes{}, nil
	}
}

func viewStreamEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(readStreamReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		s, err := svc.ViewStream(req.id)
		if err != nil {
			return nil, err
		}

		res := viewStreamRes{
			Stream: s,
		}
		return res, nil
	}
}

func removeStreamEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteStreamReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.RemoveStream(req.owner, req.id); err != nil {
			return nil, err
		}

		return removeStreamRes{}, nil
	}
}

func searchStreamEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(searchStreamReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		q := streams.Query{
			Name:       req.Name,
			StreamType: req.StreamType,
			Coords:     req.Coords,
			Page:       req.Page,
			Limit:      req.Limit,
			MinPrice:   req.MinPrice,
			MaxPrice:   req.MaxPrice,
		}

		page, err := svc.SearchStreams(q)
		if err != nil {
			return nil, err
		}

		res := searchStreamRes{page}
		return res, nil
	}
}
