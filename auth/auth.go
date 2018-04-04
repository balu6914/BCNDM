package auth

var _ Service = (*authService)(nil)

type authService struct {
	users    UserRepository
	hasher   Hasher
	idp      IdentityProvider
}

// New instantiates the domain service implementation.
func New(users UserRepository, clients ClientRepository, channels ChannelRepository, hasher Hasher, idp IdentityProvider) Service {
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

	if _, err := ms.users.One(sub); err != nil {
		return ErrUnauthorizedAccess
	}

	client.Owner = sub

	return ms.clients.Update(client)
}

func (ms *authService) View(key string, id string) (User, error) {
	sub, err := ms.idp.Identity(key)
	if err != nil {
		return Client{}, err
	}

	if _, err := ms.users.One(sub); err != nil {
		return Client{}, ErrUnauthorizedAccess
	}

	return ms.clients.One(sub, id)
}

func (ms *authService) List(key string) ([]User, error) {
	sub, err := ms.idp.Identity(key)
	if err != nil {
		return nil, err
	}

	if _, err := ms.users.One(sub); err != nil {
		return nil, ErrUnauthorizedAccess
	}

	return ms.clients.All(sub), nil
}

func (ms *authService) Remove(key string, id string) error {
	sub, err := ms.idp.Identity(key)
	if err != nil {
		return err
	}

	if _, err := ms.users.One(sub); err != nil {
		return ErrUnauthorizedAccess
	}

	return ms.clients.Remove(sub, id)
}

func (ms *authService) CanAccess(key, channel string) (string, error) {
	client, err := ms.idp.Identity(key)
	if err != nil {
		return "", err
	}

	if !ms.channels.HasClient(channel, client) {
		return "", ErrUnauthorizedAccess
	}

	return client, nil
}
