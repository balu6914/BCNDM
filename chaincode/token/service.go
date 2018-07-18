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
)

// Service defined ERC20 compliant interface.
type Service interface {
	// Init set initial token supply.
	Init(shim.ChaincodeStubInterface) error

	// TotalSupply returns total token supply.
	TotalSupply(shim.ChaincodeStubInterface) (uint64, error)

	// BalanceOf returns wanted account balance.
	BalanceOf(shim.ChaincodeStubInterface, string) (uint64, error)

	// Transfer given amount from callers account to specified account. Returns
	// true only if transaction can be executed.
	Transfer(shim.ChaincodeStubInterface, string, uint64) bool

	// Transfer given amount from specified to specified account. Returns true
	// only if transaction can be executed.
	TransferFrom(shim.ChaincodeStubInterface, string, string, uint64) bool

	// Approve allows specified spender to withdraw from callers account
	// multiple times up to specified value.
	Approve(shim.ChaincodeStubInterface, string, uint64) bool

	// Allowance returns the amount which specified spender is still allowed to
	// withdraw from specified owner.
	Allowance(shim.ChaincodeStubInterface, string, string) (uint64, error)
}
