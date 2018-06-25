package auth

var _ Service = (*authService)(nil)

type authService struct {
	users  UserRepository
	hasher Hasher
	idp    IdentityProvider
	fabric FabricNetwork
}

// New instantiates the domain service implementation.
func New(users UserRepository, hasher Hasher, idp IdentityProvider, fn FabricNetwork) Service {
	return &authService{
		users:  users,
		hasher: hasher,
		idp:    idp,
		fabric: fn,
	}
}

func (ms *authService) Register(user User) error {
	hash, err := ms.hasher.Hash(user.Password)
	if err != nil {
		return ErrMalformedEntity
	}

	user.Password = hash
	if err = ms.users.Save(user); err != nil {
		return err
	}

	u, err := ms.users.OneByEmail(user.Email)
	if err != nil {
		return err
	}

	// Create New user in Fabric network calling fabric-ca
	return ms.fabric.CreateUser(u.ID.Hex(), u.Password)
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
		return ErrUnauthorizedAccess
	}

	u, err := ms.users.OneByID(id)
	if err != nil {
		return ErrUnauthorizedAccess
	}

	hash, err := ms.hasher.Hash(user.Password)
	if err != nil {
		return ErrMalformedEntity
	}
	user.Password = hash
	user.ID = u.ID

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

	// Get balance and update user
	balance, err := ms.fabric.Balance(user.ID.Hex())
	if err != nil {
		return User{}, ErrFetchingBalance
	}
	user.Balance = balance

	return user, nil
}

func (ms *authService) Delete(key string) error {
	id, err := ms.idp.Identity(key)
	if err != nil {
		return ErrUnauthorizedAccess
	}

	user, err := ms.users.OneByID(id)
	if err != nil {
		return ErrUnauthorizedAccess
	}

	return ms.users.Remove(user.ID.Hex())
}

func (ms *authService) Identify(key string) (string, error) {
	id, err := ms.idp.Identity(key)
	if err != nil {
		return "", err
	}

	return id, nil
}
