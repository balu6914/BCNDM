package api

import (
	"gopkg.in/mgo.v2/bson"

	"monetasa/streams"
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
	owner  string
	stream streams.Stream
}

func (req createStreamReq) validate() error {
	return req.stream.Validate()
}

type updateStreamReq struct {
	owner  string
	id     string
	stream streams.Stream
}

func (req updateStreamReq) validate() error {
	if req.id == "" {
		return streams.ErrMalformedData
	}

	if !bson.IsObjectIdHex(req.id) {
		return streams.ErrMalformedData
	}

	return req.stream.Validate()
}

type createBulkStreamReq struct {
	owner   string
	Streams []streams.Stream
}

func (req createBulkStreamReq) validate() error {
	for _, stream := range req.Streams {
		if err := stream.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type readStreamReq struct {
	streamID string
}

func (req readStreamReq) validate() error {
	if req.streamID == "" {
		return streams.ErrMalformedData
	}

	if !bson.IsObjectIdHex(req.streamID) {
		return streams.ErrMalformedData
	}

	return nil
}

type deleteStreamReq struct {
	owner    string
	streamID string
}

func (req deleteStreamReq) validate() error {
	if req.streamID == "" {
		return streams.ErrMalformedData
	}

	if !bson.IsObjectIdHex(req.streamID) {
		return streams.ErrMalformedData
	}

	return nil
}

type searchStreamReq struct {
	locationType string
	// Points are represented as Longitude and Latitude respectively.
	// It's important to note what the order of coordinates is.
	points [][]float64
}

func (req searchStreamReq) validate() error {
	if req.locationType != typeGeo {
		return streams.ErrUnknownType
	}

	n := len(req.points)
	for i := 0; i < n; i++ {
		p := req.points[i]
		if p[0] < minLongitude || p[0] > maxLongitude ||
			p[1] < minLatitude || p[1] > maxLatitude {
			return streams.ErrMalformedData
		}
	}

	return nil
}
