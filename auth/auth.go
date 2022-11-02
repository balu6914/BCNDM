package auth

import (
	"fmt"
	"io"
	"math/rand"
	"regexp"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/asaskevich/govalidator"
)

const (
	version           = "1.0.0"
	nbAttempet        = 5
	numbers           = `[0-9]{1}`
	lowerLetters      = `[a-z]{1}`
	capitalLetters    = `[A-Z]{1}`
	symbol            = `[!@#~$%^&*()+|_]{1}`
	minPasswordLength = 9
)

var (
	// ErrPassLength indicates that password lenght is lower-than minPasswordLength
	ErrPassLength = fmt.Errorf("password length is lower than %d", minPasswordLength)
	// ErrPassContainNum indicates that password don't contain a number
	ErrPassContainNum = fmt.Errorf("password must contain at least a number")
	// ErrPassContainLowCase indicates that password don't contain a character between a and z
	ErrPassContainLowCase = fmt.Errorf("password must contain a character between a and z")
	// ErrPassContainUpCase indicates that password don't contain a character between a and z
	ErrPassContainUpCase = fmt.Errorf("password must contain a character between A and Z")
	// ErrPassContainSymbol indicates that password don't contain a symbol
	ErrPassContainSymbol = fmt.Errorf("password must contain a symbol")
	// ErrPasswordRecoveryExpired indicates that password recovery time of 30 minutes has expired
	ErrPasswordRecoveryExpired = fmt.Errorf("password recovery time expired")
	// ErrMailNotSent that there was a problem with mail server while sending the password recovery email
	ErrMailNotSent = fmt.Errorf("email transport error")
)

var _ Service = (*authService)(nil)

type entropy struct {
	mu sync.Mutex
	t  time.Time
	r  io.Reader
}

func (e *entropy) Read(p []byte) (n int, err error) {
	e.mu.Lock()
	n, err = e.r.Read(p)
	e.mu.Unlock()
	return
}

type authService struct {
	users    UserRepository
	recovery RecoveryTokenProvider
	mailsvc  MailService
	hasher   Hasher
	idp      IdentityProvider
	ts       TransactionsService
	policies PolicyRepository
	access   AccessControl
	entropy  *entropy
}

var symbols = []rune("01233456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(b)
}

// New instantiates the domain service implementation.
func New(users UserRepository, policies PolicyRepository, hasher Hasher, idp IdentityProvider, ts TransactionsService, access AccessControl, recovery RecoveryTokenProvider, mailsvc MailService) Service {
	t := time.Now()
	mt := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	e := &entropy{r: mt, t: t}
	return &authService{
		users:    users,
		recovery: recovery,
		mailsvc:  mailsvc,
		hasher:   hasher,
		idp:      idp,
		ts:       ts,
		access:   access,
		policies: policies,
		entropy:  e,
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
	policy := policies[AdminRole]
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

	// Remove Admin policy because it's already saved.
	delete(policies, AdminRole)
	for _, p := range policies {
		p.Owner = uid
		if _, err := as.policies.Save(p); err != nil {
			return err
		}
	}

	return nil
}

// CheckPasswordLevel validate all password regex
func CheckPasswordLevel(ps string) error {
	if len(ps) < minPasswordLength {
		return ErrPassLength
	}
	if matched, err := regexp.MatchString(numbers, ps); !matched || err != nil {
		return ErrPassContainNum
	}
	if matched, err := regexp.MatchString(lowerLetters, ps); !matched || err != nil {
		return ErrPassContainLowCase
	}
	if matched, err := regexp.MatchString(capitalLetters, ps); !matched || err != nil {
		return ErrPassContainUpCase
	}
	if matched, err := regexp.MatchString(symbol, ps); !matched || err != nil {
		return ErrPassContainSymbol
	}
	return nil
}

func (as *authService) Register(key string, user User) (string, error) {
	_, err := as.Authorize(key, Create, user)

	if err != nil {
		return "", err
	}
	if err := CheckPasswordLevel(user.Password); err != nil {
		return "", err
	}

	hash, err := as.hasher.Hash(user.Password)
	if err != nil {
		return "", ErrMalformedEntity
	}

	user.Password = hash
	// Add new password to history
	user.PasswordHistory = append(user.PasswordHistory, user.Password)

	// If there is no User policy, attach one depending on the role.
	if len(user.Policies) == 0 {
		policy, err := as.policies.OneByName(user.Role)
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

	if dbu.Locked {
		return "", ErrUserAccountLocked
	}
	if dbu.Attempt >= nbAttempet && dbu.Locked != true {
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
	return as.idp.TemporaryKey(dbu.ID, dbu.Role)
}

func (as *authService) checkAndHashPassword(storedUser User, newUser *User) error {
	hash, err := as.hasher.Hash(newUser.Password)
	if err != nil {
		return ErrMalformedEntity
	}
	// Check if password is already used in the Last 5 passwords
	for _, hpassword := range storedUser.PasswordHistory {
		if err := as.hasher.Compare(newUser.Password, hpassword); err == nil {
			return ErrUserPasswordHistory
		}
	}
	newUser.Password = hash
	newUser.PasswordHistory = storedUser.PasswordHistory
	if len(newUser.PasswordHistory) > 5 {
		newUser.PasswordHistory = newUser.PasswordHistory[1:]
	}
	newUser.PasswordHistory = append(newUser.PasswordHistory, newUser.Password)

	return nil
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
	if user.Role != "" {
		p, err := as.policies.OneByName(user.Role)
		if err != nil {
			return err
		}
		user.Policies = []Policy{p}
	}

	// If password supplied, hash it and check against latest 5 passwords
	if user.Password != "" {
		err := as.checkAndHashPassword(u, &user)
		if err != nil {
			return err
		}
	}
	if user.ContactEmail != "" && !govalidator.IsEmail(user.ContactEmail) {
		return ErrMalformedEntity
	}
	return as.users.Update(user)
}

func (as *authService) RecoverPassword(email string) error {
	user, err := as.users.OneByEmail(email)
	if err != nil {
		return err
	}

	secret := randomString(32)
	user.PasswordResetSecret = secret
	tokenString, err := as.recovery.CreateTokenString(user.ID, secret)
	if err != nil {
		return err
	}

	if updateErr := as.users.Update(user); updateErr != nil {
		return updateErr
	}

	templateData := map[string]interface{}{
		"Token": tokenString,
		"ID":    user.ID,
	}
	emailSubject := "Datapace password recovery."

	if mailErr := as.mailsvc.SendRecoveryEmail(email, emailSubject, templateData); mailErr != nil {
		return mailErr
	}

	return nil
}

func (as *authService) ValidateRecoveryToken(token string, id string) error {
	user, err := as.users.OneByID(id)
	if err != nil {
		return ErrNotFound
	}
	_, parseErr := as.recovery.ParseToken(token, user.PasswordResetSecret)
	if parseErr != nil {
		return parseErr
	}

	return nil
}

func (as *authService) UpdatePassword(token string, id string, password string) error {
	storedUser, err := as.users.OneByID(id)
	if err != nil {
		return ErrNotFound
	}
	_, parseErr := as.recovery.ParseToken(token, storedUser.PasswordResetSecret)
	if parseErr != nil {
		return parseErr
	}
	newUser := storedUser
	newUser.Password = password
	if checkErr := as.checkAndHashPassword(storedUser, &newUser); checkErr != nil {
		return checkErr
	}
	newUser.PasswordResetSecret = ""
	return as.users.Update(newUser)
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

func (as *authService) ViewUserById(id string) (User, error) {
	user, err := as.users.OneByID(id)
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
	return as.ViewUserById(id)
}

func (as *authService) ListUsers(key string) ([]User, error) {
	if _, err := as.Authorize(key, List, User{}); err != nil {
		return []User{}, err
	}

	role, err := as.idp.Role(key)
	if err != nil {
		return []User{}, err
	}

	var filters AdminFilters
	switch role {
	case AdminRole:
		filters.Roles = []string{UserRole, AdminUserRole, AdminWalletRole}
		filters.Locked = true
		filters.Disabled = true
	case AdminUserRole:
		filters.Roles = []string{UserRole}
		filters.Locked = true
		filters.Disabled = true
	case AdminWalletRole:
		filters.Roles = []string{UserRole}
	default:
		return []User{}, nil
	}

	users, err := as.users.RetrieveAll(filters)
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
