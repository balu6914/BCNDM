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
