package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"monetasa"
	"monetasa/subscriptions"
)

var auth monetasa.AuthServiceClient

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc subscriptions.Service, ac monetasa.AuthServiceClient) http.Handler {
	auth = ac

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	r.Get("/subscriptions", kithttp.NewServer(
		readSubsEndpoint(svc),
		decodeReadSubsRequest,
		encodeResponse,
		opts...,
	))

	r.Post("/subscriptions", kithttp.NewServer(
		createSubsEndpoint(svc),
		decodeCreateSubsRequest,
		encodeResponse,
		opts...,
	))

	r.GetFunc("/version", monetasa.Version())

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
	id, err := auth.Identify(ctx, &monetasa.Token{Value: token})
	if err != nil {
		e, ok := status.FromError(err)
		if ok && e.Code() == codes.Unauthenticated {
			return "", subscriptions.ErrUnauthorizedAccess
		}
		return "", err
	}

	return id.GetValue(), nil
}

func decodeCreateSubsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	userID, err := authorize(r)
	if err != nil {
		return nil, err
	}

	sub := subscriptions.Subscription{}
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		return nil, err
	}

	req := subscriptionReq{
		UserID:       userID,
		Subscription: sub,
	}
	return req, nil
}

func decodeReadSubsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	userID, err := authorize(r)
	if err != nil {
		return nil, err
	}

	req := listSubsReq{
		UserID: userID,
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
	case subscriptions.ErrUnsupportedContentType:
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
