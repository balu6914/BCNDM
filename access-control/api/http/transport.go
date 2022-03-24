package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/datapace/datapace"
	"github.com/datapace/datapace/access-control"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var errUnsupportedContentType = errors.New("unsupported content type")

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc access.Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	r.Post("/access-requests", kithttp.NewServer(
		requestAccessEndpoint(svc),
		decodeRequestAccess,
		encodeResponse,
		opts...,
	))

	r.Put("/access-requests/:id/approve", kithttp.NewServer(
		approveAccessEndpoint(svc),
		decodeApproveAccess,
		encodeResponse,
		opts...,
	))

	r.Put("/access-requests/:id/revoke", kithttp.NewServer(
		revokeAccessEndpoint(svc),
		decodeRevokeAccess,
		encodeResponse,
		opts...,
	))

	r.Get("/access-requests/sent", kithttp.NewServer(
		listSentRequestsEndpoint(svc),
		decodeListAccessRequests,
		encodeResponse,
		opts...,
	))

	r.Get("/access-requests/received", kithttp.NewServer(
		listReceivedRequestsEndpoint(svc),
		decodeListAccessRequests,
		encodeResponse,
		opts...,
	))

	r.Put("/access/grant/:uid", kithttp.NewServer(
		grantAccessEndpoint(svc),
		decodeGrantRequest,
		encodeResponse,
		opts...,
	))

	r.GetFunc("/version", datapace.Version())
	r.Handle("/metrics", promhttp.Handler())

	return r
}

func decodeRequestAccess(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Header.Get("Content-Type") != contentType {
		return nil, errUnsupportedContentType
	}

	var req requestAccessReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	req.key = r.Header.Get("Authorization")

	return req, nil
}

func decodeApproveAccess(_ context.Context, r *http.Request) (interface{}, error) {
	req := approveAccessReq{
		key: r.Header.Get("Authorization"),
		id:  bone.GetValue(r, "id"),
	}

	return req, nil
}

func decodeRevokeAccess(_ context.Context, r *http.Request) (interface{}, error) {
	req := revokeAccessReq{
		key: r.Header.Get("Authorization"),
		id:  bone.GetValue(r, "id"),
	}

	return req, nil
}

func decodeListAccessRequests(_ context.Context, r *http.Request) (interface{}, error) {
	vals := bone.GetQuery(r, "state")
	if len(vals) > 1 {
		return nil, access.ErrMalformedEntity
	}

	state := access.State("")
	if len(vals) == 1 {
		state = access.State(vals[0])
	}

	req := listAccessRequestsReq{
		key:   r.Header.Get("Authorization"),
		state: state,
	}

	return req, nil
}

func decodeGrantRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := grantAccessReq{
		key:        r.Header.Get("Authorization"),
		dstUserIid: bone.GetValue(r, "uid"),
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
	case access.ErrMalformedEntity:
		w.WriteHeader(http.StatusBadRequest)
	case access.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case access.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case access.ErrConflict:
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
