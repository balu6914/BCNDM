package api

import (
	"context"
	"encoding/json"
	// "fmt"
	"net/http"
	// "strconv"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"

	"monetasa/auth"
	"monetasa/dapp"
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(sr dapp.StreamRepository) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	r.Get("/version", kithttp.NewServer(
		versionEndpoint(),
		decodeVersionRequest,
		encodeResponse,
		opts...,
	))

	r.Post("/streams", kithttp.NewServer(
		saveStreamEndpoint(sr),
		decodeModifyStreamRequest,
		encodeResponse,
		opts...,
	))

	r.Put("/streams/:id", kithttp.NewServer(
		updateStreamEndpoint(sr),
		decodeModifyStreamRequest,
		encodeResponse,
		opts...,
	))

	r.Get("/streams/:id", kithttp.NewServer(
		oneStreamEndpoint(sr),
		decodeReadStreamRequest,
		encodeResponse,
		opts...,
	))

	r.Delete("/streams/:id", kithttp.NewServer(
		removeStreamEndpoint(sr),
		decodeModifyStreamRequest,
		encodeResponse,
		opts...,
	))

	// r.Get("/streams/search", kithttp.NewServer(
	// 	searchStreamEndpoint(sr),
	// 	decodeSearchStreamRequest,
	// 	encodeResponse,
	// opts...,
	// ))

	// r.Post("/streams/purch", kithttp.NewServer(
	// 	purchaseStreamEndpoint(sr),
	// 	decodePurchaseStreamRequest,
	// 	encodeResponse,
	// opts...,
	// ))

	return r
}

func decodeVersionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeModifyStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req modifyStreamReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.Id = bone.GetValue(r, "id")
	return req, nil
}

func decodeReadStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req readStreamReq
	req.Id = bone.GetValue(r, "id")
	return req, nil
}

// func decodeSearchStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
// 	var req searchStreamReq
// 	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 	// 	return nil, errp
// 	// }
// 	q := r.URL.Query()

// 	req.Type = q["type"][0]
// 	req.x0, _ = strconv.Atoi(q["x0"][0])
// 	req.y0, _ = strconv.Atoi(q["y0"][0])
// 	req.x1, _ = strconv.Atoi(q["x1"][0])
// 	req.y1, _ = strconv.Atoi(q["y1"][0])
// 	req.x2, _ = strconv.Atoi(q["x2"][0])
// 	req.y2, _ = strconv.Atoi(q["y2"][0])
// 	req.x3, _ = strconv.Atoi(q["x3"][0])
// 	req.y3, _ = strconv.Atoi(q["y3"][0])

// 	fmt.Println("decodeSearchStreamRequest")
// 	fmt.Println(req)
// 	return req, nil
// }

// func decodePurchaseStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
// 	var req purchaseStreamReq
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		return nil, err
// 	}
// 	return req, nil
// }

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", contentType)

	switch err {
	case auth.ErrMalformedEntity:
		w.WriteHeader(http.StatusBadRequest)
	case auth.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case auth.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case auth.ErrConflict:
		w.WriteHeader(http.StatusConflict)
	default:
		if _, ok := err.(*json.SyntaxError); ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
	}
}
