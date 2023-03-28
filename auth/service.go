package auth

import "errors"

var (
	// ErrConflict indicates usage of the existing email during account
	// registration.
	ErrConflict = errors.New("unique index violation")

	// ErrMalformedEntity indicates malformed entity specification (e.g.
	// invalid username or password).
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrUserAccountDisabled indicates user account being disabled
	ErrUserAccountDisabled = errors.New("user account disabled")

	// ErrUserAccountLocked indicates user account being locked
	ErrUserAccountLocked = errors.New("user account locked")

	// ErrUserPasswordHistory indicates user password exists already in last 5 password history
	ErrUserPasswordHistory = errors.New("password exists in the last 5 password history")

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
	InitAdmin(User, map[string]Policy) error

	// Login authenticates the user given its credentials. Successful
	// authentication generates new access token. Failed invocations are
	// identified by the non-nil error values in the response.
	Login(User) (string, error)

	// Update updates user account. In case of the failed update, a
	// non-nil error value is returned.
	UpdateUser(string, User) error

	// RecoverPassword sends recovery URL to user's email.
	// If there's no user with such email in database, a non-nil error is returned.
	RecoverPassword(string) error

	// ValidateRecoveryToken checks if token is valid and not expired.
	// In case that such token doesn't exist or is expired, a non-nil error value is returned.
	ValidateRecoveryToken(string, string) error

	// UpdatePassword changes the password of a user that the password recovery token is assigned to.
	// Sets the password provided via request along with the valid token.
	// In case that the token is expired/invalid or password doesn't match length criteria, a non-nil error value is returned.
	UpdatePassword(string, string, string) error

	// ViewUser retrieves data about the client identified with the provided ID
	// Key provided must have same ID (user viewing his own data) or admin role.
	ViewUser(string, string) (User, error)

	// ViewUserById retrieves data about the user by its id.
	ViewUserById(id string) (User, error)

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

	// Exists checks if user with specified id exists. If it doesn't then error
	// is returned.
	Exists(string) error

	// Authorize authorizes user with provided token for executing
	// given action over given resource.
	Authorize(key string, action Action, resource Resource) (string, error)

	// AddPolicy creates a new Policy.
	AddPolicy(key string, policy Policy) (string, error)

	// ViewPolicy returns a policy with the given ID.
	ViewPolicy(key, id string) (Policy, error)

	// ListPolicies lists all the policies that belong to the given user.
	ListPolicies(key string) ([]Policy, error)

	// RemovePolicy removes an existing policy.
	RemovePolicy(key, id string) error

	// AttachPolicy attaches a policy to the user.
	AttachPolicy(key, policyID, userID string) error

	// DetachPolicy removes policy from the user.
	DetachPolicy(key, policyID, userID string) error
}
