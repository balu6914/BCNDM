package auth

var _ Service = (*authService)(nil)

type managerService struct {
	users    UserRepository
	hasher   Hasher
	idp      IdentityProvider
}

// New instantiates the domain service implementation.
func New(users UserRepository, clients ClientRepository, channels ChannelRepository, hasher Hasher, idp IdentityProvider) Service {
	return &managerService{
		users:    users,
		hasher:   hasher,
		idp:      idp,
	}
}

func (ms *managerService) Register(user User) error {
	hash, err := ms.hasher.Hash(user.Password)
	if err != nil {
		return ErrMalformedEntity
	}

	user.Password = hash
	return ms.users.Save(user)
}

func (ms *managerService) Login(user User) (string, error) {
	dbUser, err := ms.users.One(user.Email)
	if err != nil {
		return "", ErrUnauthorizedAccess
	}

	if err := ms.hasher.Compare(user.Password, dbUser.Password); err != nil {
		return "", ErrUnauthorizedAccess
	}

	return ms.idp.TemporaryKey(user.Email)
}

func (ms *managerService) Update(user User) error {
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

func (ms *managerService) View(key, id string) (User, error) {
	sub, err := ms.idp.Identity(key)
	if err != nil {
		return Client{}, err
	}

	if _, err := ms.users.One(sub); err != nil {
		return Client{}, ErrUnauthorizedAccess
	}

	return ms.clients.One(sub, id)
}

func (ms *managerService) List(key string) ([]User, error) {
	sub, err := ms.idp.Identity(key)
	if err != nil {
		return nil, err
	}

	if _, err := ms.users.One(sub); err != nil {
		return nil, ErrUnauthorizedAccess
	}

	return ms.clients.All(sub), nil
}

func (ms *managerService) Remove(key, id string) error {
	sub, err := ms.idp.Identity(key)
	if err != nil {
		return err
	}

	if _, err := ms.users.One(sub); err != nil {
		return ErrUnauthorizedAccess
	}

	return ms.clients.Remove(sub, id)
}

func (ms *managerService) Identity(key string) (string, error) {
	client, err := ms.idp.Identity(key)
	if err != nil {
		return "", err
	}

	return client, nil
}

func (ms *managerService) CanAccess(key, channel string) (string, error) {
	client, err := ms.idp.Identity(key)
	if err != nil {
		return "", err
	}

	if !ms.channels.HasClient(channel, client) {
		return "", ErrUnauthorizedAccess
	}

	return client, nil
}
