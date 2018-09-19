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
	Init(shim.ChaincodeStubInterface, Fee) error

	// Fee returns current system fee.
	Fee(shim.ChaincodeStubInterface) Fee

	// SetFee sets system fee value.
	SetFee(shim.ChaincodeStubInterface, Fee) error

	// Transfer given amount from callers account to specified account and
	// fee to platform owner account. Returns error only if transaction can't be
	// executed.
	Transfer(shim.ChaincodeStubInterface, string, uint64, ...Transfer) error
}

// Transfer contains transfer data.
type Transfer struct {
	To    string
	Value uint64
}

// Fee contains system owner CN and fee value that will go to owner.
type Fee struct {
	Owner string `json:"owner"`
	Value uint64 `json:"value"`
}
