package auth

var _ Service = (*authService)(nil)

type authService struct {
	users  UserRepository
	hasher Hasher
	idp    IdentityProvider
	ts     TransactionsService
	access AccessControl
}

// New instantiates the domain service implementation.
func New(users UserRepository, hasher Hasher, idp IdentityProvider, ts TransactionsService, access AccessControl) Service {
	return &authService{
		users:  users,
		hasher: hasher,
		idp:    idp,
		ts:     ts,
		access: access,
	}
}

func (as *authService) Register(key string, user User) (string, error) {
	isAdmin, err := as.isAdmin(key)
	if err != nil {
		return "", ErrMalformedEntity
	}
	if !isAdmin {
		return "", ErrUnauthorizedAccess
	}
	hash, err := as.hasher.Hash(user.Password)
	if err != nil {
		return "", ErrMalformedEntity
	}

	user.Password = hash

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

func (as *authService) InitAdmin(user User) error {
	hash, err := as.hasher.Hash(user.Password)
	if err != nil {
		return ErrMalformedEntity
	}
	user.Password = hash
	_, err = as.users.OneByEmail(user.Email)
	if err != nil && err != ErrNotFound {
		return err
	}

	//User already exists so just return
	if err != ErrNotFound {
		return nil
	}
	id, err := as.users.Save(user)
	if err != nil {
		return err
	}
	if err := as.ts.CreateUser(id); err != nil {
		as.users.Remove(id)
		return err
	}
	return nil
}

func (as *authService) Login(user User) (string, error) {
	dbu, err := as.users.OneByEmail(user.Email)
	if err != nil {
		return "", ErrUnauthorizedAccess
	}

	if err := as.hasher.Compare(user.Password, dbu.Password); err != nil {
		return "", ErrUnauthorizedAccess
	}

	return as.idp.TemporaryKey(dbu.ID, dbu.Roles)
}

func (as *authService) Update(key string, user User) error {
	id, err := as.idp.Identity(key)
	if err != nil {
		return ErrUnauthorizedAccess
	}

	user.ID = id

	return as.users.Update(user)
}

func (as *authService) UpdatePassword(key string, old string, user User) error {
	id, err := as.idp.Identity(key)
	if err != nil {
		return ErrUnauthorizedAccess
	}
	user.ID = id
	// Validate current password with hashed record in DB
	dbu, err := as.users.OneByID(user.ID)
	if err != nil {
		return ErrNotFound
	}

	if err := as.hasher.Compare(old, dbu.Password); err != nil {
		return ErrMalformedEntity
	}

	hash, err := as.hasher.Hash(user.Password)

	if err != nil {
		return ErrMalformedEntity
	}
	user.Password = hash

	return as.users.Update(user)
}

func (as *authService) View(key, userID string) (User, error) {
	id, err := as.idp.Identity(key)
	if err != nil {
		return User{}, ErrUnauthorizedAccess
	}
	isAdmin, err := as.isAdmin(key)
	if err != nil {
		return User{}, ErrMalformedEntity
	}
	if id != userID && !isAdmin {
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
	isAdmin, err := as.isAdmin(key)
	if err != nil {
		return []User{}, ErrMalformedEntity
	}
	if !isAdmin {
		return []User{}, ErrUnauthorizedAccess
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

func (as *authService) Exists(id string) error {
	if _, err := as.users.OneByID(id); err != nil {
		return ErrNotFound
	}

	return nil
}

func (as *authService) isAdmin(key string) (bool, error) {
	roles, err := as.idp.Roles(key)
	if err != nil {
		return false, ErrMalformedEntity
	}
	isAdmin := false
	for _, role := range roles {
		if role == "admin" {
			isAdmin = true
		}
	}
	return isAdmin, nil
}
