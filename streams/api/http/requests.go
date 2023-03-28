package http

import (
	"fmt"

	"github.com/datapace/datapace/errors"
	"gopkg.in/mgo.v2/bson"

	"github.com/datapace/datapace/streams"
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

func (req *addStreamReq) validate() error {
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
	user        string
	Name        string      `alias:"name" json:"name"`
	Owner       string      `alias:"owner" json:"owner"`
	StreamType  string      `alias:"type" json:"type"`
	Coords      [][]float64 `alias:"coords" json:"coords"`
	Page        uint64      `alias:"page" json:"page"`
	Limit       uint64      `alias:"limit" json:"limit"`
	MinPrice    *uint64     `alias:"minPrice" json:"minPrice"`
	MaxPrice    *uint64     `alias:"maxPrice" json:"maxPrice"`
	SubCategory string      `alias:"subcategory" json:"subcategory"`

	// Metadata is the stream metadata constraint
	Metadata map[string]interface{} `json:"metadata"`
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
		return errors.Wrap(errors.New("stream id is missing in the request"), streams.ErrMalformedData)
	}

	if req.stream.ID != "" && req.id != req.stream.ID {
		return errors.Wrap(fmt.Errorf("stream id %s change to %s is not allowed", req.stream.ID, req.id), streams.ErrMalformedData)
	}

	if req.stream.Owner != "" && req.owner != req.stream.Owner {
		return errors.Wrap(fmt.Errorf("stream owner %s change to %s is not allowed", req.stream.Owner, req.owner), streams.ErrMalformedData)
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

type exportStreamsReq struct {
	owner string
}

func (req exportStreamsReq) validate() error {
	if !bson.IsObjectIdHex(req.owner) {
		return streams.ErrMalformedData
	}
	return nil
}

type addCategoryReq struct {
	key              string
	CategoryName     string   `json:"categoryname"`
	SubCategoryNames []string `json:"subcategorynames"`
}

func (req addCategoryReq) validate() error {
	if req.CategoryName == "" {
		return streams.ErrMalformedData
	}

	return nil
}

type listCategoryReq struct {
	key string
}

func (req listCategoryReq) validate() error {
	if req.key == "" {
		return streams.ErrMalformedData
	}

	return nil
}
