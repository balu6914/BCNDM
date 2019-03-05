package auth

var _ Service = (*authService)(nil)

type authService struct {
	users  UserRepository
	hasher Hasher
	idp    IdentityProvider
	ts     TransactionsService
	access AccessControl
	cipher Cipher
}

// New instantiates the domain service implementation.
func New(users UserRepository, hasher Hasher, idp IdentityProvider, ts TransactionsService, access AccessControl, cipher Cipher) Service {
	return &authService{
		users:  users,
		hasher: hasher,
		idp:    idp,
		ts:     ts,
		access: access,
		cipher: cipher,
	}
}

func (as *authService) Register(user User) error {
	hash, err := as.hasher.Hash(user.Password)
	if err != nil {
		return ErrMalformedEntity
	}

	user.Password = hash
	user, err = as.cipher.Encrypt(user)
	if err != nil {
		return ErrEncrypt
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

	return as.idp.TemporaryKey(dbu.ID)
}

func (as *authService) Update(key string, user User) error {
	id, err := as.idp.Identity(key)
	if err != nil {
		return ErrUnauthorizedAccess
	}

	user.ID = id

	user, err = as.cipher.Encrypt(user)
	if err != nil {
		return ErrEncrypt
	}

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

	user, err = as.cipher.Encrypt(user)
	if err != nil {
		return ErrEncrypt
	}

	return as.users.Update(user)
}

func (as *authService) View(key string) (User, error) {
	id, err := as.idp.Identity(key)
	if err != nil {
		return User{}, ErrUnauthorizedAccess
	}

	user, err := as.users.OneByID(id)
	if err != nil {
		return User{}, ErrUnauthorizedAccess
	}

	user, err = as.cipher.Decrypt(user)
	if err != nil {
		return User{}, ErrDecrypt
	}

	return user, nil
}

func (as *authService) ListUsers(key string) ([]User, error) {
	if _, err := as.idp.Identity(key); err != nil {
		return nil, ErrUnauthorizedAccess
	}

	users, err := as.users.AllExcept([]string{})
	if err != nil {
		return []User{}, err
	}

	for i, u := range users {
		users[i], err = as.cipher.Decrypt(u)
		if err != nil {
			return []User{}, ErrDecrypt
		}
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

	for i, u := range users {
		users[i], err = as.cipher.Decrypt(u)
		if err != nil {
			return []User{}, ErrDecrypt
		}
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
