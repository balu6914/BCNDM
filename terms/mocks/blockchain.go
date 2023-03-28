package mocks

import (
	t "github.com/datapace/datapace/terms"
)

type mockTermsLedger struct {
	terms map[string]string
}

// NewTermsLedger returns mock instance of terms ledger.
func NewTermsLedger() t.TermsLedger {
	return mockTermsLedger{terms: map[string]string{}}
}

// CreateTerms creates new Terms.
func (tl mockTermsLedger) CreateTerms(terms t.Terms) error {
	tl.terms[terms.StreamID] = terms.TermsHash
	return nil
}

// ValidateTerms validates existing Terms.
func (tl mockTermsLedger) ValidateTerms(terms t.Terms) (bool, error) {
	if tl.terms[terms.StreamID] == terms.TermsHash {
		return true, nil
	}
	return false, nil
}
