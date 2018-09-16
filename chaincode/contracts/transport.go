package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var errIncorrectFunc = errors.New("incorrect function name")

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
	function, _ := stub.GetFunctionAndParameters()

	if function != "init" {
		return shim.Error(ErrInvalidFuncCall.Error())
	}

	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode.
func (cr chaincodeRouter) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	// call routing
	switch function {
	case "createContracts":
		return cr.createContracts(stub, args)
	case "signContract":
		return cr.signContract(stub, args)
	case "transfer":
		return cr.transfer(stub, args)
	}

	return shim.Error(errIncorrectFunc.Error())
}

func (cr chaincodeRouter) createContracts(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(ErrInvalidNumOfArgs.Error())
	}

	var req []contractReq
	if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
		return shim.Error(ErrInvalidArgument.Error())
	}

	contracts := []Contract{}
	for _, c := range req {
		contract := Contract{
			StreamID:  c.StreamID,
			StartTime: c.StartTime,
			EndTime:   c.EndTime,
			OwnerID:   c.OwnerID,
			PartnerID: c.PartnerID,
			Share:     c.Share,
		}
		contracts = append(contracts, contract)
	}

	if err := cr.svc.CreateContracts(stub, contracts...); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (cr chaincodeRouter) signContract(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(ErrInvalidNumOfArgs.Error())
	}

	var req contractReq
	if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
		return shim.Error(ErrInvalidArgument.Error())
	}

	contract := Contract{
		StreamID: req.StreamID,
		EndTime:  req.EndTime,
	}
	if err := cr.svc.SignContract(stub, contract); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (cr chaincodeRouter) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(ErrInvalidNumOfArgs.Error())
	}

	var req transferReq
	if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
		return shim.Error(ErrInvalidArgument.Error())
	}

	if err := cr.svc.Transfer(stub, req.StreamID, req.To, req.Time, req.Value); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
