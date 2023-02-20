package main

import (
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var (
	// ErrInvalidFuncCall indicates invalid function name invocation.
	ErrInvalidFuncCall = errors.New("called invalid function")

	// ErrInvalidNumOfArgs indicates invalid number of arguments.
	ErrInvalidNumOfArgs = errors.New("invalid number of arguments")

	// ErrInvalidArgument indicates invalid argument format.
	ErrInvalidArgument = errors.New("invalid argument passed")

	// ErrSettingState indicates that setting state failed.
	ErrSettingState = errors.New("failed to set state")

	// ErrDeletingState indicates that deletion of state failed.
	ErrDeletingState = errors.New("failed to delete state")

	// ErrGettingState indicates that getting state failed.
	ErrGettingState = errors.New("failed to get state")

	// ErrInvalidStateData indicates that chaincode read invalid state data.
	ErrInvalidStateData = errors.New("received invalid state data")

	// ErrUnauthorized indicates that reading data from certificate failed.
	ErrUnauthorized = errors.New("failed to read data from certificate")

	// ErrFailedBalanceSet indicates that setting balance failed.
	ErrFailedBalanceSet = errors.New("failed to set balance")

	// ErrFailedKeyCreation indicates that creating composite key failed.
	ErrFailedKeyCreation = errors.New("failed to create key")

	// ErrNotEnoughTokens indicates that spender doesn't have enough tokens.
	ErrNotEnoughTokens = errors.New("not enough tokens")

	// ErrOverflow indicates that overflow error happened while calculating new
	// balance.
	ErrOverflow = errors.New("overflow error")

	// ErrFailedSerialization indicates that object serialization failed.
	ErrFailedSerialization = errors.New("failed to serialize response data")

	// ErrFailedRichQuery indicates that it fails to execute mongo query on world state
	ErrFailedRichQuery = errors.New("failed to execute rich query")
)

// Service defines ERC20 compliant interface.
type Service interface {
	// Init set initial token supply.
	Init(shim.ChaincodeStubInterface, TokenInfo) error

	// TotalSupply returns total token supply.
	TotalSupply(shim.ChaincodeStubInterface) (uint64, error)

	// BalanceOf returns wanted account balance & list of deltas to calculate balance.
	BalanceOf(shim.ChaincodeStubInterface, string) (uint64, []string, error)

	// Transfer given amount from callers account to specified account. Returns
	// true only if transaction can be executed.
	Transfer(shim.ChaincodeStubInterface, string, string, uint64) bool

	// Transfer given amount from specified to specified account. Returns true
	// only if transaction can be executed.
	TransferFrom(shim.ChaincodeStubInterface, string, string, string, uint64) bool

	// Approve allows specified spender to withdraw from callers account
	// multiple times up to specified value.
	Approve(shim.ChaincodeStubInterface, string, uint64) bool

	// Allowance returns the amount which specified spender is still allowed to
	// withdraw from specified owner.
	Allowance(shim.ChaincodeStubInterface, string, string) (uint64, error)

	// GroupTransfer given amount of tokens from callers account to
	GroupTransfer(shim.ChaincodeStubInterface, ...Transfer) error

	// TxHistory returns list of transactions
	TxHistory(shim.ChaincodeStubInterface, string, string, string, string) ([]TransferFrom, *TokenInfo, error)

	// CollectDeltasForTreasury collects all deltas for treasury account and combines them
	// it is recommended to execute this method at a regular interval depending upon tx load in the system
	CollectDeltasForTreasury(shim.ChaincodeStubInterface) error
}

type TokenInfo struct {
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	Decimals      uint8  `json:"decimals"`
	ContractOwner string `json:"contractOwner"`
}

type TokenDelta struct {
	Value     uint64 `json:"value"`
	Operation string `json:"operation"`
}

type Balance struct {
	User  string `json:"user"`
	Value uint64 `json:"value"`
}

type TransferFrom struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Value     uint64 `json:"value"`
	DateTime  string `json:"dateTime"` // dateTime should be added at middleware level in format: DD-MM-YYYY hh:mm:ss
	TxType    string `json:"txType"`
	EpochTime int64  `json:"epochTime"`
}

type Approve struct {
	Spender string `json:"spender"`
	Value   uint64 `json:"value"`
}

// Transfer contrains data necessary to transfer tokens.
type Transfer struct {
	To       string
	Value    uint64
	DateTime string `json:"dateTime"` // dateTime should be added at middleware level in format: DD-MM-YYYY hh:mm:ss
	TxType   string `json:"txType"`
}
