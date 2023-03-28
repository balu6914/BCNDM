package api

import (
	"context"

	"github.com/datapace/datapace/subscriptions"

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

			s, err := svc.AddSubscription(req.UserID, req.UserToken, sub)
			if err != nil {
				resps.Responses = append(resps.Responses, addSubResp{
					StreamID:     sub.StreamID,
					ErrorMessage: err.Error(),
				})
				continue
			}

			resps.Responses = append(resps.Responses, addSubResp{
				StreamID:       sub.StreamID,
				SubscriptionID: s.ID.Hex(),
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

		page, err := svc.SearchSubscriptions(req.Query)
		if err != nil {
			return nil, err
		}

		res := searchSubsRes{page}
		return res, nil
	}
}

func reportSubsEndpoint(svc subscriptions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(searchSubsReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		resp, err := svc.Report(req.Query, req.owner)
		if err != nil {
			return nil, err
		}

		res := reportResponse(resp)
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

func viewSubByUserAndStreamEndpoint(svc subscriptions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(viewSubByUserAndStreamReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		s, err := svc.ViewSubByUserAndStream(req.userID, req.streamID)
		if err != nil {
			return nil, err
		}

		res := viewSubRes{
			Subscription: s,
		}
		return res, nil
	}
}
