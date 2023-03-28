package main

import (
	"errors"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var (
	// ErrInvalidNumOfArgs indicates invalid number of arguments.
	ErrInvalidNumOfArgs = errors.New("invalid number of arguments")

	// ErrInvalidFuncCall indicates invalid function name invocation.
	ErrInvalidFuncCall = errors.New("called invalid function")

	// ErrInvalidArgument indicates invalid argument format.
	ErrInvalidArgument = errors.New("invalid argument passed")

	// ErrFailedKeyCreation indicates that creating composite key failed.
	ErrFailedKeyCreation = errors.New("failed to create key")

	// ErrSettingState indicates that setting state failed.
	ErrSettingState = errors.New("failed to set state")

	// ErrGettingState indicates that getting state failed.
	ErrGettingState = errors.New("failed to get state")

	// ErrNotFound indicates that wanted entity doesn't exist.
	ErrNotFound = errors.New("entity not found")

	// ErrInvalidStateData indicates that chaincode read invalid state data.
	ErrInvalidStateData = errors.New("received invalid state data")
)

// Service defines API for creating contract and applying it on transfer.
type Service interface {
	// CreateContract creates new user defined contract.
	CreateContracts(shim.ChaincodeStubInterface, ...Contract) error

	// SignContract signs existing user defined contract.
	SignContract(shim.ChaincodeStubInterface, Contract) error

	// Transfers tokens and applies user defined contract.
	Transfer(shim.ChaincodeStubInterface, string, string, string, time.Time, uint64) error
}

// Contract contains user defined contract data.
type Contract struct {
	StreamID  string    `json:"stream_id,omitempty"`
	StartTime time.Time `json:"start_time,omitempty"`
	EndTime   time.Time `json:"end_time,omitempty"`
	OwnerID   string    `json:"owner_id,omitempty"`
	PartnerID string    `json:"partner_id,omitempty"`
	Share     uint64    `json:"share,omitempty"`
	Signed    bool      `json:"signed,omitempty"`
}
