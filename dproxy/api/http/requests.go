package http

import (
	"datapace/dproxy"
)

type requestCreateToken struct {
	URL      string `json:"url"`
	TTL      int    `json:"ttl"`
	MaxCalls int    `json:"max_calls"`
}

func (req requestCreateToken) validate() error {
	if req.URL == "" {
		return dproxy.ErrMalformedEntity
	}
	return nil
}
