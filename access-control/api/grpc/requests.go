package grpc

import access "datapace/access-control"

type partnersReq struct {
	id string
}

func (req partnersReq) validate() error {
	if req.id == "" {
		return access.ErrMalformedEntity
	}

	return nil
}
