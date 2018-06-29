package transactions_test

import (
	"fmt"
	"monetasa/transactions"
	"monetasa/transactions/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	userID  = "user"
	chanID  = "chan"
	secret  = "secret"
	balance = 42
)

func newService() transactions.Service {
	bn := mocks.NewBlockchainNetwork(map[string]uint64{
		userID: balance,
	})

	return transactions.New(bn)
}

func TestCreateUser(t *testing.T) {
	svc := newService()

	cases := map[string]struct {
		id     string
		secret string
		key    []byte
		err    error
	}{
		"register new user": {
			id:     "nonexistent_id",
			secret: secret,
			key:    []byte(secret),
			err:    nil,
		},
		"register existing user": {
			id:     userID,
			secret: secret,
			key:    []byte{},
			err:    transactions.ErrFailedUserCreation,
		},
	}

	for desc, tc := range cases {
		key, err := svc.CreateUser(tc.id, tc.secret)
		assert.Equal(t, tc.key, key, fmt.Sprintf("%s: expected %s got %s", desc, tc.key, key))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}
}

func TestBalance(t *testing.T) {
	svc := newService()

	cases := map[string]struct {
		userID  string
		chanID  string
		balance uint64
		err     error
	}{
		"fetch balance of existing user": {
			userID:  userID,
			chanID:  chanID,
			balance: balance,
			err:     nil,
		},
		"fetch balance of nonexistent user": {
			userID:  "nonexistent_user",
			chanID:  chanID,
			balance: 0,
			err:     transactions.ErrFailedBalanceFetch,
		},
	}

	for desc, tc := range cases {
		balance, err := svc.Balance(tc.userID, tc.chanID)
		assert.Equal(t, tc.balance, balance, fmt.Sprintf("%s: expected %d got %d", desc, tc.balance, balance))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}
}
