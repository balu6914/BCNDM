package api

import (
	"context"
	"datapace/subscriptions"

	"github.com/go-kit/kit/endpoint"
)

func addSubEndpoint(svc subscriptions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {

		req := request.(addSubsReq)

		resps := addSubsRes{}

		for _, sub := range req.Subscriptions {

			if err := sub.Validate(); err != nil {
				resps.Responses = append(resps.Responses, addSubResp{
					StreamID:     sub.StreamID,
					ErrorMessage: err.Error(),
				})
				continue
			}

			id, err := svc.AddSubscription(req.UserID, req.UserToken, sub)
			if err != nil {
				resps.Responses = append(resps.Responses, addSubResp{
					StreamID:     sub.StreamID,
					ErrorMessage: err.Error(),
				})
				continue
			}

			resps.Responses = append(resps.Responses, addSubResp{
				StreamID:       sub.StreamID,
				SubscriptionID: id,
			})
		}
		
		return resps, nil
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
