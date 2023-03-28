package grpc

import (
	"context"
	"github.com/datapace/datapace/terms"
	"github.com/go-kit/kit/endpoint"
)

func createTermsEndpoint(svc terms.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(termsReq)

		terms := terms.Terms{
			StreamID: req.streamId,
			TermsURL: req.termsUrl,
		}

		if _, err := svc.CreateTerms(terms); err != nil {
			return nil, err
		}

		return createRes{}, nil
	}
}
