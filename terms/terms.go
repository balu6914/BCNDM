package terms

import (
	"gopkg.in/mgo.v2/bson"
)

// Terms represents users terms for his stream.
type Terms struct {
	ID        bson.ObjectId `json:"id,omitempty"`
	StreamID  string        `json:"stream_id,omitempty"`
	TermsURL  string        `json:"terms_url,omitempty"`
	TermsHash string        `json:"terms_hash,omitempty"`
}

// Validate returns an error if representation is invalid.
func (sub *Terms) Validate() error {
	return nil
}

// Page represents paged result for list response.
type Page struct {
	Page    uint64  `json:"page"`
	Limit   uint64  `json:"limit"`
	Total   uint64  `json:"total"`
	Content []Terms `json:"content"`
}

// TermsRepository specifies a Terms persistence API.
type TermsRepository interface {
	// Save persists terms.
	Save(Terms) (string, error)

	// One retrieves Terms by its ID.
	One(string) (Terms, error)
}

// TermsLedger specifies access terms writer API.
type TermsLedger interface {
	// CreateTerms creates new terms for stream.
	CreateTerms(Terms) error

	// ValidateTerms validates terms for stream.
	ValidateTerms(Terms) (bool, error)
}
