package http

import (
	"datapace/dproxy"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var errTokenNotFound = errors.New("token not found in url")

type ReverseProxy struct {
	svc dproxy.Service
	p   *httputil.ReverseProxy
}

func NewReverseProxy(svc dproxy.Service) *ReverseProxy {
	return &ReverseProxy{svc: svc, p: &httputil.ReverseProxy{}}
}

func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	rp.p.Director, err = rp.makeDirector(r)
	switch err {
	case nil:
	case dproxy.ErrInvalidToken:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		return
	case dproxy.ErrResourceNotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	rp.p.ServeHTTP(w, r)
}

func (rp *ReverseProxy) makeDirector(originalReq *http.Request) (func(r *http.Request), error) {
	t := originalReq.Header.Get("Authorization")
	if t == "" {
		return nil, errTokenNotFound
	}
	targetURL, err := rp.svc.GetTargetURL(t)
	if err != nil {
		return nil, err
	}
	return func(r *http.Request) {
		u, _ := url.Parse(targetURL)
		r.URL.Host = u.Host
		r.URL.Scheme = u.Scheme
		r.RequestURI = u.RequestURI()
		r.URL.RawQuery = u.RawQuery
		r.URL.RawPath = u.RawPath
		r.URL.Path = u.Path
		r.Host = u.Host
	}, nil
}
