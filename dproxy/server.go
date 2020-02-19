package dproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func MakeReverseProxy(target string) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		u, _ := url.Parse(target)
		proxy := httputil.NewSingleHostReverseProxy(u)
		req.URL.Host = u.Host
		req.URL.Scheme = u.Scheme
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = u.Host
		proxy.ServeHTTP(rw, req)
	}
}
