package http

import (
	"context"
	"datapace/dproxy"
	"datapace/errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
)

func createTokenEndpoint(svc dproxy.Service, responseType, dProxyRootUrl, fsPrefix, httpPrefix string) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(requestCreateToken)
		if err := req.validate(); err != nil {
			return nil, err
		}
		token, err := svc.CreateToken(req.URL, req.TTL, req.MaxCalls)
		if err != nil {
			return nil, err
		}
		switch responseType {
		case urloutput:
			url := fmt.Sprintf("%s%s/%s", dProxyRootUrl, fsPrefix, token)
			if req.URL[0:4] == "http" {
				url = fmt.Sprintf("%s%s/%s", dProxyRootUrl, httpPrefix, token)
			}
			res := createUrlRes{
				URL: url,
			}
			return res, nil
		case tokenoutput:
			res := createTokenRes{
				Token: token,
			}
			return res, nil
		default:
			return nil, errors.New("unknown output type")
		}
	}
}
