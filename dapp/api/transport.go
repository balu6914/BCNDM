package api

import (
	"context"
	"encoding/json"
	"io"
	// "fmt"
	"net/http"
	// "strconv"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"

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
		decodeCreateStreamRequest,
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
		oneStreamEndpoint(sr),
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

func decodeCreateStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var stream dapp.Stream
	if err := json.NewDecoder(r.Body).Decode(&stream); err != nil {
		return nil, err
	}

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

// func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
// 	return json.NewEncoder(w).Encode(response)
// }

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
