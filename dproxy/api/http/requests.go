package http

import (
	"github.com/datapace/datapace/dproxy"
)

type requestCreateToken struct {
	URL      string `json:"url"`
	TTL      int    `json:"ttl"`
	MaxCalls int    `json:"max_calls,omitempty"`
	MaxUnit  string `json:"max_unit,omitempty"`
}

func (req requestCreateToken) validate() error {
	if req.URL == "" {
		return dproxy.ErrMalformedEntity
	}
	return nil
}
