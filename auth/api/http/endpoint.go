package http

import (
	"context"

	"github.com/datapace/datapace/auth"

	"github.com/go-kit/kit/endpoint"
)

func registerEndpoint(svc auth.Service) endpoint.Endpoint {
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
			Role:         req.Role,
			Metadata:     req.Metadata,
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

func updateUserEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		if err := svc.UpdateUser(req.key, req.toUser()); err != nil {
			return nil, err
		}
		return okRes{}, nil
	}
}

func recoverPasswordEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(recoverReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		if err := svc.RecoverPassword(req.Email); err != nil {
			return nil, err
		}
		return okRes{}, nil
	}
}

func validateRecoveryTokenEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(recoveryTokenReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		if err := svc.ValidateRecoveryToken(req.token, req.id); err != nil {
			return nil, err
		}
		return okRes{}, nil
	}
}

func updatePasswordEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(recoveryPasswordReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		if err := svc.UpdatePassword(req.token, req.id, req.Password); err != nil {
			return nil, err
		}
		return okRes{}, nil
	}
}

func viewUserEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(identityReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		user, err := svc.ViewUser(req.key, req.ID)
		if err != nil {
			return nil, err
		}
		var policies []interface{}
		for _, p := range user.Policies {
			policies = append(policies, p)
		}

		res := viewUserRes{
			ID:           user.ID,
			Email:        user.Email,
			ContactEmail: user.ContactEmail,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Company:      user.Company,
			Address:      user.Address,
			Phone:        user.Phone,
			Role:         user.Role,
			CreatedDate:  user.CreatedDate,
			Policies:     policies,
			Metadata:     user.Metadata,
		}

		return res, nil
	}
}

func listUsersEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(identityReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		users, err := svc.ListUsers(req.key)
		if err != nil {
			return nil, err
		}

		res := listUsersRes{
			Users: []viewUserRes{},
		}
		for _, user := range users {
			res.Users = append(res.Users, viewUserRes{
				ID:           user.ID,
				FirstName:    user.FirstName,
				LastName:     user.LastName,
				Email:        user.Email,
				ContactEmail: user.ContactEmail,
				Company:      user.Company,
				Address:      user.Address,
				Phone:        user.Phone,
				Disabled:     user.Disabled,
				Locked:       user.Locked,
				Role:         user.Role,
				CreatedDate:  user.CreatedDate,
				Metadata:     user.Metadata,
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

		res := listUsersRes{
			Users: []viewUserRes{},
		}
		for _, user := range users {
			res.Users = append(res.Users, viewUserRes{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			})
		}

		return res, nil
	}
}

func addPolicyEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(policyRequest)
		if err := req.validate(); err != nil {
			return nil, err
		}

		policy := auth.Policy{
			Version: req.Version,
			Name:    req.Name,
			Rules:   []auth.Rule{},
		}
		for _, r := range req.Rules {
			rule := auth.Rule{
				Type:   r.Type,
				Action: r.Action,
			}
			if r.Condition != nil {
				rule.Condition = auth.SimpleCondition{
					Key:   r.Condition.Key,
					Value: r.Condition.Value,
				}
			}
			policy.Rules = append(policy.Rules, rule)
		}

		id, err := svc.AddPolicy(req.key, policy)
		if err != nil {
			return nil, err
		}

		return createRes{ID: id}, nil
	}
}

func viewPolicyEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(identityReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		policy, err := svc.ViewPolicy(req.key, req.ID)
		if err != nil {
			return nil, err
		}
		res := viewPolicyRes{
			ID:      policy.ID,
			Version: policy.Version,
			Owner:   policy.Owner,
			Name:    policy.Name,
			Rules:   []rule{},
		}
		for _, r := range policy.Rules {
			rule := rule{
				Action: r.Action,
				Type:   r.Type,
			}
			if c, ok := r.Condition.(auth.SimpleCondition); ok {
				rule.Condition = &condition{
					Key: c.Key,
				}
				if c.Value != "" {
					rule.Condition.Value = c.Value
				}
			}
			res.Rules = append(res.Rules, rule)
		}

		return res, nil
	}
}

func listPoliciesEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(identityReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		policies, err := svc.ListPolicies(req.key)
		if err != nil {
			return nil, err
		}
		res := listPoliciesRes{
			Policies: []viewPolicyRes{},
		}
		for _, policy := range policies {
			p := viewPolicyRes{
				ID:      policy.ID,
				Version: policy.Version,
				Owner:   policy.Owner,
				Name:    policy.Name,
				Rules:   []rule{},
			}
			for _, r := range policy.Rules {
				rule := rule{
					Action: r.Action,
					Type:   r.Type,
				}
				if c, ok := r.Condition.(auth.SimpleCondition); ok {
					rule.Condition = &condition{
						Key: c.Key,
					}
					if c.Value != "" {
						rule.Condition.Value = c.Value
					}
				}
				p.Rules = append(p.Rules, rule)
			}

			res.Policies = append(res.Policies, p)
		}

		return res, nil
	}
}

func removePolicyEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(identityReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.RemovePolicy(req.key, req.ID); err != nil {
			return nil, err
		}

		return removeRes{}, nil
	}
}

func attachPolicyEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(attachReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.AttachPolicy(req.key, req.policyID, req.userID); err != nil {
			return nil, err
		}

		return okRes{}, nil
	}
}

func detachPolicyEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(attachReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.DetachPolicy(req.key, req.policyID, req.userID); err != nil {
			return nil, err
		}

		return okRes{}, nil
	}
}
