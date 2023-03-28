package terms

// TermsRepository specifies a Terms persistence API.
type TermsRepository interface {
	// Save persists terms.
	Save(Terms) (string, error)

	// One retrieves Terms by its ID.
	One(string) (Terms, error)
}
