package api

import (
	"context"
	"encoding/json"
	"io"
	"monetasa"
	"net/http"

	"monetasa/auth"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc auth.Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	r.Post("/user", kithttp.NewServer(
		registrationEndpoint(svc),
		decodeCredentials,
		encodeResponse,
		opts...,
	))

	r.Get("/user", kithttp.NewServer(
		viewEndpoint(svc),
		decodeIdentity,
		encodeResponse,
		opts...,
	))

	r.Put("/user", kithttp.NewServer(
		updateEndpoint(svc),
		decodeUpdate,
		encodeResponse,
		opts...,
	))

	r.Delete("/user", kithttp.NewServer(
		deleteEndpoint(svc),
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

	r.Get("/access-grant", kithttp.NewServer(
		canAccessEndpoint(svc),
		decodeIdentity,
		encodeResponse,
		opts...,
	))

	r.GetFunc("/version", monetasa.Version())
	r.Handle("/metrics", promhttp.Handler())

	return r
}

func decodeIdentity(_ context.Context, r *http.Request) (interface{}, error) {
	req := identityReq{
		key: r.Header.Get("Authorization"),
	}

	return req, nil
}

func decodeCredentials(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Header.Get("Content-Type") != contentType {
		return nil, auth.ErrUnsupportedContentType
	}

	var user auth.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return nil, err
	}

	return userReq{user}, nil
}

func decodeUpdate(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Header.Get("Content-Type") != contentType {
		return nil, auth.ErrUnsupportedContentType
	}

	var user auth.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return nil, err
	}

	req := updateReq{
		key:  r.Header.Get("Authorization"),
		user: user,
	}

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
	case auth.ErrUnsupportedContentType:
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
