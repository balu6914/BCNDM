package http

import (
	"context"
	"monetasa/transactions"

	"github.com/go-kit/kit/endpoint"
)

func buyEndpoint(svc transactions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(buyReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.BuyTokens(req.userID, req.Amount); err != nil {
			return nil, err
		}

		return buyRes{}, nil
	}
}

func balanceEndpoint(svc transactions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(balanceReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		balance, err := svc.Balance(req.userID)
		if err != nil {
			return nil, err
		}

		return balanceRes{Balance: balance}, nil
	}
}

func withdrawEndpoint(svc transactions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(withdrawReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.WithdrawTokens(req.userID, req.Amount); err != nil {
			return nil, err
		}

		return withdrawRes{}, nil
	}
}
