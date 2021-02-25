package http

import (
	"errors"

	"github.com/datapace/datapace/auth"

	"github.com/asaskevich/govalidator"
)

// User validation errors
var (
	errInvalidEmail         = errors.New("invalid user email")
	errInvalidPassword      = errors.New("invalid password length")
	errInvalidFirstName     = errors.New("invalid first name length")
	errInvalidLastName      = errors.New("invalid last name length")
	errInvalidCompany       = errors.New("invalid company name length")
	errInvalidPhone         = errors.New("invalid phone number length")
	errInvalidAddress       = errors.New("invalid address length")
	errInvalidRole          = errors.New("invalid role")
	errInvalidPolicyRules   = errors.New("invalid policy rules list")
	errInvalidPolicyVersion = errors.New("invalid policy version")
)

type apiReq interface {
	validate() error
}

type registerReq struct {
	key          string
	Email        string `json:"email"`
	Password     string `json:"password"`
	ContactEmail string `json:"contact_email,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Company      string `json:"company,omitempty"`
	Address      string `json:"address,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Role         string `json:"role,omitempty"`
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
		return errInvalidEmail
	}

	if req.Password == "" || len(req.Password) < minPasswordLength ||
		len(req.Password) > maxPasswordLength {
		return errInvalidPassword
	}

	if len(req.FirstName) > maxNameLength {
		return errInvalidFirstName
	}

	if len(req.LastName) > maxNameLength {
		return errInvalidLastName
	}

	if len(req.Company) > maxCompanyLength {
		return errInvalidCompany
	}

	if len(req.Phone) > maxPhoneLength {
		return errInvalidPhone
	}

	if len(req.Address) > maxAddressLength {
		return errInvalidAddress
	}

	if req.Role != auth.UserRole && req.Role != auth.AdminRole {
		return errInvalidRole
	}

	if !govalidator.IsEmail(req.Email) {
		return errInvalidEmail
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
	key          string
	id           string
	ContactEmail *string `json:"contact_email,omitempty"`
	FirstName    *string `json:"first_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
	Company      *string `json:"company,omitempty"`
	Address      *string `json:"address,omitempty"`
	Phone        *string `json:"phone,omitempty"`
	Role         *string `json:"role,omitempty"`
	Password     *string `json:"password"`
	Disabled     *bool   `json:"disabled"`
	Locked       *bool   `json:"locked"`
}

func (req updateReq) validate() error {
	if req.id == "" {
		return auth.ErrMalformedEntity
	}
	if req.ContactEmail != nil && (!govalidator.IsEmail(*req.ContactEmail) ||
		len(*req.ContactEmail) > maxEmailLength) {
		return auth.ErrMalformedEntity
	}
	return nil
}

func (req updateReq) toUser() auth.User {
	ret := auth.User{
		ID: req.id,
	}
	if req.Address != nil {
		ret.Address = *req.Address
	}
	if req.Company != nil {
		ret.Company = *req.Company
	}
	if req.ContactEmail != nil {
		ret.ContactEmail = *req.ContactEmail
	}
	if req.Disabled != nil {
		ret.Disabled = *req.Disabled
	}
	if req.Locked != nil {
		ret.Locked = *req.Locked
	}
	if req.FirstName != nil {
		ret.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		ret.LastName = *req.LastName
	}
	if req.Password != nil {
		ret.Password = *req.Password
	}
	if req.Phone != nil {
		ret.Phone = *req.Phone
	}
	if req.Role != nil {
		ret.Role = *req.Role
	}
	return ret
}

type policyRequest struct {
	key     string
	Version string `json:"version,omitempty"`
	Name    string `json:"name,omitempty"`
	Rules   []rule `json:"rules,omitempty"`
}

type rule struct {
	Action    auth.Action `json:"action"`
	Type      string      `json:"type"`
	Condition *condition  `json:"condition,omitempty"`
}

type condition struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func (pr policyRequest) validate() error {
	if pr.Version == "" {
		return errInvalidPolicyVersion
	}
	if pr.Rules == nil || len(pr.Rules) == 0 {
		return errInvalidPolicyRules
	}
	return nil
}

type attachReq struct {
	key      string
	policyID string
	userID   string
}

func (ar attachReq) validate() error {
	if ar.key == "" || ar.policyID == "" || ar.userID == "" {
		return auth.ErrMalformedEntity
	}
	return nil
}
