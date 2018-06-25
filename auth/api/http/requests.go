package http

import (
	"monetasa/auth"

	"github.com/asaskevich/govalidator"
)

type apiReq interface {
	validate() error
}

type userReq struct {
	user auth.User
}

func (req userReq) validate() error {
	return req.user.Validate()
}

type identityReq struct {
	key string
}

func (req identityReq) validate() error {
	if req.key == "" {
		return auth.ErrUnauthorizedAccess
	}

	return nil
}

type updateReq struct {
	key  string
	user auth.User
}

func (req updateReq) validate() error {
	if req.key == "" {
		return auth.ErrUnauthorizedAccess
	}

	return req.user.Validate()
}

type viewReq struct {
	key string
	id  string
}

func (req viewReq) validate() error {
	if req.key == "" {
		return auth.ErrUnauthorizedAccess
	}

	if !govalidator.IsUUID(req.id) {
		return auth.ErrNotFound
	}

	return nil
}
