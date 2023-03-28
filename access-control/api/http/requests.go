package http

import (
	access "github.com/datapace/datapace/access-control"
	"github.com/datapace/datapace/auth"
)

type apiReq interface {
	validate() error
}

type requestAccessReq struct {
	key      string
	Receiver string `json:"receiver"`
}

func (req requestAccessReq) validate() error {
	if req.key == "" {
		return auth.ErrUnauthorizedAccess
	}

	if req.Receiver == "" {
		return auth.ErrMalformedEntity
	}

	return nil
}

type approveAccessReq struct {
	key string
	id  string
}

func (req approveAccessReq) validate() error {
	if req.key == "" {
		return auth.ErrUnauthorizedAccess
	}

	if req.id == "" {
		return access.ErrMalformedEntity
	}

	return nil
}

type revokeAccessReq struct {
	key string
	id  string
}

func (req revokeAccessReq) validate() error {
	if req.key == "" {
		return auth.ErrUnauthorizedAccess
	}

	if req.id == "" {
		return auth.ErrMalformedEntity
	}

	return nil
}

type listAccessRequestsReq struct {
	key   string
	state access.State
}

func (req listAccessRequestsReq) validate() error {
	if req.key == "" {
		return access.ErrUnauthorizedAccess
	}

	if req.state != access.Pending &&
		req.state != access.Approved &&
		req.state != access.Revoked &&
		req.state != access.State("") {
		return access.ErrMalformedEntity
	}

	return nil
}

type grantAccessReq struct {
	key        string
	dstUserIid string
}

func (req grantAccessReq) validate() error {
	if req.key == "" {
		return auth.ErrUnauthorizedAccess
	}

	if req.dstUserIid == "" {
		return auth.ErrMalformedEntity
	}

	return nil
}
