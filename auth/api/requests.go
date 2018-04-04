package api

import (
	"github.com/asaskevich/govalidator"
	"monetasa/monetasa"
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

type addClientReq struct {
	key    string
	client auth.Client
}

func (req addClientReq) validate() error {
	if req.key == "" {
		return auth.ErrUnauthorizedAccess
	}

	return req.client.Validate()
}

type updateClientReq struct {
	key    string
	id     string
	client auth.Client
}

func (req updateClientReq) validate() error {
	if req.key == "" {
		return auth.ErrUnauthorizedAccess
	}

	if !govalidator.IsUUID(req.id) {
		return auth.ErrNotFound
	}

	return req.client.Validate()
}

type createChannelReq struct {
	key     string
	channel auth.Channel
}

func (req createChannelReq) validate() error {
	if req.key == "" {
		return auth.ErrUnauthorizedAccess
	}

	return nil
}

type updateChannelReq struct {
	key     string
	id      string
	channel auth.Channel
}

func (req updateChannelReq) validate() error {
	if req.key == "" {
		return auth.ErrUnauthorizedAccess
	}

	if !govalidator.IsUUID(req.id) {
		return auth.ErrNotFound
	}

	return nil
}

type viewResourceReq struct {
	key string
	id  string
}

func (req viewResourceReq) validate() error {
	if req.key == "" {
		return auth.ErrUnauthorizedAccess
	}

	if !govalidator.IsUUID(req.id) {
		return auth.ErrNotFound
	}

	return nil
}

type listResourcesReq struct {
	key    string
	size   int
	offset int
}

func (req listResourcesReq) validate() error {
	if req.key == "" {
		return auth.ErrUnauthorizedAccess
	}

	if req.size > 0 && req.offset >= 0 {
		return nil
	}

	return auth.ErrMalformedEntity
}

type connectionReq struct {
	key      string
	chanId   string
	clientId string
}

func (req connectionReq) validate() error {
	if req.key == "" {
		return auth.ErrUnauthorizedAccess
	}

	if !govalidator.IsUUID(req.chanId) && !govalidator.IsUUID(req.clientId) {
		return auth.ErrNotFound
	}

	return nil
}
