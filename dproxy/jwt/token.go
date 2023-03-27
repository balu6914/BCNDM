package jwt

import (
	"github.com/datapace/datapace/dproxy"
)

type token struct {
	url      string
	uid      string
	maxCalls int
	maxUnit  string
	subID    string
}

var _ dproxy.Token = (*token)(nil)

func NewToken(uid, url string, maxCalls int, maxUnit, subID string) dproxy.Token {
	return &token{uid: uid, url: url, maxCalls: maxCalls, maxUnit: maxUnit, subID: subID}
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

func (t token) MaxUnit() string {
	return t.maxUnit
}

func (t token) Subid() string {
	return t.subID
}
