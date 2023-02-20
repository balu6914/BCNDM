package mocks

import (
	"sync"
	"time"

	"github.com/datapace/datapace/transactions"
)

var _ transactions.TokenLedger = (*mockNetwork)(nil)

const (
	dateFormat = "2006-01-02 15:04:05"
)

type mockNetwork struct {
	users     map[string]uint64
	txHistory map[string][]transactions.TransferFrom
	remaining uint64
	mutex     *sync.Mutex
}

// NewTokenLedger returns mock instance of blockchain network.
func NewTokenLedger(users map[string]uint64, remaining uint64) transactions.TokenLedger {
	txHistory := make(map[string][]transactions.TransferFrom)
	for user, _ := range users {
		txHistory[user] = []transactions.TransferFrom{}
	}
	return &mockNetwork{
		users:     users,
		txHistory: txHistory,
		remaining: remaining,
		mutex:     &sync.Mutex{},
	}
}

func (mn *mockNetwork) CreateUser(id, secret string) error {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	if _, ok := mn.users[id]; ok {
		return transactions.ErrFailedUserCreation
	}

	mn.users[id] = 0
	mn.txHistory[id] = []transactions.TransferFrom{}
	return nil
}

func (mn *mockNetwork) Balance(name string) (uint64, error) {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	balance, ok := mn.users[name]
	if !ok {
		return 0, transactions.ErrNotFound
	}

	return balance, nil
}

func (mn *mockNetwork) Transfer(stream, from, to string, value uint64) error {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	balance, ok := mn.users[from]
	if !ok {
		return transactions.ErrFailedTransfer
	}

	if balance < value {
		return transactions.ErrFailedTransfer
	}

	mn.users[to] = mn.users[to] + value
	mn.users[from] = mn.users[from] - value
	transfer := transactions.TransferFrom{
		From:     from,
		To:       to,
		Value:    value,
		DateTime: time.Now().UTC().Format(dateFormat),
	}

	mn.txHistory[to] = append(mn.txHistory[to], transfer)
	mn.txHistory[from] = append(mn.txHistory[from], transfer)
	return nil
}

func (mn *mockNetwork) BuyTokens(account string, value uint64) error {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	if value > mn.remaining {
		return transactions.ErrFailedTransfer
	}

	mn.users[account] += value
	mn.remaining -= value
	transfer := transactions.TransferFrom{
		From:     "treasury",
		To:       account,
		Value:    value,
		DateTime: time.Now().UTC().Format(dateFormat),
	}

	mn.txHistory[account] = append(mn.txHistory[account], transfer)
	return nil
}

func (mn *mockNetwork) WithdrawTokens(account string, value uint64) error {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	if value > mn.users[account] {
		return transactions.ErrFailedTransfer
	}

	mn.users[account] -= value
	mn.remaining += value
	transfer := transactions.TransferFrom{
		From:     account,
		To:       "treasury",
		Value:    value,
		DateTime: time.Now().UTC().Format(dateFormat),
	}

	mn.txHistory[account] = append(mn.txHistory[account], transfer)
	return nil
}

func (mn *mockNetwork) TxHistory(name, fromDateTime, toDateTime, txType string) (transactions.TokenTxHistory, error) {
	txHis := new(transactions.TokenTxHistory)
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	txHistory, ok := mn.txHistory[name]
	if !ok {
		return *txHis, transactions.ErrNotFound
	}

	txHistoryRes := transactions.TokenTxHistory{
		TokenInfo: transactions.TokenInfo{
			Name:          "token name",
			Symbol:        "token",
			Decimals:      8,
			ContractOwner: "treasury",
		},
		TxList: txHistory,
	}

	return txHistoryRes, nil
}
