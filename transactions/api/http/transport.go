package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"monetasa"
	"monetasa/transactions"
	"net/http"
	"time"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const contentType = "application/json"

var (
	errMalformedEntity        = errors.New("malformed entity")
	errUnauthorizedAccess     = errors.New("missing or invalid credentials provided")
	errUnsupportedContentType = errors.New("unsupported content type")
	auth                      monetasa.AuthServiceClient
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc transactions.Service, ac monetasa.AuthServiceClient) http.Handler {
	auth = ac

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()
	r.Get("/tokens", kithttp.NewServer(
		balanceEndpoint(svc),
		decodeBalanceReq,
		encodeResponse,
		opts...,
	))
	r.Post("/tokens/buy", kithttp.NewServer(
		buyEndpoint(svc),
		decodeBuyReq,
		encodeResponse,
		opts...,
	))
	r.Post("/tokens/withdraw", kithttp.NewServer(
		withdrawEndpoint(svc),
		decodeWithdrawReq,
		encodeResponse,
		opts...,
	))

	r.GetFunc("/version", monetasa.Version())
	r.Handle("/metrics", promhttp.Handler())

	return r
}

func decodeBalanceReq(_ context.Context, r *http.Request) (interface{}, error) {
	id, err := authorize(r)
	if err != nil {
		return nil, err
	}

	req := balanceReq{userID: id}

	return req, nil
}

func decodeBuyReq(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Header.Get("Content-Type") != contentType {
		return nil, errUnsupportedContentType
	}

	id, err := authorize(r)
	if err != nil {
		return nil, err
	}

	var req buyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	req.userID = id

	return req, nil
}

func decodeWithdrawReq(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Header.Get("Content-Type") != contentType {
		return nil, errUnsupportedContentType
	}

	id, err := authorize(r)
	if err != nil {
		return nil, err
	}

	var req withdrawReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	req.userID = id

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

func authorize(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return "", errUnauthorizedAccess
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	id, err := auth.Identify(ctx, &monetasa.Token{Value: token})
	if err != nil {
		e, ok := status.FromError(err)
		if ok && e.Code() == codes.Unauthenticated {
			return "", errUnauthorizedAccess
		}
		return "", err
	}

	return id.GetValue(), nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", contentType)

	switch err {
	case errMalformedEntity:
		w.WriteHeader(http.StatusBadRequest)
	case transactions.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
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
