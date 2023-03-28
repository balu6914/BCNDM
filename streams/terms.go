package streams

// TermsService contains API definition for Terms service.
type TermsService interface {
	// CreateTerms creates new terms on Terms service.
	CreateTerms(Stream) error
}
