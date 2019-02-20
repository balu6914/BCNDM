package auth

var _ Service = (*authService)(nil)

type authService struct {
	users    UserRepository
	hasher   Hasher
	idp      IdentityProvider
	ts       TransactionsService
	accesses AccessRequestRepository
}

// New instantiates the domain service implementation.
func New(users UserRepository, hasher Hasher, idp IdentityProvider, ts TransactionsService, accesses AccessRequestRepository) Service {
	return &authService{
		users:    users,
		hasher:   hasher,
		idp:      idp,
		ts:       ts,
		accesses: accesses,
	}
}

func (as *authService) Register(user User) error {
	hash, err := as.hasher.Hash(user.Password)
	if err != nil {
		return ErrMalformedEntity
	}

	user.Password = hash
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

func (as *authService) View(key string) (User, error) {
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

func (as *authService) List(key string) ([]User, error) {
	if _, err := as.idp.Identity(key); err != nil {
		return nil, ErrUnauthorizedAccess
	}

	return as.users.List()
}

func (as *authService) RequestAccess(key, partner string) (string, error) {
	if _, err := as.users.OneByID(partner); err != nil {
		return "", ErrNotFound
	}

	id, err := as.idp.Identity(key)
	if err != nil {
		return "", ErrUnauthorizedAccess
	}

	return as.accesses.RequestAccess(id, partner)
}

func (as *authService) ListSentAccessRequests(key string, state State) ([]AccessRequest, error) {
	id, err := as.idp.Identity(key)
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}

	return as.accesses.ListSentAccessRequests(id, state)
}

func (as *authService) ListReceivedAccessRequests(key string, state State) ([]AccessRequest, error) {
	id, err := as.idp.Identity(key)
	if err != nil {
		return nil, ErrUnauthorizedAccess
	}

	return as.accesses.ListReceivedAccessRequests(id, state)
}

func (as *authService) ListPartners(id string) ([]string, error) {
	requests, err := as.accesses.ListSentAccessRequests(id, Approved)
	if err != nil {
		return nil, err
	}

	partners := []string{}
	for _, req := range requests {
		partners = append(partners, req.Receiver)
	}

	return partners, nil
}

func (as *authService) ApproveAccessRequest(key, id string) error {
	uid, err := as.idp.Identity(key)
	if err != nil {
		return ErrUnauthorizedAccess
	}

	return as.accesses.ApproveAccessRequest(uid, id)
}

func (as *authService) RejectAccessRequest(key, id string) error {
	uid, err := as.idp.Identity(key)
	if err != nil {
		return ErrUnauthorizedAccess
	}

	return as.accesses.RejectAccessRequest(uid, id)
}

func (as *authService) Identify(key string) (string, error) {
	id, err := as.idp.Identity(key)
	if err != nil {
		return "", err
	}

	return id, nil
}
