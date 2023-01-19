package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/datapace/datapace"

	"github.com/datapace/datapace/auth"
	authproto "github.com/datapace/datapace/proto/auth"
	"github.com/datapace/datapace/proto/common"
	"github.com/datapace/datapace/subscriptions"
)

const resourceType = "subscription"

var (
	errUnsupportedContentType = errors.New("unsupported content type")
	authClient                authproto.AuthServiceClient
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc subscriptions.Service, ac authproto.AuthServiceClient) http.Handler {
	authClient = ac

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

	r.Get("/owner/:ownerID/stream/:streamID/subscriptions", kithttp.NewServer(
		viewSubByUserAndStreamEndpoint(svc),
		decodeViewSubByUserAndStreamRequest,
		encodeResponse,
		opts...,
	))

	r.Get("/subscriptions/bought", kithttp.NewServer(
		searchSubsEndpoint(svc),
		decodeBoughtSubsRequest,
		encodeResponse,
		opts...,
	))

	r.Get("/subscriptions/report", kithttp.NewServer(
		reportSubsEndpoint(svc),
		decodeReportRequest,
		encodeBinaryResponse,
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

func decodeSearch(r *http.Request) (searchSubsReq, error) {
	q := r.URL.Query()
	req := searchSubsReq{
		Query: subscriptions.Query{
			Limit: 20,
			Page:  20,
		},
	}
	for k, v := range q {
		if len(v) != 0 && v[0] != "" {
			switch k {
			case "page":
				p, err := strconv.ParseUint(v[0], 10, 64)
				if err != nil {
					return req, err
				}
				req.Page = p
			case "limit":
				l, err := strconv.ParseUint(v[0], 10, 64)
				if err != nil {
					return req, err
				}
				req.Limit = l
			case "streamId":
				req.StreamID = v[0]
			case "startTime":
				st, err := time.Parse(time.ANSIC, v[0])
				if err != nil {
					return req, err
				}
				req.StartTime = &st
			}
		}
	}

	return req, nil
}

func decodeAddSubRequest(_ context.Context, r *http.Request) (interface{}, error) {
	ar := &authproto.AuthRequest{
		Action: int64(auth.CreateBulk),
		Type:   resourceType,
		Token:  r.Header.Get("Authorization"),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	userID, err := authClient.Authorize(ctx, ar)
	if err != nil {
		return nil, subscriptions.ErrUnauthorizedAccess
	}

	subs := []subscriptions.Subscription{}

	if err := json.NewDecoder(r.Body).Decode(&subs); err != nil {
		return nil, err
	}
	uid := userID.GetValue()

	for i := range subs {
		subs[i].UserID = uid
	}

	req := addSubsReq{
		UserID:        uid,
		UserToken:     r.Header.Get("Authorization"),
		Subscriptions: subs,
	}

	return req, nil
}

func authorize(r *http.Request) (*common.ID, error) {
	ar := &authproto.AuthRequest{
		Action: int64(auth.List),
		Type:   resourceType,
		Token:  r.Header.Get("Authorization"),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return authClient.Authorize(ctx, ar)
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

	req.UserID = userID.GetValue()
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

	req.StreamOwner = userID.GetValue()
	return req, nil
}

func decodeReportRequest(_ context.Context, r *http.Request) (interface{}, error) {
	userID, err := authorize(r)
	if err != nil {
		return nil, err
	}

	req, err := decodeSearch(r)
	if err != nil {
		return nil, err
	}
	id := userID.GetValue()
	q := r.URL.Query().Get("type")
	req.owner = id
	switch q {
	case "bought":
		req.UserID = id
	case "sold":
		req.StreamOwner = id
	default:
		req.UserID = id
		req.StreamOwner = id
	}
	return req, nil
}

func decodeViewSubRequest(_ context.Context, r *http.Request) (interface{}, error) {
	ar := &authproto.AuthRequest{
		Action: int64(auth.Read),
		Type:   resourceType,
		Token:  r.Header.Get("Authorization"),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	userID, err := authClient.Authorize(ctx, ar)
	if err != nil {
		return nil, subscriptions.ErrUnauthorizedAccess
	}

	req := viewSubReq{
		userID:         userID.GetValue(),
		subscriptionID: bone.GetValue(r, "id"),
	}
	return req, nil
}

func decodeViewSubByUserAndStreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	userID := bone.GetValue(r, "ownerID")
	streamID := bone.GetValue(r, "streamID")

	req := viewSubByUserAndStreamReq{
		userID:   userID,
		streamID: streamID,
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

func encodeBinaryResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	if ar, ok := response.(apiRes); ok {
		for k, v := range ar.headers() {
			w.Header().Set(k, v)
		}

		w.WriteHeader(ar.code())

		if ar.empty() {
			return nil
		}
	}
	buff := response.(reportResponse)
	_, err := w.Write(buff)
	return err
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
	case subscriptions.ErrStreamAccess:
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
