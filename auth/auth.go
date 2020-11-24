package auth

import (
	"github.com/asaskevich/govalidator"
)

const version = "1.0.0"
const nbAttempet = 5

var _ Service = (*authService)(nil)

type authService struct {
	users    UserRepository
	hasher   Hasher
	idp      IdentityProvider
	ts       TransactionsService
	policies PolicyRepository
	access   AccessControl
}

// New instantiates the domain service implementation.
func New(users UserRepository, policies PolicyRepository, hasher Hasher, idp IdentityProvider, ts TransactionsService, access AccessControl) Service {
	return &authService{
		users:    users,
		hasher:   hasher,
		idp:      idp,
		ts:       ts,
		access:   access,
		policies: policies,
	}
}

func (as *authService) InitAdmin(user User, policies map[string]Policy) error {
	_, err := as.users.OneByEmail(user.Email)
	if err != ErrNotFound {
		return err
	}

	hash, err := as.hasher.Hash(user.Password)
	if err != nil {
		return ErrMalformedEntity
	}
	user.Password = hash

	uid, err := as.users.Save(user)
	if err != nil {
		return err
	}
	policy := policies["admin"]
	// Save policy. Use admin as an owner.
	// There is no need to check if the policy exists,
	// because Name and Owner fields are used as a compound index.
	policy.Owner = uid
	pid, err := as.policies.Save(policy)
	if err != nil {
		return err
	}
	if err := as.policies.Attach(pid, uid); err != nil {
		return err
	}
	up := policies["user"]
	up.Owner = uid
	if _, err := as.policies.Save(up); err != nil {
		return err
	}

	return nil
}

func (as *authService) Register(key string, user User) (string, error) {
	owner, err := as.Authorize(key, Create, user)

	if err != nil {
		return "", err
	}

	hash, err := as.hasher.Hash(user.Password)
	if err != nil {
		return "", ErrMalformedEntity
	}

	user.Password = hash
	// Add new password to history
	user.PasswordHistory = append(user.PasswordHistory, user.Password)

	if len(user.Policies) == 0 {
		policy, err := as.policies.OneByName(owner, "user")
		if err != nil {
			return "", err
		}
		user.Policies = append(user.Policies, policy)
	}

	id, err := as.users.Save(user)
	if err != nil {
		return "", err
	}

	if err := as.ts.CreateUser(id); err != nil {
		as.users.Remove(id)
		return "", err
	}

	return id, nil
}

func (as *authService) Login(user User) (string, error) {
	dbu, err := as.users.OneByEmail(user.Email)
	if err != nil {
		return "", ErrUnauthorizedAccess
	}

	if dbu.Attempt >= nbAttempet && dbu.Locked != false {
		dbu.Locked = true
		as.users.Update(dbu)
		return "", ErrUserAccountLocked
	}
	if err := as.hasher.Compare(user.Password, dbu.Password); err != nil {
		dbu.Attempt = dbu.Attempt + 1
		as.users.Update(dbu)
		return "", ErrUnauthorizedAccess
	}
	if dbu.Disabled {
		return "", ErrUserAccountDisabled
	}
	// Reset number of attempt with wrong password
	dbu.Attempt = 0
	as.users.Update(dbu)
	return as.idp.TemporaryKey(dbu.ID, dbu.Roles)
}

// Update updates existing user. Key(token) supplied needs to either have admin role or
// it needs to contain user.ID same with the user that is being updated (self update)
// if non empty password is supplied, then password is hashed and saved instead of clear text
func (as *authService) UpdateUser(key string, user User) error {
	u, err := as.users.OneByID(user.ID)
	if err != nil {
		return err
	}

	if _, err := as.Authorize(key, Update, u); err != nil {
		return ErrUnauthorizedAccess
	}
	// If password supplied, hash it
	if user.Password != "" {
		hash, err := as.hasher.Hash(user.Password)
		if err != nil {
			return ErrMalformedEntity
		}
		// Check if password is already used in the Last 5 passwords
		for _, hpassword := range u.PasswordHistory {
			if err := as.hasher.Compare(user.Password, hpassword); err != nil {
				return ErrUserPasswordHistory
			}
		}
		user.Password = hash
		user.PasswordHistory = u.PasswordHistory
		if len(user.PasswordHistory) > 5 {
			user.PasswordHistory = user.PasswordHistory[1:]
		}
		user.PasswordHistory = append(user.PasswordHistory, user.Password)
	}
	if user.ContactEmail != "" && !govalidator.IsEmail(user.ContactEmail) {
		return ErrMalformedEntity
	}
	return as.users.Update(user)
}

func (as *authService) ViewUser(key, userID string) (User, error) {
	if _, err := as.Authorize(key, Read, User{ID: userID}); err != nil {
		return User{}, ErrUnauthorizedAccess
	}
	user, err := as.users.OneByID(userID)
	if err != nil {
		return User{}, ErrUnauthorizedAccess
	}
	return user, nil
}

func (as *authService) ViewEmail(key string) (User, error) {
	id, err := as.idp.Identity(key)
	if err != nil {
		return User{}, ErrUnauthorizedAccess
	}
	user, err := as.users.OneByID(id)
	if err != nil {
		return User{}, ErrUnauthorizedAccess
	}
	return user, nil
}

func (as *authService) ListUsers(key string) ([]User, error) {
	if _, err := as.Authorize(key, List, User{}); err != nil {
		return []User{}, err
	}
	users, err := as.users.AllExcept([]string{})
	if err != nil {
		return []User{}, err
	}
	return users, nil
}

func (as *authService) ListNonPartners(key string) ([]User, error) {
	id, err := as.idp.Identity(key)
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}

	plist, err := as.access.PotentialPartners(id)
	if err != nil {
		return []User{}, err
	}

	plist = append(plist, id)

	users, err := as.users.AllExcept(plist)
	if err != nil {
		return []User{}, err
	}

	return users, nil
}

func (as *authService) Identify(key string) (string, error) {
	id, err := as.idp.Identity(key)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (as *authService) Authorize(key string, action Action, resource Resource) (string, error) {
	id, err := as.idp.Identity(key)
	if err != nil {
		return "", ErrUnauthorizedAccess
	}

	user, err := as.users.OneByID(id)
	if err != nil {
		return "", err
	}

	if user.Disabled {
		return "", ErrUserAccountDisabled
	}

	policies := user.Policies
	for _, p := range policies {
		if p.Validate(user, action, resource) {
			return id, nil
		}
	}

	return "", ErrUnauthorizedAccess
}

func (as *authService) Exists(id string) error {
	if _, err := as.users.OneByID(id); err != nil {
		return ErrNotFound
	}

	return nil
}

func (as *authService) AddPolicy(key string, policy Policy) (string, error) {
	id, err := as.Authorize(key, Create, policy)
	if err != nil {
		return "", ErrUnauthorizedAccess
	}
	policy.Owner = id
	return as.policies.Save(policy)
}

func (as *authService) ViewPolicy(key, id string) (Policy, error) {
	_, err := as.Authorize(key, Read, Policy{})
	if err != nil {
		return Policy{}, ErrUnauthorizedAccess
	}
	return as.policies.OneByID(id)
}

func (as *authService) ListPolicies(key string) ([]Policy, error) {
	owner, err := as.Authorize(key, Read, Policy{})
	if err != nil {
		return []Policy{}, ErrUnauthorizedAccess
	}
	return as.policies.List(owner)
}

func (as *authService) RemovePolicy(key string, id string) error {
	id, err := as.idp.Identity(key)
	if err != nil {
		return ErrUnauthorizedAccess
	}

	return as.policies.Remove(id)
}

func (as *authService) AttachPolicy(key, policyID, userID string) error {
	_, err := as.idp.Identity(key)
	if err != nil {
		return ErrUnauthorizedAccess
	}
	return as.policies.Attach(policyID, userID)
}

func (as *authService) DetachPolicy(key, policyID, userID string) error {
	_, err := as.idp.Identity(key)
	if err != nil {
		return ErrUnauthorizedAccess
	}
	return as.policies.Detach(policyID, userID)
}
