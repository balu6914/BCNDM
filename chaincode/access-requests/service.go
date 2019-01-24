package main

import (
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	// Pending represenets pending state of access request.
	Pending = "pending"
	// Approved represents accepted state of access request.
	Approved = "approved"
	// Revoked represents revoked state of access request.
	Revoked = "revoked"
)

var (
	// ErrConflict indicates that entity already exists and it cannot be created
	// again.
	ErrConflict = errors.New("access request already exist")

	// ErrInvalidStateTransition indicates that access request is not in
	// appropriate state in order for operation to be executed.
	ErrInvalidStateTransition = errors.New("access request not in appropriate state")

	// ErrNotFound indicates that access request doesn't exist.
	ErrNotFound = errors.New("access request doesn't exist")

	// ErrFailedKeyCreation indicates that creating composite key failed.
	ErrFailedKeyCreation = errors.New("failed to create key")

	// ErrGettingState indicates that getting state failed.
	ErrGettingState = errors.New("failed to get state")

	// ErrSettingState indicates that setting state failed.
	ErrSettingState = errors.New("failed to set state")
)

// State represents access request state.
type State string

// Service contains access request chaincode API definition.
type Service interface {
	// RequestAccess creates access request entity with pending state inside
	// blockchain. If access request already exists, it will move it to pending
	// state.
	RequestAccess(shim.ChaincodeStubInterface, string) error

	// ApproveAccess changes access request state from pending to approved.
	ApproveAccess(shim.ChaincodeStubInterface, string) error

	// RevokeAccess changes access request state from approved to revoked.
	RevokeAccess(shim.ChaincodeStubInterface, string) error

	// ListAccess returns page of access requests.
	ListAccess(shim.ChaincodeStubInterface) ([]Access, error)
}

// Access represents access request and its status.
type Access struct {
	requester string
	receiver  string
	state     State
}
