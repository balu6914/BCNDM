package proxy

type registerReq struct {
	TTL uint64 `json:"ttl"`
	URL string `json:"url"`
}
