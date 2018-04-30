package api

import (
	"context"
	"encoding/json"
	// "fmt"
	"io"
	"net/http"
	"strconv"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"gopkg.in/mgo.v2/bson"

	"monetasa"
	"monetasa/dapp"
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(sr dapp.Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	r.Post("/streams", kithttp.NewServer(
		addStreamEndpoint(sr),
		decodeCreateStreamRequest,
		encodeResponse,
		opts...,
	))

	r.Get("/streams/search", kithttp.NewServer(
		searchStreamEndpoint(sr),
		decodeSearchStreamRequest,
		encodeResponse,
		opts...,
	))

	r.Put("/streams/:id", kithttp.NewServer(
		updateStreamEndpoint(sr),
		decodeUpdateStreamRequest,
		encodeResponse,
		opts...,
	))

	r.Get("/streams/:id", kithttp.NewServer(
		viewStreamEndpoint(sr),
		decodeReadStreamRequest,
		encodeResponse,
		opts...,
	))

	r.Delete("/streams/:id", kithttp.NewServer(
		removeStreamEndpoint(sr),
		decodeDeleteStreamRequest,
		encodeResponse,
		opts...,
	))

	r.GetFunc("/version", monetasa.Version())

	return r
}

func decodeCreateStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var stream dapp.Stream
	if err := json.NewDecoder(r.Body).Decode(&stream); err != nil {
		return nil, err
	}
	stream.ID = bson.NewObjectId()

	return createStreamReq{stream}, nil
}

func decodeUpdateStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var stream dapp.Stream
	if err := json.NewDecoder(r.Body).Decode(&stream); err != nil {
		return nil, err
	}

	req := updateStreamReq{
		Id:     bone.GetValue(r, "id"),
		Stream: stream,
	}
	return req, nil
}

func decodeReadStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := readStreamReq{
		Id: bone.GetValue(r, "id"),
	}
	return req, nil
}

func decodeDeleteStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := deleteStreamReq{
		Id: bone.GetValue(r, "id"),
	}
	return req, nil

}

func decodeSearchStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	q := r.URL.Query()

	var req searchStreamReq
	req.Type = q["type"][0]
	req.x0, _ = strconv.ParseFloat(q["x0"][0], 64)
	req.y0, _ = strconv.ParseFloat(q["y0"][0], 64)
	req.x1, _ = strconv.ParseFloat(q["x1"][0], 64)
	req.y1, _ = strconv.ParseFloat(q["y1"][0], 64)
	req.x2, _ = strconv.ParseFloat(q["x2"][0], 64)
	req.y2, _ = strconv.ParseFloat(q["y2"][0], 64)
	req.x3, _ = strconv.ParseFloat(q["x3"][0], 64)
	req.y3, _ = strconv.ParseFloat(q["y3"][0], 64)

	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", contentType)

	if ar, ok := response.(apiRes); ok {
		for k, v := range ar.headers() {
			w.Header().Set(k, v)
		}

		w.WriteHeader(ar.code())

		if ar.empty() {
			return nil
		}
	}

	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", contentType)

	switch err {
	case dapp.ErrMalformedEntity:
		w.WriteHeader(http.StatusBadRequest)
	case dapp.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case dapp.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case dapp.ErrConflict:
		w.WriteHeader(http.StatusConflict)
	case dapp.ErrUnsupportedContentType:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	case io.ErrUnexpectedEOF:
		w.WriteHeader(http.StatusBadRequest)
	case io.EOF:
		w.WriteHeader(http.StatusBadRequest)
	default:
		switch err.(type) {
		case *json.SyntaxError:
			w.WriteHeader(http.StatusBadRequest)
		case *json.UnmarshalTypeError:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
