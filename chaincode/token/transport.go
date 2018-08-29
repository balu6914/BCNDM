package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var (
	errFailedTransfer = errors.New("failed to transfer tokens")
	errFailedApprove  = errors.New("failed to approve allowance")
	errIncorrectFunc  = errors.New("incorrect function name")
)

var _ shim.Chaincode = (*chaincodeRouter)(nil)

type chaincodeRouter struct {
	svc Service
}

// NewChaincode returns new token chaincode API instace.
func NewChaincode(svc Service) shim.Chaincode {
	return chaincodeRouter{svc: svc}
}

// Init is called during chaincode instantiation to initialize token.
// Note that chaincode upgrade also calls this function to reset or to
// migrate data, so be careful when initializing token.
func (cr chaincodeRouter) Init(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function != "init" {
		return shim.Error(ErrInvalidFuncCall.Error())
	}

	if len(args) != 1 {
		return shim.Error(ErrInvalidNumOfArgs.Error())
	}

	// get token data from JSON
	var ti tokenInfo
	if err := json.Unmarshal([]byte(args[0]), &ti); err != nil {
		return shim.Error(ErrInvalidArgument.Error())
	}

	info := TokenInfo{
		Name:        ti.Name,
		Symbol:      ti.Symbol,
		Decimals:    ti.Decimals,
		TotalSupply: ti.TotalSupply,
	}

	if err := cr.svc.Init(stub, info); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode.
func (cr chaincodeRouter) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	// call routing
	switch function {
	// Returns a Token totalSupply
	case "totalSupply":
		return cr.totalSupply(stub)
	// Transfers _value amount of tokens to address _to
	case "transfer":
		return cr.transfer(stub, args)
	// Returns the account balance of another account.
	case "balanceOf":
		return cr.balance(stub, args)
	// Allows _spender to withdraw from your account multiple times, up to the _value amount.
	// If this function is called again it overwrites the current allowance with _value.
	case "approve":
		return cr.approve(stub, args)
	// Returns the amount which _spender is still allowed to withdraw from _owner.
	case "allowance":
		return cr.allowance(stub, args)
	// Transfers _value amount of tokens from address _from to address _to
	case "transferFrom":
		return cr.transferFrom(stub, args)
	case "groupTransfer":
		return cr.groupTransfer(stub, args)
	}
	return shim.Error(errIncorrectFunc.Error())
}

func (cr chaincodeRouter) totalSupply(stub shim.ChaincodeStubInterface) pb.Response {
	total, err := cr.svc.TotalSupply(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	balance := balanceRes{Value: total}
	payload, err := json.Marshal(balance)
	if err != nil {
		return shim.Error(ErrFailedSerialization.Error())
	}

	return shim.Success(payload)
}

func (cr chaincodeRouter) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(ErrInvalidNumOfArgs.Error())
	}

	var transfer transferReq
	if err := json.Unmarshal([]byte(args[0]), &transfer); err != nil {
		return shim.Error(ErrInvalidArgument.Error())
	}

	if ok := cr.svc.Transfer(stub, transfer.To, transfer.Value); !ok {
		return shim.Error(errFailedTransfer.Error())
	}

	return shim.Success(nil)
}

func (cr chaincodeRouter) balance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(ErrInvalidNumOfArgs.Error())
	}

	var req balanceReq
	if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
		return shim.Error(ErrInvalidArgument.Error())
	}

	balance, err := cr.svc.BalanceOf(stub, req.Owner)
	if err != nil {
		return shim.Error(err.Error())
	}

	res := balanceRes{Value: balance}

	payload, err := json.Marshal(res)
	if err != nil {
		return shim.Error(ErrFailedSerialization.Error())
	}

	return shim.Success(payload)
}

func (cr chaincodeRouter) approve(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(ErrInvalidNumOfArgs.Error())
	}

	var req approveReq
	if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
		return shim.Error(ErrInvalidArgument.Error())
	}

	if ok := cr.svc.Approve(stub, req.Spender, req.Value); !ok {
		return shim.Error(errFailedApprove.Error())
	}

	return shim.Success(nil)
}

func (cr chaincodeRouter) allowance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(ErrInvalidNumOfArgs.Error())
	}

	var req allowanceReq
	if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
		return shim.Error(ErrInvalidArgument.Error())
	}

	allowance, err := cr.svc.Allowance(stub, req.Owner, req.Spender)
	if err != nil {
		return shim.Error(err.Error())
	}

	res := allowanceRes{Value: allowance}
	payload, err := json.Marshal(res)
	if err != nil {
		return shim.Error(ErrFailedSerialization.Error())
	}

	return shim.Success(payload)
}

func (cr chaincodeRouter) transferFrom(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(ErrInvalidNumOfArgs.Error())
	}

	var req transferFromReq
	if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
		return shim.Error(ErrInvalidArgument.Error())
	}

	if ok := cr.svc.TransferFrom(stub, req.From, req.To, req.Value); !ok {
		return shim.Error(errFailedTransfer.Error())
	}

	return shim.Success(nil)
}

func (cr chaincodeRouter) groupTransfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(ErrInvalidNumOfArgs.Error())
	}

	var req []transferReq
	if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
		return shim.Error(ErrInvalidArgument.Error())
	}

	transfers := []Transfer{}
	for _, tr := range req {
		transfers = append(transfers, Transfer{
			To:    tr.To,
			Value: tr.Value,
		})
	}

	if err := cr.svc.GroupTransfer(stub, transfers...); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
