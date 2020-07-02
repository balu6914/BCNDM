package http

import (
	"context"
	"encoding/json"
	"github.com/datapace/datapace/auth"
	"github.com/go-zoo/bone"
	"net/http"

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
		id, err := svc.Register(req.key, user)
		return createRes{ID: id}, err
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
		r := request.(*http.Request)
		current, err := svc.View(r.Header.Get("Authorization"), bone.GetValue(r, "id"))
		// Filling updateReq with current user data
		updateReq := updateReq{
			ContactEmail: current.ContactEmail,
			FirstName:    current.FirstName,
			LastName:     current.LastName,
			Company:      current.Company,
			Address:      current.Address,
			Phone:        current.Phone,
			Roles:        current.Roles,
			Password:     "",
			Disabled:     current.Disabled,
		}
		if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
			return updateRes{}, err
		}

		if err != nil {
			return updateRes{}, err
		}

		user := auth.User{
			ID:           current.ID,
			Email:        current.Email,
			ContactEmail: updateReq.ContactEmail,
			Password:     updateReq.Password,
			FirstName:    updateReq.FirstName,
			LastName:     updateReq.LastName,
			Company:      updateReq.Company,
			Address:      updateReq.Address,
			Phone:        updateReq.Phone,
			Roles:        updateReq.Roles,
			Disabled:     updateReq.Disabled,
		}
		if err := svc.Update(r.Header.Get("Authorization"), user); err != nil {
			return updateRes{}, err
		}
		return updateRes{}, nil
	}
}

func viewEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(identityReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		user, err := svc.View(req.key, req.ID)
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
				ID:           user.ID,
				FirstName:    user.FirstName,
				LastName:     user.LastName,
				Email:        user.Email,
				ContactEmail: user.ContactEmail,
				Company:      user.Company,
				Address:      user.Address,
				Phone:        user.Phone,
				Disabled:     user.Disabled,
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
