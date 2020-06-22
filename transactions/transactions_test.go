package transactions_test

import (
	"fmt"
	"testing"

	"github.com/datapace/transactions"
	"github.com/datapace/transactions/mocks"

	"github.com/stretchr/testify/assert"
)

const (
	userID1         = "5281b83afbb7f35cb62d0834"
	userID2         = "5281b83afbb7f35cb62d0835"
	chanID          = "chan"
	secret          = "secret"
	balance         = 42
	remainingTokens = 100
)

func newService() transactions.Service {
	ur := mocks.NewUserRepository(map[string]string{
		userID1: secret,
		userID2: secret,
	})
	tl := mocks.NewTokenLedger(map[string]uint64{
		userID1: balance,
		userID2: balance,
	}, remainingTokens)
	cl := mocks.NewContractLedger()
	cr := mocks.NewContractRepository()
	streams := mocks.NewStreamsService(map[string]transactions.Stream{})

	return transactions.New(ur, tl, cl, cr, streams)
}

func TestCreateUser(t *testing.T) {
	svc := newService()

	cases := map[string]struct {
		id  string
		err error
	}{
		"register new user": {
			id:  "5281b83afbb7f35cb62d0836",
			err: nil,
		},
		"register existing user": {
			id:  userID1,
			err: transactions.ErrConflict,
		},
	}

	for desc, tc := range cases {
		err := svc.CreateUser(tc.id)
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
			userID:  userID1,
			balance: balance,
			err:     nil,
		},
		"fetch balance of nonexistent user": {
			userID:  "nonexistent_user",
			balance: 0,
			err:     transactions.ErrFailedBalanceFetch,
		},
	}

	for desc, tc := range cases {
		balance, err := svc.Balance(tc.userID)
		assert.Equal(t, tc.balance, balance, fmt.Sprintf("%s: expected %d got %d", desc, tc.balance, balance))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}
}

func TestTransfer(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc  string
		from  string
		to    string
		value uint64
		err   error
	}{
		{
			desc:  "transfer money from one account to another",
			from:  userID1,
			to:    userID2,
			value: balance,
			err:   nil,
		},
		{
			desc:  "transfer too much money from one account to another",
			from:  userID1,
			to:    userID2,
			value: balance,
			err:   transactions.ErrFailedTransfer,
		},
	}

	for _, tc := range cases {
		err := svc.Transfer("", tc.from, tc.to, tc.value)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
	}
}

func TestBuyTokens(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc  string
		to    string
		value uint64
		err   error
	}{
		{
			desc:  "transfer token to user account",
			to:    userID1,
			value: 1,
			err:   nil,
		},
		{
			desc:  "transfer too many tokens to user account",
			to:    userID1,
			value: remainingTokens,
			err:   transactions.ErrFailedTransfer,
		},
	}

	for _, tc := range cases {
		err := svc.BuyTokens(tc.to, tc.value)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
	}
}

func TestWithdrawTokens(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc    string
		account string
		value   uint64
		err     error
	}{
		{
			desc:    "transfer token to admin account",
			account: userID1,
			value:   1,
			err:     nil,
		},
		{
			desc:    "transfer too many tokens to admin account",
			account: userID1,
			value:   remainingTokens,
			err:     transactions.ErrFailedTransfer,
		},
	}

	for _, tc := range cases {
		err := svc.WithdrawTokens(tc.account, tc.value)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
	}
}
