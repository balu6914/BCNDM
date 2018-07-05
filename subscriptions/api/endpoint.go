package api

import (
	"context"
	"monetasa/subscriptions"

	"github.com/go-kit/kit/endpoint"
)

func createSubsEndpoint(svc subscriptions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(subscriptionReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.CreateSubscription(req.UserID, req.Subscription); err != nil {
			return nil, err
		}

		return subscriptionRes{}, nil
	}
}

func readSubsEndpoint(svc subscriptions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(listSubsReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		subs, err := svc.ReadSubscriptions(req.UserID)
		if err != nil {
			return nil, err
		}

		res := listSubsRes{
			Subscriptions: subs,
		}

		return res, nil
	}
}
