package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/datapace/datapace"

	"github.com/datapace/datapace/executions"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const contentType = "application/json"

var (
	errUnauthorizedAccess     = errors.New("missing or invalid credentials provided")
	errUnsupportedContentType = errors.New("unsupported content type")
	auth                      datapace.AuthServiceClient
)

// MakeHandler returns HTTP handler for executions serivce.
func MakeHandler(svc executions.Service, ac datapace.AuthServiceClient) http.Handler {
	auth = ac

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	r.Post("/executions", kithttp.NewServer(
		startExecutionEndpoint(svc),
		decodeStartExecutionRequest,
		encodeResponse,
		opts...,
	))

	r.Get("/executions/:id/result", kithttp.NewServer(
		resultEndpoint(svc),
		decodeViewRequest,
		encodeResponse,
		opts...,
	))

	r.Get("/executions/:id", kithttp.NewServer(
		viewEndpoint(svc),
		decodeViewRequest,
		encodeResponse,
		opts...,
	))

	r.Get("/executions", kithttp.NewServer(
		listEndpoint(svc),
		decodeListRequest,
		encodeResponse,
		opts...,
	))

	r.GetFunc("/version", datapace.Version())

	return r
}

func decodeStartExecutionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	owner, err := authorize(r)
	if err != nil {
		return nil, err
	}

	req := startExecutionReq{
		owner: owner,
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return req, nil
}

func decodeViewRequest(_ context.Context, r *http.Request) (interface{}, error) {
	owner, err := authorize(r)
	if err != nil {
		return nil, err
	}

	req := viewReq{
		owner: owner,
		id:    bone.GetValue(r, "id"),
	}

	return req, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	owner, err := authorize(r)
	if err != nil {
		return nil, err
	}

	req := listReq{
		owner: owner,
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
	case executions.ErrMalformedData:
		w.WriteHeader(http.StatusBadRequest)
	case errUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
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

func authorize(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return "", errUnauthorizedAccess
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	id, err := auth.Identify(ctx, &datapace.Token{Value: token})
	if err != nil {
		e, ok := status.FromError(err)
		if ok && e.Code() == codes.Unauthenticated {
			return "", errUnauthorizedAccess
		}
		return "", err
	}

	return id.GetValue(), nil
}
