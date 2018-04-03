package auth

import "errors"

var (
	// ErrConflict indicates usage of the existing email during account
	// registration.
	ErrConflict error = errors.New("email already taken")

	// ErrMalformedEntity indicates malformed entity specification (e.g.
	// invalid username or password).
	ErrMalformedEntity error = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess error = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound error = errors.New("non-existent entity")
)


// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// Register creates new user account. In case of the failed registration, a
	// non-nil error value is returned.
	Register(User) error

	// Update updates user account. In case of the failed update, a
	// non-nil error value is returned.
	Update(User) error

	// Deletedeletes user account. In case of the failed deletion, a
	// non-nil error value is returned.
	Delete(User) error

	// Login authenticates the user given its credentials. Successful
	// authentication generates new access token. Failed invocations are
	// identified by the non-nil error values in the response.
	Login(User) (string, error)

	// CanAccess determines whether the channel can be accessed using the
	// provided key.
	CanAccess(string, string) (string, error)
}
