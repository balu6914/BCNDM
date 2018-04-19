package auth

var _ Service = (*authService)(nil)

type authService struct {
	users    UserRepository
	hasher   Hasher
	idp      IdentityProvider
}

// New instantiates the domain service implementation.
func New(users UserRepository, hasher Hasher, idp IdentityProvider) Service {
	return &authService{
		users:    users,
		hasher:   hasher,
		idp:      idp,
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
	dbUser, err := ms.users.One(user.Email)
	if err != nil {
		return "", ErrUnauthorizedAccess
	}

	if err := ms.hasher.Compare(user.Password, dbUser.Password); err != nil {
		return "", ErrUnauthorizedAccess
	}

	return ms.idp.TemporaryKey(user.Email)
}

func (ms *authService) Update(key string, user User) error {
	sub, err := ms.idp.Identity(key)
	if err != nil {
		return err
	}

	u, err := ms.users.One(sub)
	if  err != nil {
		return ErrUnauthorizedAccess
	}

	if u.Email != user.Email {
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
	sub, err := ms.idp.Identity(key)
	if err != nil {
		return User{}, err
	}

	user, err := ms.users.One(sub);
	if err != nil {
		return User{}, ErrUnauthorizedAccess
	}

	return user, nil
}

func (ms *authService) List(key string) ([]User, error) {
	sub, err := ms.idp.Identity(key)
	if err != nil {
		return nil, err
	}

	if _, err := ms.users.One(sub); err != nil {
		return nil, ErrUnauthorizedAccess
	}

	return ms.users.All()
}

func (ms *authService) Delete(key string) error {
	sub, err := ms.idp.Identity(key)
	if err != nil {
		return err
	}

	user, err := ms.users.One(sub)
	if err != nil {
		return ErrUnauthorizedAccess
	}

	return ms.users.Remove(user.Email)
}

func (ms *authService) Identity(key string) (string, error) {
	user, err := ms.idp.Identity(key)
	if err != nil {
		return "", err
	}

	return user, nil
}

func (ms *authService) CanAccess(key, channel string) (string, error) {
	client, err := ms.idp.Identity(key)
	if err != nil {
		return "", err
	}

	return client, nil
}
