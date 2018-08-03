package grpc

import "monetasa/streams"

type oneReq struct {
	id string
}

func (req oneReq) validate() error {
	if req.id == "" {
		return streams.ErrMalformedData
	}

	return nil
}
