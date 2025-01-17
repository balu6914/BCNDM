package http

import (
	"context"
	"time"

	"github.com/datapace/datapace/transactions"

	"github.com/go-kit/kit/endpoint"
)

func buyEndpoint(svc transactions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(buyReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		if req.FundID != "" {
			req.userID = req.FundID
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
		if req.ID != "" {
			req.userID = req.ID
		}
		balance, err := svc.Balance(req.userID)
		if err != nil {
			return nil, err
		}

		return balanceRes{Balance: balance}, nil
	}
}

func txHistoryEndpoint(svc transactions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(txHistoryReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		txHistory, err := svc.TxHistory(req.userID, req.fromDateTime, req.toDateTime, req.txType)
		if err != nil {
			return nil, err
		}

		return txHistory, nil
	}
}

func withdrawEndpoint(svc transactions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(withdrawReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		if req.FundID != "" {
			req.userID = req.FundID
		}
		if err := svc.WithdrawTokens(req.userID, req.Amount); err != nil {
			return nil, err
		}

		return withdrawRes{}, nil
	}
}

func createContractsEndpoint(svc transactions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(createContractsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		now := time.Now()

		contracts := []transactions.Contract{}
		for _, item := range req.Items {
			contracts = append(contracts, transactions.Contract{
				StreamID:    req.StreamID,
				OwnerID:     req.ownerID,
				StartTime:   now,
				EndTime:     req.EndTime,
				PartnerID:   item.PartnerID,
				Share:       item.Share,
				Description: req.Description,
			})
		}
		if err := svc.CreateContracts(contracts...); err != nil {
			return nil, err
		}

		return createContractsRes{}, nil
	}
}

func signContractEndpoint(svc transactions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(signContractReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		contract := transactions.Contract{
			StreamID:  req.StreamID,
			EndTime:   req.EndTime,
			PartnerID: req.partnerID,
		}

		if err := svc.SignContract(contract); err != nil {
			return nil, err
		}

		return signContractRes{}, nil
	}
}

func listContractsEndpoint(svc transactions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(listContractsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		page := svc.ListContracts(req.userID, req.page, req.limit, req.role)
		res := listContractsRes{
			Page:      page.Page,
			Limit:     page.Limit,
			Total:     page.Total,
			Contracts: []contractView{},
		}

		for _, contract := range page.Contracts {
			res.Contracts = append(res.Contracts, contractView{
				StreamID:    contract.StreamID,
				StreamName:  contract.StreamName,
				Description: contract.Description,
				StartTime:   contract.StartTime,
				EndTime:     contract.EndTime,
				OwnerID:     contract.OwnerID,
				PartnerID:   contract.PartnerID,
				Share:       float64(contract.Share) / shareScale,
				Signed:      contract.Signed,
			})
		}

		return res, nil
	}
}
