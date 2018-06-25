package auth

import "errors"

var (
	// ErrConflict indicates usage of the existing email during account
	// registration.
	ErrConflict = errors.New("email already taken")

	// ErrMalformedEntity indicates malformed entity specification (e.g.
	// invalid username or password).
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrFetchingBalance indicates failure while fetching users balance.
	ErrFetchingBalance = errors.New("failed to fetch balance")
)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// Register creates new user account. In case of the failed registration, a
	// non-nil error value is returned.
	Register(User) error

	// Update updates user account. In case of the failed update, a
	// non-nil error value is returned.
	Update(string, User) error

	// ViewClient retrieves data about the client identified with the provided
	// ID, that belongs to the user identified by the provided key.
	View(string) (User, error)

	// Delete deletes user account. In case of the failed deletion, a
	// non-nil error value is returned.
	Delete(string) error

	// Login authenticates the user given its credentials. Successful
	// authentication generates new access token. Failed invocations are
	// identified by the non-nil error values in the response.
	Login(User) (string, error)

	// Identity retrieves Client ID for provided client token.
	Identify(string) (string, error)
}
