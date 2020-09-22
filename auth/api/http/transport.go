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
		registerEndpoint(svc),
		decodeRegister,
		encodeResponse,
		opts...,
	))

	r.Get("/users", kithttp.NewServer(
		listUsersEndpoint(svc),
		decodeIdentity,
		encodeResponse,
		opts...,
	))

	r.Get("/users/non-partners", kithttp.NewServer(
		nonPartnersEndpoint(svc),
		decodeIdentity,
		encodeResponse,
		opts...,
	))

	r.Get("/users/:id", kithttp.NewServer(
		viewUserEndpoint(svc),
		decodeIdentity,
		encodeResponse,
		opts...,
	))

	r.Patch("/users/:id", kithttp.NewServer(
		updateUserEndpoint(svc),
		decodeUpdate,
		encodeResponse,
		opts...,
	))

	r.Post("/tokens", kithttp.NewServer(
		loginEndpoint(svc),
		decodeCredentials,
		encodeResponse,
		opts...,
	))

	// Policies API
	r.Post("/policies", kithttp.NewServer(
		addPolicyEndpoint(svc),
		decodePolicy,
		encodeResponse,
		opts...,
	))

	r.Get("/policies/:id", kithttp.NewServer(
		viewPolicyEndpoint(svc),
		decodeIdentity,
		encodeResponse,
		opts...,
	))

	r.Get("/policies", kithttp.NewServer(
		listPoliciesEndpoint(svc),
		decodeIdentity,
		encodeResponse,
		opts...,
	))

	r.Delete("/policies/:id", kithttp.NewServer(
		removePolicyEndpoint(svc),
		decodeIdentity,
		encodeResponse,
		opts...,
	))

	r.Put("/policies/:policy/users/:user", kithttp.NewServer(
		attachPolicyEndpoint(svc),
		decodeAttach,
		encodeResponse,
		opts...,
	))

	r.Delete("/policies/:policy/users/:user", kithttp.NewServer(
		detachPolicyEndpoint(svc),
		decodeAttach,
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
	if r.Header.Get("Content-Type") != contentType {
		return nil, errUnsupportedContentType
	}
	req := updateReq{
		id:  bone.GetValue(r, "id"),
		key: r.Header.Get("Authorization"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

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

func decodePolicy(_ context.Context, r *http.Request) (interface{}, error) {
	req := policyRequest{
		key: r.Header.Get("Authorization"),
	}
	if r.Header.Get("Content-Type") != contentType {
		return nil, errUnsupportedContentType
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeAttach(_ context.Context, r *http.Request) (interface{}, error) {
	req := attachReq{
		key: r.Header.Get("Authorization"),
	}
	req.policyID = bone.GetValue(r, "policy")
	req.userID = bone.GetValue(r, "user")
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
	case errInvalidEmail, errInvalidPassword, errInvalidFirstName,
		errInvalidLastName, errInvalidCompany, errInvalidAddress,
		errInvalidPolicyRules, errInvalidPolicyVersion:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	case auth.ErrUserAccountDisabled:
		w.WriteHeader(http.StatusForbidden)
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
