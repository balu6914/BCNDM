package grpc

import "github.com/datapace/datapace/auth"

type identityReq struct {
	token string
}

func (req identityReq) validate() error {
	if req.token == "" {
		return auth.ErrMalformedEntity
	}

	return nil
}

type existsReq struct {
	id string
}

func (req existsReq) validate() error {
	if req.id == "" {
		return auth.ErrMalformedEntity
	}

	return nil
}

type authReq struct {
	token        string
	action       int64
	resourceType string
	attributes   map[string]string
}

// authRequest implements auth.Resource interface.
func (ar authReq) Attributes() map[string]string {
	return ar.attributes
}

func (ar authReq) ResourceType() string {
	return ar.resourceType
}

func (ar authReq) validate() error {
	if ar.token == "" {
		return auth.ErrUnauthorizedAccess
	}

	return nil
}
