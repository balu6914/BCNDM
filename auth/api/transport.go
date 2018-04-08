package api

import (
	"context"
	"io"
	"io/ioutil"
	"monetasa"
	"net/http"
	"strings"

	"monetasa/auth"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc auth.Service, mc auth.AuthClient) http.Handler {
	auth = mc

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	r.Get("/users", kithttp.NewServer(
		getUsersEndpoint(svc),
		decodeRequest,
		encodeResponse,
		opts...,
	))

	r.Get("/users/:id", kithttp.NewServer(
		getUserEndpoint(svc),
		decodeRequest,
		encodeResponse,
		opts...,
	))

	r.Post("/users", kithttp.NewServer(
		createUserEndpoint(svc),
		decodeRequest,
		encodeResponse,
		opts...,
	))

	r.Delete("/users/:id", kithttp.NewServer(
		deletetUserEndpoint(svc),
		decodeRequest,
		encodeResponse,
		opts...,
	))

	r.Post("/login", kithttp.NewServer(
		loginEndpoint(svc),
		decodeRequest,
		encodeResponse,
		opts...,
	))

	r.Post("/auth", kithttp.NewServer(
		authEndpoint(svc),
		decodeRequest,
		encodeResponse,
		opts...,
	))

	r.GetFunc("/version", monetasa.Version())
	r.Handle("/metrics", promhttp.Handler())

	return r
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	ct, err := checkContentType(r)
	if err != nil {
		return nil, err
	}

	publisher, err := authorize(r)
	if err != nil {
		return nil, err
	}

	payload, err := decodePayload(r.Body)
	if err != nil {
		return nil, err
	}

	channel := bone.GetValue(r, "id")

	msg := writer.RawMessage{
		Publisher:   publisher,
		Protocol:    protocol,
		ContentType: ct,
		Channel:     channel,
		Payload:     payload,
	}

	return msg, nil
}

func authorize(r *http.Request) (string, error) {
	apiKey := r.Header.Get("Authorization")

	if apiKey == "" {
		return "", errUnauthorizedAccess
	}

	// Path is `/channels/:id/messages`, we need chanID.
	c := strings.Split(r.URL.Path, "/")[2]

	id, err := auth.CanAccess(c, apiKey)
	if err != nil {
		return "", err
	}

	return id, nil
}

func checkContentType(r *http.Request) (string, error) {
	ct := r.Header.Get("Content-Type")

	if ct != ctJson {
		return "", errUnknownType
	}

	return ct, nil
}

func decodePayload(body io.ReadCloser) ([]byte, error) {
	payload, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, errMalformedData
	}
	defer body.Close()

	return payload, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.WriteHeader(http.StatusAccepted)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	switch err {
	case errMalformedData:
		w.WriteHeader(http.StatusBadRequest)
	case errUnknownType:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	case auth.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
