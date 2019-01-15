package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"datapace"
	"datapace/subscriptions"
)

var (
	errUnsupportedContentType = errors.New("unsupported content type")
	auth                      datapace.AuthServiceClient
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc subscriptions.Service, ac datapace.AuthServiceClient) http.Handler {
	auth = ac

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	r.Get("/subscriptions/owned", kithttp.NewServer(
		searchSubsEndpoint(svc),
		decodeOwnedSubsRequest,
		encodeResponse,
		opts...,
	))

	r.Get("/subscriptions/bought", kithttp.NewServer(
		searchSubsEndpoint(svc),
		decodeBoughtSubsRequest,
		encodeResponse,
		opts...,
	))

	r.Get("/subscriptions/:id", kithttp.NewServer(
		viewSubEndpoint(svc),
		decodeViewSubRequest,
		encodeResponse,
		opts...,
	))

	r.Post("/subscriptions", kithttp.NewServer(
		addSubEndpoint(svc),
		decodeAddSubRequest,
		encodeResponse,
		opts...,
	))

	r.GetFunc("/version", datapace.Version())
	r.Handle("/metrics", promhttp.Handler())

	return r
}

func authorize(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return "", subscriptions.ErrUnauthorizedAccess
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// TODO: move this code to auth interface in service package root.
	id, err := auth.Identify(ctx, &datapace.Token{Value: token})
	if err != nil {
		e, ok := status.FromError(err)
		if ok && e.Code() == codes.Unauthenticated {
			return "", subscriptions.ErrUnauthorizedAccess
		}
		return "", err
	}

	return id.GetValue(), nil
}

func decodeSearch(r *http.Request) (searchSubsReq, error) {
	q := r.URL.Query()

	req := searchSubsReq{
		Limit: 20,
	}

	if err := searchFields(&req, q); err != nil {
		return searchSubsReq{}, err
	}

	return req, nil
}

func decodeAddSubRequest(_ context.Context, r *http.Request) (interface{}, error) {
	userID, err := authorize(r)
	if err != nil {
		return nil, err
	}

	sub := subscriptions.Subscription{}
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		return nil, err
	}

	sub.UserID = userID

	req := addSubReq{
		UserID:       userID,
		UserToken:    r.Header.Get("Authorization"),
		Subscription: sub,
	}

	return req, nil
}

func decodeBoughtSubsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	userID, err := authorize(r)
	if err != nil {
		return nil, err
	}

	req, err := decodeSearch(r)
	if err != nil {
		return nil, err
	}

	req.UserID = userID

	return req, nil
}

func decodeOwnedSubsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	userID, err := authorize(r)
	if err != nil {
		return nil, err
	}

	req, err := decodeSearch(r)
	if err != nil {
		return nil, err
	}

	req.StreamOwner = userID

	return req, nil
}

func decodeViewSubRequest(_ context.Context, r *http.Request) (interface{}, error) {
	userID, err := authorize(r)
	if err != nil {
		return nil, err
	}

	req := viewSubReq{
		userID:         userID,
		subscriptionID: bone.GetValue(r, "id"),
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
	case subscriptions.ErrMalformedEntity:
		w.WriteHeader(http.StatusBadRequest)
	case subscriptions.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case subscriptions.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case subscriptions.ErrConflict:
		w.WriteHeader(http.StatusConflict)
	case subscriptions.ErrNotEnoughTokens:
		w.WriteHeader(http.StatusPaymentRequired)
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
