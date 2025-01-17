package auth

// IdentityProvider specifies an API for identity management via security
// tokens.
type IdentityProvider interface {
	// TemporaryKey generates the temporary access token.
	TemporaryKey(string, string) (string, error)

	// Identity extracts the entity identifier given its secret key.
	Identity(string) (string, error)

	// Role extracts the user roles given its secret key.
	Role(string) (string, error)
}
