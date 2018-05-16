package api

import (
	"github.com/asaskevich/govalidator"

	"monetasa/dapp"
)

const (
	minLongitude = -180
	maxLongitude = 180
	minLatitude  = -90
	maxLatitude  = 90
	typeGeo      = "geo"
)

type apiReq interface {
	validate() error
}

type createStreamReq struct {
	User   string
	Stream dapp.Stream
}

func (req createStreamReq) validate() error {
	return req.Stream.Validate()
}

type updateStreamReq struct {
	User     string
	StreamId string
	Stream   dapp.Stream
}

func (req updateStreamReq) validate() error {
	if req.StreamId == "" {
		return dapp.ErrMalformedData
	}

	if !govalidator.IsHexadecimal(req.StreamId) {
		return dapp.ErrMalformedData
	}

	return req.Stream.Validate()
}

type readStreamReq struct {
	StreamId string
}

func (req readStreamReq) validate() error {
	if req.StreamId == "" {
		return dapp.ErrMalformedData
	}

	if !govalidator.IsHexadecimal(req.StreamId) {
		return dapp.ErrMalformedData
	}

	return nil
}

type deleteStreamReq struct {
	User     string
	StreamId string
}

func (req deleteStreamReq) validate() error {
	if req.StreamId == "" {
		return dapp.ErrMalformedData
	}

	if !govalidator.IsHexadecimal(req.StreamId) {
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
	if req.Type != typeGeo {
		return dapp.ErrUnknownType
	}

	if req.x0 < minLongitude || req.x0 > maxLongitude ||
		req.x1 < minLongitude || req.x1 > maxLongitude ||
		req.x2 < minLongitude || req.x2 > maxLongitude ||
		req.x3 < minLongitude || req.x3 > maxLongitude ||
		req.y0 < minLatitude || req.y0 > maxLatitude ||
		req.y1 < minLatitude || req.y1 > maxLatitude ||
		req.y2 < minLatitude || req.y2 > maxLatitude ||
		req.y3 < minLatitude || req.y3 > maxLatitude {
		return dapp.ErrMalformedData
	}

	return nil
}
