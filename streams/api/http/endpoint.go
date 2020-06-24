package http

import (
	"context"

	"github.com/datapace/datapace/streams"

	"github.com/go-kit/kit/endpoint"
)

func addStreamEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(addStreamReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		id, err := svc.AddStream(req.stream)
		if err != nil {
			return nil, err
		}

		res := addStreamRes{
			ID: id,
		}
		return res, nil
	}
}

func addBulkStreamsEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(addBulkStreamsReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.AddBulkStreams(req.Streams); err != nil {
			switch err.(type) {
			case streams.ErrBulkConflict:
				return conflictBulkStreamsRes{err.(streams.ErrBulkConflict)}, nil
			default:
				return nil, err
			}
		}

		return addBulkStreamsRes{}, nil
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
			req.stream.ID = req.id
		}

		if err := svc.UpdateStream(req.stream); err != nil {
			return nil, err
		}

		return updateStreamRes{}, nil
	}
}

func viewStreamEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(viewStreamReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		s, err := svc.ViewStream(req.id, req.owner)
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
		req := request.(removeStreamReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.RemoveStream(req.owner, req.id); err != nil {
			return nil, err
		}

		return removeStreamRes{}, nil
	}
}

func searchStreamsEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(searchStreamsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		q := streams.Query{
			Name:       req.Name,
			Owner:      req.Owner,
			StreamType: req.StreamType,
			Coords:     req.Coords,
			Page:       req.Page,
			Limit:      req.Limit,
			MinPrice:   req.MinPrice,
			MaxPrice:   req.MaxPrice,
		}

		page, err := svc.SearchStreams(req.user, q)
		if err != nil {
			return nil, err
		}

		res := searchStreamsRes{page}
		return res, nil
	}
}
