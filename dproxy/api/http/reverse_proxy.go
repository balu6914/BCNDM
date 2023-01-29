package http

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/datapace/datapace/dproxy"
	log "github.com/datapace/datapace/logger"
)

var errTokenNotFound = errors.New("token not found in url")

type ReverseProxy struct {
	svc        dproxy.Service
	p          *httputil.ReverseProxy
	logger     log.Logger
	logPrefix  string
	PathPrefix string
}

func NewReverseProxy(svc dproxy.Service, pathPrefix string, tlsSkipVerify bool, logger log.Logger) *ReverseProxy {
	return &ReverseProxy{
		svc: svc,
		p: &httputil.ReverseProxy{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: tlsSkipVerify,
				},
			},
		},
		PathPrefix: pathPrefix,
		logger:     logger,
		logPrefix:  "rp",
	}
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
	case dproxy.ErrTokenExpired:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		return
	case dproxy.ErrQuotaExceeded:
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(http.StatusText(http.StatusTooManyRequests)))
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
	//if there is no token in authorization header, try token in the url
	if t == "" {
		t = originalReq.URL.Path
	}
	targetURL, err := rp.svc.GetTargetURL(t)
	if err != nil {
		return nil, err
	}
	rp.logger.Info(fmt.Sprintf("%s: proxying request from %s to %s", rp.logPrefix, originalReq.RemoteAddr, targetURL))
	return func(r *http.Request) {
		u, _ := url.Parse(targetURL)
		r.URL.Host = u.Host
		r.URL.Scheme = u.Scheme
		r.RequestURI = u.RequestURI()
		r.URL.RawQuery = u.RawQuery
		r.URL.RawPath = u.RawPath
		r.URL.Path = u.Path
		r.Host = u.Host
		anonymizeRequest(r)
	}, nil
}

func anonymizeRequest(req *http.Request) {
	req.Header["X-Forwarded-For"] = nil
	req.Header["X-Original-Forwarded-For"] = nil
	req.Header["X-Real-Ip"] = nil
	req.Header["User-Agent"] = nil
}
