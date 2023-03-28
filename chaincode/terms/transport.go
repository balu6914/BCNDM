package main

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var (
	errIncorrectFunc    = errors.New("incorrect function name")
	errInvalidNumOfArgs = errors.New("invalid number of arguments")
)

var _ shim.Chaincode = (*chaincodeRouter)(nil)

type chaincodeRouter struct {
	svc Service
}

// NewChaincode returns new terms chaincode API instance.
func NewChaincode(svc Service) shim.Chaincode {
	return chaincodeRouter{svc: svc}
}

func (cr chaincodeRouter) Init(stub shim.ChaincodeStubInterface) pb.Response {
	function, _ := stub.GetFunctionAndParameters()

	if function != "init" {
		return shim.Error(errIncorrectFunc.Error())
	}

	return shim.Success(nil)
}

func (cr chaincodeRouter) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	// call routing
	switch function {
	case "storeTerms":
		return cr.storeTerms(stub, args)
	case "validateTerms":
		return cr.validateTerms(stub, args)
	}

	return shim.Error(errIncorrectFunc.Error())
}

func (cr chaincodeRouter) storeTerms(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(errInvalidNumOfArgs.Error())
	}

	var req termsReq
	if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
		return shim.Error(err.Error())
	}

	if err := cr.svc.StoreTerms(stub, Terms{
		StreamID:  req.StreamID,
		TermsURL:  req.TermsURL,
		TermsHash: req.TermsHash,
	}); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (cr chaincodeRouter) validateTerms(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(errInvalidNumOfArgs.Error())
	}

	var req termsReq
	if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
		return shim.Error(err.Error())
	}

	res, err := cr.svc.ValidateTerms(stub, Terms{
		StreamID:  req.StreamID,
		TermsURL:  req.TermsURL,
		TermsHash: req.TermsHash,
	})
	if err != nil {
		return shim.Error(err.Error())
	}
	vres := validationRes{Valid: res}
	payload, err := json.Marshal(vres)
	return shim.Success(payload)
}
