package auth

var _ Service = (*authService)(nil)

type authService struct {
	users  UserRepository
	hasher Hasher
	idp    IdentityProvider
}

// New instantiates the domain service implementation.
func New(users UserRepository, hasher Hasher, idp IdentityProvider) Service {
	return &authService{
		users:  users,
		hasher: hasher,
		idp:    idp,
	}
}

func (ms *authService) Register(user User) error {
	hash, err := ms.hasher.Hash(user.Password)
	if err != nil {
		return ErrMalformedEntity
	}

	user.Password = hash
	return ms.users.Save(user)
}

func (ms *authService) Login(user User) (string, error) {
	dbu, err := ms.users.OneByEmail(user.Email)
	if err != nil {
		return "", ErrUnauthorizedAccess
	}

	if err := ms.hasher.Compare(user.Password, dbu.Password); err != nil {
		return "", ErrUnauthorizedAccess
	}

	return ms.idp.TemporaryKey(dbu.ID.Hex())
}

func (ms *authService) Update(key string, user User) error {
	id, err := ms.idp.Identity(key)
	if err != nil {
		return err
	}

	u, err := ms.users.OneById(id)
	if err != nil {
		return ErrUnauthorizedAccess
	}
	if u.ID.Hex() != user.ID.Hex() {
		return ErrUnauthorizedAccess
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
		return User{}, err
	}

	user, err := ms.users.OneById(id)
	if err != nil {
		return User{}, ErrUnauthorizedAccess
	}

	return user, nil
}

func (ms *authService) Delete(key string) error {
	id, err := ms.idp.Identity(key)
	if err != nil {
		return err
	}

	user, err := ms.users.OneById(id)
	if err != nil {
		return ErrUnauthorizedAccess
	}

	return ms.users.Remove(user.ID.Hex())
}

func (ms *authService) Identity(key string) (string, error) {
	id, err := ms.idp.Identity(key)
	if err != nil {
		return "", err
	}

	return id, nil
}
