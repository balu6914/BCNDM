package auth

import (
	"fmt"
	"monetasa/auth/fabric"
)

var _ Service = (*authService)(nil)

type authService struct {
	users  UserRepository
	hasher Hasher
	idp    IdentityProvider
	fabric fabric.Fabric
}

// New instantiates the domain service implementation.
func New(users UserRepository, hasher Hasher, idp IdentityProvider, fs fabric.Fabric) Service {
	return &authService{
		users:  users,
		hasher: hasher,
		idp:    idp,
		fabric: fs,
	}
}

func (ms *authService) Register(user User) error {
	hash, err := ms.hasher.Hash(user.Password)
	if err != nil {
		return ErrMalformedEntity
	}

	user.Password = hash
	err = ms.users.Save(user)
	if err != nil {
		return err
	}

	u, err := ms.users.OneByEmail(user.Email)
	if err != nil {
		return err
	}

	// Create New user in Fabric network calling fabric-ca
	newUser, err := fabric.CreateUser(u.ID.Hex(), u.Password, ms.fabric)
	if err != nil {
		return fmt.Errorf("Unable to create a user in the fabric-ca %v\n", err)
	}
	fmt.Printf("User created!: %v\n", newUser)

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
		return User{}, err
	}

	user, err := ms.users.OneById(id)
	if err != nil {
		return User{}, ErrUnauthorizedAccess
	}

	fs := ms.fabric
	// Get balance and update user
	balance, err := fabric.Balance(user.ID.Hex(), fs)
	if err != nil {
		return User{}, fmt.Errorf("Error fetching balance: %v\n", err)
	}
	user.Balance = balance

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
