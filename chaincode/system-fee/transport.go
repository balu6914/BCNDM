package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var (
	errFailedSerialization = errors.New("failed to serialize response data")
	errIncorrectFunc       = errors.New("incorrect function name")
	errFailedTransfer      = errors.New("failed to transfer tokens")
)

var _ shim.Chaincode = (*chaincodeRouter)(nil)

type chaincodeRouter struct {
	svc Service
}

// NewChaincode returns new fee chaincode API instace.
func NewChaincode(svc Service) shim.Chaincode {
	return chaincodeRouter{svc: svc}
}

// Init is called during chaincode instantiation to initialize fee.
// Note that chaincode upgrade also calls this function to reset or to
// migrate data, so be careful when initializing fee.
func (cr chaincodeRouter) Init(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function != "init" {
		return shim.Error(ErrInvalidFuncCall.Error())
	}

	if len(args) != 1 {
		return shim.Error(ErrInvalidNumOfArgs.Error())
	}

	var req feeReq
	if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
		return shim.Error(ErrInvalidArgument.Error())
	}

	fee := Fee{
		Owner: req.Owner,
		Value: req.Value,
	}

	if err := cr.svc.Init(stub, fee); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode.
func (cr chaincodeRouter) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	// call routing
	switch function {
	case "fee":
		return cr.fee(stub)
	case "setFee":
		return cr.setFee(stub, args)
	case "transfer":
		return cr.transfer(stub, args)
	}
	return shim.Error(errIncorrectFunc.Error())
}

func (cr chaincodeRouter) fee(stub shim.ChaincodeStubInterface) pb.Response {
	fee := cr.svc.Fee(stub)

	res := feeRes{
		Owner: fee.Owner,
		Value: fee.Value,
	}
	payload, err := json.Marshal(res)
	if err != nil {
		return shim.Error(errFailedSerialization.Error())
	}

	return shim.Success(payload)
}

func (cr chaincodeRouter) setFee(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(ErrInvalidNumOfArgs.Error())
	}

	var req feeReq
	if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
		return shim.Error(ErrInvalidArgument.Error())
	}

	fee := Fee{
		Owner: req.Owner,
		Value: req.Value,
	}

	if err := cr.svc.SetFee(stub, fee); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (cr chaincodeRouter) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(ErrInvalidNumOfArgs.Error())
	}

	var transfer transferReq
	if err := json.Unmarshal([]byte(args[0]), &transfer); err != nil {
		return shim.Error(ErrInvalidArgument.Error())
	}

	if err := cr.svc.Transfer(stub, transfer.To, transfer.Value); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
