package http

import (
	"github.com/datapace/datapace/auth"

	"github.com/asaskevich/govalidator"
)

type apiReq interface {
	validate() error
}

type registerReq struct {
	key          string
	Email        string   `json:"email"`
	Password     string   `json:"password"`
	ContactEmail string   `json:"contact_email,omitempty"`
	FirstName    string   `json:"first_name,omitempty"`
	LastName     string   `json:"last_name,omitempty"`
	Company      string   `json:"company,omitempty"`
	Address      string   `json:"address,omitempty"`
	Phone        string   `json:"phone,omitempty"`
	Roles        []string `json:"roles,omitempty"`
}

const (
	maxEmailLength    = 254
	minPasswordLength = 8
	maxPasswordLength = 32
	maxNameLength     = 32
	maxCompanyLength  = 32
	maxPhoneLength    = 32
	maxAddressLength  = 128
)

func (req registerReq) validate() error {
	if req.Email == "" || len(req.Email) > maxEmailLength {
		return auth.ErrMalformedEntity
	}

	if req.Password == "" || len(req.Password) < minPasswordLength ||
		len(req.Password) > maxPasswordLength {
		return auth.ErrMalformedEntity
	}

	if len(req.FirstName) > maxNameLength {
		return auth.ErrMalformedEntity
	}

	if len(req.LastName) > maxNameLength {
		return auth.ErrMalformedEntity
	}

	if len(req.Company) > maxCompanyLength {
		return auth.ErrMalformedEntity
	}

	if len(req.Phone) > maxPhoneLength {
		return auth.ErrMalformedEntity
	}

	if len(req.Address) > maxAddressLength {
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
	ID  string
}

func (req identityReq) validate() error {
	if req.key == "" {
		return auth.ErrUnauthorizedAccess
	}

	return nil
}

type updateReq struct {
	ContactEmail string   `json:"contact_email,omitempty"`
	FirstName    string   `json:"first_name,omitempty"`
	LastName     string   `json:"last_name,omitempty"`
	Company      string   `json:"company,omitempty"`
	Address      string   `json:"address,omitempty"`
	Phone        string   `json:"phone,omitempty"`
	Roles        []string `json:"roles,omitempty"`
	Password     string   `json:"password"`
	Disabled     bool     `json:"disabled"`
}

func (req updateReq) validate() error {
	if !govalidator.IsEmail(req.ContactEmail) ||
		len(req.ContactEmail) > maxEmailLength {
		return auth.ErrMalformedEntity
	}
	return nil
}
