package grpc

import "monetasa/auth"

type identityReq struct {
	token string
}

func (req identityReq) validate() error {
	if req.token == "" {
		return auth.ErrMalformedEntity
	}
	return nil
}

type emailReq struct {
	id string
}

func (req emailReq) validate() error {
	if req.id == "" {
		return auth.ErrMalformedEntity
	}
	return nil
}
