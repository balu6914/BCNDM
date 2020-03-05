package http

import (
	"context"
	"datapace/dproxy"
	"github.com/go-kit/kit/endpoint"
)

func createTokenEndpoint(svc dproxy.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(requestCreateToken)
		if err := req.validate(); err != nil {
			return nil, err
		}
		token, err := svc.CreateToken(req.URL, req.TTL, req.MaxCalls)
		if err != nil {
			return nil, err
		}
		res := createTokenRes{
			Token: token,
		}
		return res, nil
	}
}
