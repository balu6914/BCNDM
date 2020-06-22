package http

import (
	"context"

	"github.com/datapace/auth"

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
			Company:      req.Company,
			Address:      req.Address,
			Phone:        req.Phone,
			Roles:        req.Roles,
		}
		err := svc.Register(req.key, user)
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
			Company:      req.Company,
			Address:      req.Address,
			Phone:        req.Phone,
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
			Company:      user.Company,
			Address:      user.Address,
			Phone:        user.Phone,
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

		users, err := svc.ListUsers(req.key)
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

func nonPartnersEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(identityReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		users, err := svc.ListNonPartners(req.key)
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
