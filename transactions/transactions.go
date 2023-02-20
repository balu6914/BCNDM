package transactions

import (
	"crypto/rand"
)

const (
	maxShare  = 90000
	letters   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	secretLen = 64
)

type TokenInfo struct {
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	Decimals      uint8  `json:"decimals"`
	ContractOwner string `json:"contractOwner"`
}

type TransferFrom struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Value    uint64 `json:"value"`
	DateTime string `json:"dateTime"` // dateTime should be added at middleware level in format: DD-MM-YYYY hh:mm:ss
	TxType   string `json:"txType"`
}

type TokenTxHistory struct {
	TokenInfo TokenInfo      `json:"tokenInfo"`
	TxList    []TransferFrom `json:"txList"`
}

var _ Service = (*transactionService)(nil)

type transactionService struct {
	users     UserRepository
	tokens    TokenLedger
	conLedger ContractLedger
	conRepo   ContractRepository
	streams   StreamsService
}

// New instantiaces domain service implementation.
func New(users UserRepository, tokens TokenLedger, conLedger ContractLedger, conRepo ContractRepository, streams StreamsService) Service {
	return transactionService{
		users:     users,
		tokens:    tokens,
		conLedger: conLedger,
		conRepo:   conRepo,
		streams:   streams,
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

	if err := ts.tokens.CreateUser(id, secret); err != nil {
		ts.users.Remove(id)
		return ErrFailedUserCreation
	}

	return nil
}

func (ts transactionService) Balance(userID string) (uint64, error) {
	balance, err := ts.tokens.Balance(userID)
	if err != nil {
		return 0, ErrFailedBalanceFetch
	}

	return balance, nil
}

func (ts transactionService) Transfer(streamID, from, to string, value uint64) error {
	if err := ts.tokens.Transfer(streamID, from, to, value); err != nil {
		if err == ErrNotEnoughTokens {
			return ErrNotEnoughTokens
		}
		return ErrFailedTransfer
	}

	return nil
}

func (ts transactionService) BuyTokens(account string, value uint64) error {
	if err := ts.tokens.BuyTokens(account, value); err != nil {
		return ErrFailedTransfer
	}

	return nil
}

func (ts transactionService) WithdrawTokens(account string, value uint64) error {
	if err := ts.tokens.WithdrawTokens(account, value); err != nil {
		return ErrFailedTransfer
	}

	return nil
}

func (ts transactionService) TxHistory(userID, fromDateTime, toDateTime, txType string) (TokenTxHistory, error) {
	txHistory, err := ts.tokens.TxHistory(userID, fromDateTime, toDateTime, txType)
	if err != nil {
		return txHistory, ErrFailedTxHistoryFetch
	}

	return txHistory, nil
}

func (ts transactionService) CreateContracts(contracts ...Contract) error {
	sum := uint64(0)
	for _, contract := range contracts {
		sum += contract.Share
	}
	if sum > maxShare {
		return ErrConflict
	}

	if len(contracts) == 0 {
		return ErrMalformedEntity
	}

	stream, err := ts.streams.One(contracts[0].StreamID)
	if err != nil {
		return ErrNotFound
	}

	for i := range contracts {
		contracts[i].StreamName = stream.Name
	}

	if err := ts.conRepo.Create(contracts...); err != nil {
		return err
	}

	if err := ts.conLedger.Create(contracts...); err != nil {
		return err
	}

	ts.conRepo.Activate(contracts...)

	return nil
}

func (ts transactionService) SignContract(contract Contract) error {
	if err := ts.conLedger.Sign(contract); err != nil {
		return err
	}

	if err := ts.conRepo.Sign(contract); err != nil {
		return err
	}

	return nil
}

func (ts transactionService) ListContracts(owner string, pageNo uint64, limit uint64, role Role) ContractPage {
	return ts.conRepo.List(owner, pageNo, limit, role)
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
