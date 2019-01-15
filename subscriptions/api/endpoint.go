package api

import (
	"context"
	"datapace/subscriptions"

	"github.com/go-kit/kit/endpoint"
)

func addSubEndpoint(svc subscriptions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(addSubReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		id, err := svc.AddSubscription(req.UserID, req.UserToken, req.Subscription)
		if err != nil {
			return nil, err
		}

		return addSubRes{id}, nil
	}
}

func searchSubsEndpoint(svc subscriptions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(searchSubsReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		q := subscriptions.Query{
			Page:        req.Page,
			Limit:       req.Limit,
			StreamID:    req.StreamID,
			StreamOwner: req.StreamOwner,
			UserID:      req.UserID,
		}

		page, err := svc.SearchSubscriptions(q)
		if err != nil {
			return nil, err
		}

		res := searchSubsRes{page}
		return res, nil
	}
}

func viewSubEndpoint(svc subscriptions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(viewSubReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		s, err := svc.ViewSubscription(req.userID, req.subscriptionID)
		if err != nil {
			return nil, err
		}

		res := viewSubRes{
			Subscription: s,
		}
		return res, nil
	}
}
