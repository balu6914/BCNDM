package http

import (
	"context"
	"datapace/auth"

	"github.com/go-kit/kit/endpoint"
)

func registrationEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(registerReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		user := auth.User{
			Email:        req.Email,
			ContactEmail: req.Email,
			Password:     req.Password,
			FirstName:    req.FirstName,
			LastName:     req.LastName,
		}
		err := svc.Register(user)
		return createRes{}, err
	}
}

func loginEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(credentialsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		user := auth.User{
			Email:    req.Email,
			Password: req.Password,
		}
		token, err := svc.Login(user)
		if err != nil {
			return nil, err
		}

		return tokenRes{Token: token}, nil
	}
}

func updateEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		user := auth.User{
			ContactEmail: req.ContactEmail,
			FirstName:    req.FirstName,
			LastName:     req.LastName,
		}
		if err := svc.Update(req.key, user); err != nil {
			return nil, err
		}

		return updateRes{}, nil
	}
}

func updatePasswordEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updatePasswordReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		user := auth.User{
			Password: req.NewPassword,
		}
		if err := svc.UpdatePassword(req.key, req.OldPassword, user); err != nil {
			return nil, err
		}

		return updatePasswordRes{}, nil
	}
}

func viewEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(identityReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		user, err := svc.View(req.key)
		if err != nil {
			return nil, err
		}

		res := viewRes{
			ID:           user.ID,
			Email:        user.Email,
			ContactEmail: user.ContactEmail,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
		}

		return res, nil
	}
}

func listEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(identityReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		users, err := svc.List(req.key)
		if err != nil {
			return nil, err
		}

		res := listRes{
			Users: []viewRes{},
		}
		for _, user := range users {
			res.Users = append(res.Users, viewRes{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			})
		}

		return res, nil
	}
}

func requestAccessEndpoint(svc auth.Service) endpoint.Endpoint {
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

func approveAccessEndpoint(svc auth.Service) endpoint.Endpoint {
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

func rejectAccessEndpoint(svc auth.Service) endpoint.Endpoint {
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

func listSentRequestsEndpoint(svc auth.Service) endpoint.Endpoint {
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

func listReceivedRequestsEndpoint(svc auth.Service) endpoint.Endpoint {
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
