package terms

// TermsLedger specifies access terms writer API.
type TermsLedger interface {
	// CreateTerms creates new terms for stream.
	CreateTerms(Terms) error

	// ValidateTerms validates terms for stream.
	ValidateTerms(Terms) (bool, error)
}
