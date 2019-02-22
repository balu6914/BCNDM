package http

import (
	"context"

	access "datapace/access-control"

	"github.com/go-kit/kit/endpoint"
)

func requestAccessEndpoint(svc access.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(requestAccessReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		id, err := svc.RequestAccess(req.key, req.Receiver)
		if err != nil {
			return nil, err
		}

		res := requestAccessRes{
			id: id,
		}

		return res, nil
	}
}

func approveAccessEndpoint(svc access.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(approveAccessReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.ApproveAccessRequest(req.key, req.id); err != nil {
			return nil, err
		}

		res := approveAccessRes{}
		return res, nil
	}
}

func rejectAccessEndpoint(svc access.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(rejectAccessReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.RejectAccessRequest(req.key, req.id); err != nil {
			return nil, err
		}

		res := rejectAccessRes{}
		return res, nil
	}
}

func listSentRequestsEndpoint(svc access.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(listAccessRequestsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		requests, err := svc.ListSentAccessRequests(req.key, req.state)
		if err != nil {
			return nil, err
		}

		res := listAccessRequestsRes{
			Requests: []viewAccessRequestRes{},
		}
		for _, r := range requests {
			res.Requests = append(res.Requests, viewAccessRequestRes{
				ID:       r.ID,
				Sender:   r.Sender,
				Receiver: r.Receiver,
				State:    r.State,
			})
		}

		return res, nil
	}
}

func listReceivedRequestsEndpoint(svc access.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(listAccessRequestsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		requests, err := svc.ListReceivedAccessRequests(req.key, req.state)
		if err != nil {
			return nil, err
		}

		res := listAccessRequestsRes{
			Requests: []viewAccessRequestRes{},
		}
		for _, r := range requests {
			res.Requests = append(res.Requests, viewAccessRequestRes{
				ID:       r.ID,
				Sender:   r.Sender,
				Receiver: r.Receiver,
				State:    r.State,
			})
		}

		return res, nil
	}
}
