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

// NewChaincode returns new access request chaincode API instace.
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
	case "requestAccess":
		return cr.requestAccess(stub, args)
	case "approveAccess":
		return cr.approveAccess(stub, args)
	case "revokeAccess":
		return cr.revokeAccess(stub, args)
	case "listAccess":
		return cr.listAccess(stub, args)
	}

	return shim.Error(errIncorrectFunc.Error())
}

func (cr chaincodeRouter) requestAccess(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(errInvalidNumOfArgs.Error())
	}

	var req accessReq
	if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
		return shim.Error(err.Error())
	}

	if err := cr.svc.RequestAccess(stub, req.Receiver); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (cr chaincodeRouter) approveAccess(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(errInvalidNumOfArgs.Error())
	}

	var req approveReq
	if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
		return shim.Error(err.Error())
	}

	if err := cr.svc.ApproveAccess(stub, req.Requester); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (cr chaincodeRouter) revokeAccess(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(errInvalidNumOfArgs.Error())
	}

	var req revokeReq
	if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
		return shim.Error(err.Error())
	}

	if err := cr.svc.RevokeAccess(stub, req.Requester); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (cr chaincodeRouter) listAccess(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	reqs, err := cr.svc.ListAccess(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	res := []accessRes{}
	for _, req := range reqs {
		res = append(res, accessRes{
			Requester: req.requester,
			Receiver:  req.receiver,
			State:     req.state,
		})
	}

	payload, err := json.Marshal(&res)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(payload)
}

func (cr chaincodeRouter) grantAccess(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(errInvalidNumOfArgs.Error())
	}

	var req grantReq
	if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
		return shim.Error(err.Error())
	}

	if err := cr.svc.GrantAccess(stub, req.Destination); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
