package auth

var _ Service = (*authService)(nil)

type authService struct {
	users  UserRepository
	hasher Hasher
	idp    IdentityProvider
	ts     TransactionsService
}

// New instantiates the domain service implementation.
func New(users UserRepository, hasher Hasher, idp IdentityProvider, ts TransactionsService) Service {
	return &authService{
		users:  users,
		hasher: hasher,
		idp:    idp,
		ts:     ts,
	}
}

func (ms *authService) Register(user User) error {
	hash, err := ms.hasher.Hash(user.Password)
	if err != nil {
		return ErrMalformedEntity
	}

	user.Password = hash
	id, err := ms.users.Save(user)
	if err != nil {
		return err
	}

	if err := ms.ts.CreateUser(id); err != nil {
		ms.users.Remove(id)
		return err
	}

	return nil

}

func (ms *authService) Login(user User) (string, error) {
	dbu, err := ms.users.OneByEmail(user.Email)
	if err != nil {
		return "", ErrUnauthorizedAccess
	}

	if err := ms.hasher.Compare(user.Password, dbu.Password); err != nil {
		return "", ErrUnauthorizedAccess
	}

	return ms.idp.TemporaryKey(dbu.ID)
}

func (ms *authService) Update(key string, user User) error {
	id, err := ms.idp.Identity(key)
	if err != nil {
		return ErrUnauthorizedAccess
	}

	user.ID = id

	return ms.users.Update(user)
}

func (ms *authService) UpdatePassword(key string, old string, user User) error {
	id, err := ms.idp.Identity(key)
	if err != nil {
		return ErrUnauthorizedAccess
	}
	user.ID = id
	// Validate current password with hashed record in DB
	dbu, err := ms.users.OneByID(user.ID)
	if err != nil {
		return ErrNotFound
	}

	if err := ms.hasher.Compare(old, dbu.Password); err != nil {
		return ErrMalformedEntity
	}

	hash, err := ms.hasher.Hash(user.Password)

	if err != nil {
		return ErrMalformedEntity
	}
	user.Password = hash

	return ms.users.Update(user)
}

func (ms *authService) View(key string) (User, error) {
	id, err := ms.idp.Identity(key)
	if err != nil {
		return User{}, ErrUnauthorizedAccess
	}

	user, err := ms.users.OneByID(id)
	if err != nil {
		return User{}, ErrUnauthorizedAccess
	}

	return user, nil
}

func (ms *authService) Identify(key string) (string, error) {
	id, err := ms.idp.Identity(key)
	if err != nil {
		return "", err
	}

	return id, nil
}
