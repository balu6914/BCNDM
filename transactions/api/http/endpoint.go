package http

import (
	"context"
	"monetasa/transactions"

	"github.com/go-kit/kit/endpoint"
)

func balanceEndpoint(svc transactions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(balanceReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		balance, err := svc.Balance(req.userID, req.chanID)
		if err != nil {
			return nil, err
		}

		return balanceRes{Balance: balance}, err
	}
}
