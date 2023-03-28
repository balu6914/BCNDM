package main

type termsReq struct {
	StreamID  string `json:"stream_id,omitempty"`
	TermsURL  string `json:"terms_url,omitempty"`
	TermsHash string `json:"terms_hash,omitempty"`
}
