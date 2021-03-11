package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/datapace/datapace"
	authproto "github.com/datapace/datapace/proto/auth"
	termsproto "github.com/datapace/datapace/proto/terms"
	"github.com/datapace/datapace/terms"
	"github.com/datapace/datapace/transactions"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
)

const (
	tokenType   = "token"
	contentType = "application/json"
)

var (
	errMalformedEntity        = errors.New("malformed entity")
	errUnauthorizedAccess     = errors.New("missing or invalid credentials provided")
	errUnsupportedContentType = errors.New("unsupported content type")
	termsClient               termsproto.TermsServiceClient
)

func MakeHandler(svc terms.Service, auth authproto.AuthServiceClient) http.Handler {

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()
	r.Get("/streams/:streamID/terms/:termsHash", kithttp.NewServer(
		validateEndpoint(svc),
		decodeValidateReq,
		encodeResponse,
		opts...,
	))
	r.GetFunc("/version", datapace.Version())
	return r
}

func decodeValidateReq(_ context.Context, r *http.Request) (interface{}, error) {
	req := validateTermsReq{
		streamID:  bone.GetValue(r, "streamID"),
		termsHash: bone.GetValue(r, "termsHash"),
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
