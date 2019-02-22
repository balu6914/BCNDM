package grpc

import "datapace/auth"

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
