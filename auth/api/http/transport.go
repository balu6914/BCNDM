package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/datapace/datapace"

	"github.com/datapace/datapace/auth"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var errUnsupportedContentType = errors.New("unsupported content type")

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc auth.Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	r.Post("/users", kithttp.NewServer(
		registrationEndpoint(svc),
		decodeRegister,
		encodeResponse,
		opts...,
	))

	r.Get("/users", kithttp.NewServer(
		listEndpoint(svc),
		decodeIdentity,
		encodeResponse,
		opts...,
	))

	r.Get("/users/:id", kithttp.NewServer(
		viewEndpoint(svc),
		decodeIdentity,
		encodeResponse,
		opts...,
	))

	r.Put("/users", kithttp.NewServer(
		updateEndpoint(svc),
		decodeUpdate,
		encodeResponse,
		opts...,
	))

	r.Get("/users/non-partners", kithttp.NewServer(
		nonPartnersEndpoint(svc),
		decodeIdentity,
		encodeResponse,
		opts...,
	))

	r.Post("/tokens", kithttp.NewServer(
		loginEndpoint(svc),
		decodeCredentials,
		encodeResponse,
		opts...,
	))

	r.Put("/users/password", kithttp.NewServer(
		updatePasswordEndpoint(svc),
		decodePasswordUpdate,
		encodeResponse,
		opts...,
	))

	r.GetFunc("/version", datapace.Version())
	r.Handle("/metrics", promhttp.Handler())

	return r
}

func decodeRegister(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Header.Get("Content-Type") != contentType {
		return nil, errUnsupportedContentType
	}
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.key = r.Header.Get("Authorization")

	return req, nil
}

func decodeIdentity(_ context.Context, r *http.Request) (interface{}, error) {
	req := identityReq{
		key: r.Header.Get("Authorization"),
	}
	req.ID = bone.GetValue(r, "id")
	return req, nil
}

func decodeUpdate(_ context.Context, r *http.Request) (interface{}, error) {
	var req updateReq
	if r.Header.Get("Content-Type") != contentType {
		return nil, errUnsupportedContentType
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.key = r.Header.Get("Authorization")

	return req, nil
}

func decodeCredentials(_ context.Context, r *http.Request) (interface{}, error) {
	var req credentialsReq
	if r.Header.Get("Content-Type") != contentType {
		return nil, errUnsupportedContentType
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodePasswordUpdate(_ context.Context, r *http.Request) (interface{}, error) {
	var req updatePasswordReq
	if r.Header.Get("Content-Type") != contentType {
		return nil, errUnsupportedContentType
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	req.key = r.Header.Get("Authorization")

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
	case auth.ErrMalformedEntity:
		w.WriteHeader(http.StatusBadRequest)
	case auth.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case auth.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case auth.ErrConflict:
		w.WriteHeader(http.StatusConflict)
	case errUnsupportedContentType:
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
