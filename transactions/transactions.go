package transactions

import "crypto/rand"

const (
	letters   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	secretLen = 64
)

var _ Service = (*transactionService)(nil)

type transactionService struct {
	users UserRepository
	bn    BlockchainNetwork
}

// New instantiaces domain service implementation.
func New(users UserRepository, bn BlockchainNetwork) Service {
	return transactionService{
		users: users,
		bn:    bn,
	}
}

func (ts transactionService) CreateUser(id string) error {
	secret := generate(secretLen)
	user := User{
		ID:     id,
		Secret: secret,
	}
	if err := ts.users.Save(user); err != nil {
		return err
	}

	if err := ts.bn.CreateUser(id, secret); err != nil {
		ts.users.Remove(id)
		return ErrFailedUserCreation
	}

	return nil
}

func (ts transactionService) Balance(userID string) (uint64, error) {
	return ts.bn.Balance(userID)
}

func (ts transactionService) Transfer(from, to string, value uint64) error {
	return ts.bn.Transfer(from, to, value)
}

func generate(n uint) string {
	output := make([]byte, n)
	randomness := make([]byte, n)

	rand.Read(randomness)

	l := len(letters)
	for pos := range output {
		random := uint8(randomness[pos])
		randomPos := random % uint8(l)
		output[pos] = letters[randomPos]
	}

	return string(output)
}
