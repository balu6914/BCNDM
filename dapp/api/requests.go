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
	x0   float64
	y0   float64
	x1   float64
	y1   float64
	x2   float64
	y2   float64
	x3   float64
	y3   float64
}

func (req searchStreamReq) validate() error {
	if req.Type != "geo" {
		return dapp.ErrUnknownType
	}

	if req.x0 < -180 || req.x0 > 180 ||
		req.x1 < -180 || req.x1 > 180 ||
		req.x2 < -180 || req.x2 > 180 ||
		req.x3 < -180 || req.x3 > 180 ||
		req.y0 < -90 || req.y0 > 90 ||
		req.y1 < -90 || req.y1 > 90 ||
		req.y2 < -90 || req.y2 > 90 ||
		req.y3 < -90 || req.y3 > 90 {
		return dapp.ErrMalformedData
	}

	return nil
}

// type purchaseStreamReq struct {
// 	Id    string
// 	Hours float64 `json:"hours"`
// }
