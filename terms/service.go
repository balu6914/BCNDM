package terms

import (
	"errors"

	"gopkg.in/mgo.v2/bson"
)

var (
	// ErrConflict indicates usage of the existing terms id for the new terms.
	ErrConflict = errors.New("Terms ID already taken")

	// ErrMalformedEntity indicates malformed entity specification.
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrFailedCreateTerms indicates that creation of terms failed.
	ErrFailedCreateTerms = errors.New("failed to create terms")

	// ErrFailedFetchTermsURL indicates that  terms url is not reachable.
	ErrFailedFetchTermsURL = errors.New("failed to get terms url")
)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).

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

// Service specifies an API that must be fulfilled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// CreateTerms creates new Terms and returns Terms hash.
	CreateTerms(Terms) (string, error)

	// ValidateTerms validates existing Terms.
	ValidateTerms(Terms) (bool, error)
}
