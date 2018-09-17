package http

import (
	"monetasa/auth"

	"github.com/asaskevich/govalidator"
)

type apiReq interface {
	validate() error
}

type registerReq struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	ContactEmail string `json:"contact_email,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
}

func (req registerReq) validate() error {
	if req.Password == "" {
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
	key          string
	ContactEmail string `json:"contact_email,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
}

func (req updateReq) validate() error {
	if req.key == "" {
		return auth.ErrUnauthorizedAccess
	}

	if !govalidator.IsEmail(req.ContactEmail) {
		return auth.ErrMalformedEntity
	}

	return nil
}

type updatePasswordReq struct {
	key         string
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	RePassword  string `json:"re_password"`
}

func (req updatePasswordReq) validate() error {
	if req.OldPassword == "" || req.NewPassword == "" || req.RePassword == "" {
		return auth.ErrMalformedEntity
	}

	if req.NewPassword != req.RePassword {
		return auth.ErrMalformedEntity
	}

	return nil

}
