package proxy

type registerReq struct {
	TTL            uint64 `json:"ttl"`
	URL            string `json:"url"`
	MaxCalls       uint64 `json:"max_calls,omitempty"`
	MaxUnit        string `json:"max_unit,omitempty"`
	SubscriptionID string `json:"subscription_id"`
}
