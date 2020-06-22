package jwt

import "github.com/datapace/dproxy"

type token struct {
	url      string
	uid      string
	maxCalls int
}

var _ dproxy.Token = (*token)(nil)

func NewToken(uid, url string, maxCalls int) dproxy.Token {
	return &token{uid: uid, url: url, maxCalls: maxCalls}
}

func (t token) Url() string {
	return t.url
}

func (t token) Uid() string {
	return t.url
}

func (t token) MaxCalls() int {
	return t.maxCalls
}
