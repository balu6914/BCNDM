package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"monetasa/subscriptions"
	"net/http"
)

const contentType = "application/json"

var _ subscriptions.Proxy = (*proxy)(nil)

type proxy struct {
	url string
}

// New receives proxy url and returns new proxy instance.
func New(url string) subscriptions.Proxy {
	return proxy{
		url: fmt.Sprintf("%s/api/register", url),
	}
}

func (p proxy) Register(ttl uint64, url string) (string, error) {
	req := registerReq{
		TTL: ttl,
		URL: url,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", subscriptions.ErrMalformedEntity
	}

	res, err := http.Post(p.url, contentType, bytes.NewReader(body))
	if err != nil {
		return "", subscriptions.ErrFailedCreateSub
	}
	defer res.Body.Close()

	var ur urlRes
	if err := json.NewDecoder(res.Body).Decode(&ur); err != nil {
		return "", subscriptions.ErrMalformedEntity
	}

	return ur.URL, nil
}
