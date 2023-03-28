package access

// AuthService contains API for fetching user related data.
type AuthService interface {
	// Identifies user by the provided key.
	Identify(string) (string, error)
	// Exists checks if user exists and returns not found error if not.
	Exists(string) error
}
