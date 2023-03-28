package streams

// AccessControl contains API specification for checking access rights.
type AccessControl interface {
	// Partners fetches a list of partners who's resources can be viewed.
	Partners(string) ([]string, error)
}
