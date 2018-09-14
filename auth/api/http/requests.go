package http

import (
	"monetasa/auth"

	"github.com/asaskevich/govalidator"
)

type apiReq interface {
	validate() error
}

type registerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name,omitempty"`
}

const (
	maxEmailLength    = 32
	minPasswordLength = 8
	maxPasswordLength = 16
	maxNameLength     = 32
)

func (req registerReq) validate() error {
	if req.Email == "" || len(req.Email) > maxEmailLength ||
		req.Password == "" || len(req.Password) < minPasswordLength ||
		len(req.Password) > maxPasswordLength ||
		len(req.Name) > maxNameLength {
		return auth.ErrMalformedEntity
	}

	if !govalidator.IsEmail(req.Email) {
		return auth.ErrMalformedEntity
	}

	return nil
}

type credentialsReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req credentialsReq) validate() error {
	if req.Email == "" || req.Password == "" {
		return auth.ErrMalformedEntity
	}

	if !govalidator.IsEmail(req.Email) {
		return auth.ErrMalformedEntity
	}

	return nil
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
	key      string
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (req updateReq) validate() error {
	if req.key == "" {
		return auth.ErrUnauthorizedAccess
	}

	if req.Email == "" || req.Password == "" {
		return auth.ErrMalformedEntity
	}

	if !govalidator.IsEmail(req.Email) {
		return auth.ErrMalformedEntity
	}

	return nil
}
