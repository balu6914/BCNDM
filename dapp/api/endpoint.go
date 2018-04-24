package api

import (
	"context"
	"github.com/go-kit/kit/endpoint"

	"monetasa/dapp"
)

func saveStreamEndpoint(svc dapp.StreamRepository) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(createStreamReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.Save(req.Stream)

		if err != nil {
			return nil, err
		}

		return createStreamRes{}, nil
	}
}

func updateStreamEndpoint(svc dapp.StreamRepository) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateStreamReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.Update(req.Id, req.Stream)

		if err != nil {
			return nil, err
		}

		return modifyStreamRes{}, nil
	}
}

func oneStreamEndpoint(svc dapp.StreamRepository) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(readStreamReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		s, err := svc.One(req.Id)

		if err != nil {
			return nil, err
		}

		res := readStreamRes{
			Name:        s.Name,
			Type:        s.Type,
			Description: s.Description,
			Price:       s.Price,
		}
		return res, nil
	}
}

func removeStreamEndpoint(svc dapp.StreamRepository) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteStreamReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.Remove(req.Id)

		if err != nil {
			return nil, err
		}

		return modifyStreamRes{}, nil
	}
}

func searchStreamEndpoint(svc dapp.StreamRepository) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(searchStreamReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		coords := [][]float64{[]float64{req.x0, req.y0}, []float64{req.x1, req.y1},
			[]float64{req.x2, req.y2}, []float64{req.x3, req.y3},
			[]float64{req.x0, req.y0}}

		streams, err := svc.Search(coords)
		if err != nil {
			return nil, err
		}

		res := searchStreamRes{
			Streams: streams,
		}
		return res, nil
	}
}
