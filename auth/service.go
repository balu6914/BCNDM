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

	// ErrUserAccountDisabled indicates user account being disabled
	ErrUserAccountDisabled = errors.New("user account disabled")

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
	Register(string, User) (string, error)

	// InitAdmin creates admin account if it does not exist already. In case of the failed creation, a
	// non-nil error value is returned.
	InitAdmin(User) error

	// Login authenticates the user given its credentials. Successful
	// authentication generates new access token. Failed invocations are
	// identified by the non-nil error values in the response.
	Login(User) (string, error)

	// Update updates user account. In case of the failed update, a
	// non-nil error value is returned.
	Update(string, User) error

	// View retrieves data about the client identified with the provided ID
	// Key provided must have same ID (user viewing his own data) or admin role.
	View(string, string) (User, error)

	// ViewEmail provides backwards compatibility for grpc which doesn't support authorization at the moment
	// It retrieves data about the client identified with the provided
	// ID, that belongs to the user identified by the provided key.
	ViewEmail(string) (User, error)

	// Identity retrieves Client ID for provided client token.
	Identify(string) (string, error)

	// ListUsers retrieves list of all system users.
	ListUsers(string) ([]User, error)

	// ListNonPartners retrieves list of all users that are not potential or
	// actual partners.
	ListNonPartners(string) ([]User, error)

	// Exists checkes if user with specified id exists. If it doesn't then error
	// is returned.
	Exists(string) error
}
