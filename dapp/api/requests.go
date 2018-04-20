package api

import (
	"github.com/asaskevich/govalidator"

	"monetasa/dapp"
)

type apiReq interface {
	validate() error
}

type createStreamReq struct {
	Stream dapp.Stream
}

func (req createStreamReq) validate() error {
	return req.Stream.Validate()
}

type updateStreamReq struct {
	Id     string
	Stream dapp.Stream
}

func (req updateStreamReq) validate() error {
	if req.Id == "" {
		return dapp.ErrMalformedData
	}

	if !govalidator.IsHexadecimal(req.Id) {
		return dapp.ErrMalformedData
	}

	return req.Stream.Validate()
}

type readStreamReq struct {
	Id string
}

func (req readStreamReq) validate() error {
	if req.Id == "" {
		return dapp.ErrMalformedData
	}

	if !govalidator.IsHexadecimal(req.Id) {
		return dapp.ErrMalformedData
	}

	return nil
}

type deleteStreamReq struct {
	Id string
}

func (req deleteStreamReq) validate() error {
	if req.Id == "" {
		return dapp.ErrMalformedData
	}

	if !govalidator.IsHexadecimal(req.Id) {
		return dapp.ErrMalformedData
	}

	return nil
}

type searchStreamReq struct {
	Type string
	x0   int
	y0   int
	x1   int
	y1   int
	x2   int
	y2   int
	x3   int
	y3   int
}

type purchaseStreamReq struct {
	Id    string
	Hours int `json:"hours"`
}
