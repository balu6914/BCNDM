package main

import (
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var (

	// ErrFailedKeyCreation indicates that creating composite key failed.
	ErrFailedKeyCreation = errors.New("failed to create key")

	// ErrStoringTerms indicates that setting terms failed.
	ErrStoringTerms = errors.New("failed to store terms")

	// ErrGettingTerms indicates that getting state failed.
	ErrGettingTerms = errors.New("failed to get terms")
)

// Service contains terms chaincode API definition.
type Service interface {
	// StoreTerms creates new terms
	StoreTerms(shim.ChaincodeStubInterface, Terms) error

	// ValidateTerms validates terms hash against blockchain
	ValidateTerms(shim.ChaincodeStubInterface, Terms) (bool, error)
}

// Terms contains terms data.
type Terms struct {
	StreamID  string `json:"stream_id"`
	TermsURL  string `json:"terms_url,omitempty"`
	TermsHash string `json:"terms_hash"`
}
