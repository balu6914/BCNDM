package api

import (
	"gopkg.in/mgo.v2/bson"

	"monetasa/streams"
)

const (
	minLongitude        = -180
	maxLongitude        = 180
	minLatitude         = -90
	maxLatitude         = 90
	defLimit     uint64 = 20
	maxLimit     uint64 = 100
)

type apiReq interface {
	validate() error
}

type createStreamReq struct {
	owner  string
	stream streams.Stream
}

func (req createStreamReq) validate() error {
	if req.stream.Owner != "" && req.stream.Owner != req.owner {
		return streams.ErrMalformedData
	}

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

	if req.stream.ID != "" && req.id != req.stream.ID.Hex() {
		return streams.ErrMalformedData
	}

	if req.stream.Owner != "" && req.owner != req.stream.Owner {
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
	id string
}

func (req readStreamReq) validate() error {
	if req.id == "" {
		return streams.ErrMalformedData
	}

	if !bson.IsObjectIdHex(req.id) {
		return streams.ErrMalformedData
	}

	return nil
}

type deleteStreamReq struct {
	owner string
	id    string
}

func (req deleteStreamReq) validate() error {
	if req.id == "" {
		return streams.ErrMalformedData
	}

	if !bson.IsObjectIdHex(req.id) {
		return streams.ErrMalformedData
	}

	return nil
}

type searchStreamReq struct {
	// Coords are represented as Longitude and Latitude respectively.
	// It's important to note what the order of coordinates is.
	Name       string      `alias:"name"`
	StreamType string      `alias:"type"`
	Coords     [][]float64 `alias:"coords"`
	Page       uint64      `alias:"page"`
	Limit      uint64      `alias:"limit"`
	MinPrice   *uint64     `alias:"minPrice"`
	MaxPrice   *uint64     `alias:"maxPrice"`
}

func newSearchStreamReq() searchStreamReq {
	// Set default page and offset.
	return searchStreamReq{Limit: defLimit}
}

func (req searchStreamReq) validate() error {
	if req.MinPrice != nil && req.MaxPrice != nil {
		if *req.MaxPrice < *req.MinPrice {
			return streams.ErrMalformedData
		}
	}

	for _, p := range req.Coords {
		if p[0] < minLongitude || p[0] > maxLongitude ||
			p[1] < minLatitude || p[1] > maxLatitude {
			return streams.ErrMalformedData
		}
	}

	return nil
}
