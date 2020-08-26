package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const objectType = "terms"

var _ Service = (*termsChaincode)(nil)

type termsChaincode struct{}

// NewService returns terms storage request chaincode implementation.
func NewService() Service {
	return termsChaincode{}
}

func (tc termsChaincode) StoreTerms(stub shim.ChaincodeStubInterface, terms Terms) error {
	if err := stub.PutState(terms.StreamID, []byte(terms.TermsHash)); err != nil {
		return ErrStoringTerms
	}
	return nil
}
func (tc termsChaincode) ValidateTerms(stub shim.ChaincodeStubInterface, terms Terms) (bool, error) {
	data, err := stub.GetState(terms.StreamID)
	if err != nil {
		return false, ErrGettingTerms
	}
	if string(data) != terms.TermsHash {
		return false, nil
	}
	return true, nil
}
