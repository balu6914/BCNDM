package http

import (
	"context"
	"github.com/datapace/datapace/terms"
	"github.com/go-kit/kit/endpoint"
)

func validateEndpoint(svc terms.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(validateTermsReq)
		isValid, err := svc.ValidateTerms(terms.Terms{
			StreamID:  req.streamID,
			TermsURL:  req.termsUrl,
			TermsHash: req.termsHash,
		})
		if err != nil {
			return false, err
		}
		return isValid, nil
	}
}
