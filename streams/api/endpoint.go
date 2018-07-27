package api

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"gopkg.in/mgo.v2/bson"

	"monetasa/streams"
)

func addStreamEndpoint(svc streams.Service) endpoint.Endpoint {
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

func addBulkStreamEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(createBulkStreamReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.AddBulkStream(req.owner, req.Streams); err != nil {
			return nil, err
		}

		return createBulkStreamResponse{}, nil
	}
}

func updateStreamEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateStreamReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		req.stream.ID = bson.ObjectIdHex(req.streamID)

		if err := svc.UpdateStream(req.owner, req.stream); err != nil {
			return nil, err
		}

		return modifyStreamRes{}, nil
	}
}

func viewStreamEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(readStreamReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		s, err := svc.ViewStream(req.streamID)
		if err != nil {
			return nil, err
		}

		res := readStreamRes{
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

		if err := svc.RemoveStream(req.owner, req.streamID); err != nil {
			return nil, err
		}

		return modifyStreamRes{}, nil
	}
}

func searchStreamEndpoint(svc streams.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(searchStreamReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		streams, err := svc.SearchStreams(req.points)
		if err != nil {
			return nil, err
		}

		res := searchStreamRes{
			Streams: streams,
		}
		return res, nil
	}
}
