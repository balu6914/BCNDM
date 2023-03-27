package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/datapace/datapace"

	"github.com/datapace/datapace/access-control"
	"github.com/datapace/datapace/dproxy"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
)

const tokenoutput = "token"
const urloutput = "url"

var errUnsupportedContentType = errors.New("unsupported content type")

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc dproxy.Service, rp *ReverseProxy, fs *FsProxy, dProxyRootUrl string) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}
	r := bone.New()
	r.GetFunc(fs.PathPrefix, fs.GetFile)
	r.GetFunc(fs.PathPrefix+"/:token", fs.GetFile)
	r.PutFunc(fs.PathPrefix, fs.PutFile)
	r.PutFunc(fs.PathPrefix+"/:token", fs.PutFile)
	r.PostFunc(fs.PathPrefix, fs.PostFile)
	r.PostFunc(fs.PathPrefix+"/:token", fs.PostFile)
	r.Get(rp.PathPrefix+"/", http.StripPrefix(rp.PathPrefix+"/", rp))
	r.Get(rp.PathPrefix, rp)
	r.Post("/api/token", kithttp.NewServer(createTokenEndpoint(svc, tokenoutput, dProxyRootUrl, fs.PathPrefix, rp.PathPrefix),
		decodeCreateToken,
		encodeTokenResponse,
		opts...))
	r.Post("/api/register", kithttp.NewServer(createTokenEndpoint(svc, urloutput, dProxyRootUrl, fs.PathPrefix, rp.PathPrefix),
		decodeCreateUrl,
		encodeUrlResponse,
		opts...))
	r.GetFunc("/version", datapace.Version())
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

func decodeCreateUrl(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Header.Get("Content-Type") != contentType {
		return nil, errUnsupportedContentType
	}
	var req requestCreateToken
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeTokenResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", contentType)
	if _, ok := response.(createTokenRes); ok {
		w.WriteHeader(http.StatusCreated)
	}
	return json.NewEncoder(w).Encode(response)
}

func encodeUrlResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", contentType)
	if _, ok := response.(createUrlRes); ok {
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
