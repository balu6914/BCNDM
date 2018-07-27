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
		if err == ErrConflict || err == ErrMalformedEntity {
			return err
		}
		return ErrFailedUserCreation
	}

	if err := ts.bn.CreateUser(id, secret); err != nil {
		ts.users.Remove(id)
		return ErrFailedUserCreation
	}

	return nil
}

func (ts transactionService) Balance(userID string) (uint64, error) {
	balance, err := ts.bn.Balance(userID)
	if err != nil {
		return 0, ErrFailedBalanceFetch
	}

	return balance, nil
}

func (ts transactionService) Transfer(from, to string, value uint64) error {
	if err := ts.bn.Transfer(from, to, value); err != nil {
		return ErrFailedTransfer
	}

	return nil
}

func (ts transactionService) BuyTokens(account string, value uint64) error {
	if err := ts.bn.BuyTokens(account, value); err != nil {
		return ErrFailedTransfer
	}

	return nil
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
