package main

import (
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var (
	// ErrInvalidNumOfArgs indicates invalid number of arguments.
	ErrInvalidNumOfArgs = errors.New("invalid number of arguments")

	// ErrInvalidFuncCall indicates invalid function name invocation.
	ErrInvalidFuncCall = errors.New("called invalid function")

	// ErrInvalidArgument indicates invalid argument format.
	ErrInvalidArgument = errors.New("invalid argument passed")

	// ErrSettingState indicates that setting state failed.
	ErrSettingState = errors.New("failed to set state")

	// ErrFailedTransfer indicates that token transfer failed.
	ErrFailedTransfer = errors.New("failed to transfer tokens")

	// ErrUnauthorized indicates that reading data from certificate failed.
	ErrUnauthorized = errors.New("failed to read data from certificate")
)

// Service defines API for taking, fetching and setting system fee.
type Service interface {
	// Init sets initial system fee.
	Init(shim.ChaincodeStubInterface) error

	// Fee returns current system fee.
	Fee(shim.ChaincodeStubInterface) uint64

	// SetFee sets system fee value.
	SetFee(shim.ChaincodeStubInterface, uint64) error

	// Transfer given amount from callers account to specified account and
	// fee to platform owner account. Returns error only if transaction can't be
	// executed.
	Transfer(shim.ChaincodeStubInterface, string, uint64) error
}
