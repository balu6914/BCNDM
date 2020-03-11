package http

import (
	"context"
	"datapace/access-control"
	"datapace/dproxy"
	"encoding/json"
	"errors"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"io"
	"net/http"
)

var errUnsupportedContentType = errors.New("unsupported content type")

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc dproxy.Service, rp *ReverseProxy, fs *FsProxy) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}
	r := bone.New()
	r.GetFunc("/fs", fs.GetFile)
	r.PutFunc("/fs", fs.PutFile)
	r.Get("/http", rp)
	r.Post("/api/register", kithttp.NewServer(createTokenEndpoint(svc),
		decodeCreateToken,
		encodeResponse,
		opts...))
	return r
}

func decodeCreateToken(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Header.Get("Content-Type") != contentType {
		return nil, errUnsupportedContentType
	}
	var req requestCreateToken
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", contentType)
	if _, ok := response.(createTokenRes); ok {
		w.WriteHeader(http.StatusCreated)
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
