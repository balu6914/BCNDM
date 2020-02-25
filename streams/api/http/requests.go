package http

import (
	"gopkg.in/mgo.v2/bson"

	"datapace/streams"
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

type addStreamReq struct {
	owner  string
	stream streams.Stream
}

func (req addStreamReq) validate() error {
	if req.stream.Owner != "" && req.stream.Owner != req.owner {
		return streams.ErrMalformedData
	}

	return req.stream.Validate()
}

type addBulkStreamsReq struct {
	owner   string
	Streams []streams.Stream
}

func (req addBulkStreamsReq) validate() error {
	for _, stream := range req.Streams {
		if err := stream.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type searchStreamsReq struct {
	// Coords are represented as Longitude and Latitude respectively.
	// It's important to note what the order of coordinates is.
	// Fields `Owner` and `user` represent owner of the Stream and
	// the user who sent request respectively.
	user       string
	Name       string      `alias:"name"`
	Owner      string      `alias:"owner"`
	StreamType string      `alias:"type"`
	Coords     [][]float64 `alias:"coords"`
	Page       uint64      `alias:"page"`
	Limit      uint64      `alias:"limit"`
	MinPrice   *uint64     `alias:"minPrice"`
	MaxPrice   *uint64     `alias:"maxPrice"`
}

func newSearchStreamReq() searchStreamsReq {
	// Set default page and offset.
	return searchStreamsReq{Limit: defLimit}
}

func (req searchStreamsReq) validate() error {
	if !bson.IsObjectIdHex(req.user) {
		return streams.ErrMalformedData
	}

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

type updateStreamReq struct {
	owner  string
	id     string
	stream streams.Stream
}

func (req updateStreamReq) validate() error {
	if req.id == "" {
		return streams.ErrMalformedData
	}

	if req.stream.ID != "" && req.id != req.stream.ID {
		return streams.ErrMalformedData
	}

	if req.stream.Owner != "" && req.owner != req.stream.Owner {
		return streams.ErrMalformedData
	}

	return req.stream.Validate()
}

type viewStreamReq struct {
	id    string
	owner string
}

func (req viewStreamReq) validate() error {
	if !bson.IsObjectIdHex(req.id) || !bson.IsObjectIdHex(req.owner) {
		return streams.ErrMalformedData
	}

	return nil
}

type removeStreamReq struct {
	owner string
	id    string
}

func (req removeStreamReq) validate() error {
	if req.id == "" {
		return streams.ErrMalformedData
	}

	if !bson.IsObjectIdHex(req.id) {
		return streams.ErrMalformedData
	}

	return nil
}
