package subscriptions

// Proxy defines API for communicating with proxy.
type Proxy interface {
	// Registers new stream with url and given time to live. It returns
	// generated hash that is used for accessing given url.
	Register(ttl uint64, url string, maxCalls uint64, maxUnit, subscriptionID string) (string, error)
}
