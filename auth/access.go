package auth

// AccessControl contains API spec for access control.
type AccessControl interface {
	// PotentialPartners fetches and returns list of potential and actual
	// partners.
	PotentialPartners(string) ([]string, error)
}
