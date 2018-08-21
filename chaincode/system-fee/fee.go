package main

import (
	"encoding/binary"
	"encoding/json"
	"math"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	feeKey  = "__fee"
	ownerCN = "Admin@org1.monetasa.com"

	// Decimals represent number of zeros in 100% value, in this case it's 5, so
	// 100% == 100000.
	decimals = 5
)

var _ Service = (*feeChaincode)(nil)

type feeChaincode struct {
	ts TransferService
}

// NewService returns fee contract implementation.
func NewService(ts TransferService) Service {
	return feeChaincode{ts: ts}
}

func (fc feeChaincode) Init(stub shim.ChaincodeStubInterface) error {
	function, args := stub.GetFunctionAndParameters()

	if function != "init" {
		return ErrInvalidFuncCall
	}

	if len(args) != 1 {
		return ErrInvalidNumOfArgs
	}

	var req feeReq
	if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
		return ErrInvalidArgument
	}

	if err := fc.SetFee(stub, req.Value); err != nil {
		return err
	}

	return nil
}

func (fc feeChaincode) Fee(stub shim.ChaincodeStubInterface) uint64 {
	data, err := stub.GetState(feeKey)
	if err != nil || data == nil {
		return 0
	}

	return binary.LittleEndian.Uint64(data)
}

func (fc feeChaincode) SetFee(stub shim.ChaincodeStubInterface, value uint64) error {
	if value > uint64(math.Pow(10, decimals)) {
		return ErrInvalidArgument
	}

	cn, err := callerCN(stub)
	if err != nil || cn != ownerCN {
		return ErrUnauthorized
	}

	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, value)
	if err := stub.PutState(feeKey, data); err != nil {
		return ErrSettingState
	}

	return nil
}

func (fc feeChaincode) Transfer(stub shim.ChaincodeStubInterface, to string, value uint64) error {
	fee := fc.Fee(stub)
	wholeValue := uint64(math.Pow(10, decimals))

	feeValue := fee * value / wholeValue

	transfers := []Transfer{
		Transfer{
			To:    to,
			Value: value - feeValue,
		},
		Transfer{
			To:    ownerCN,
			Value: feeValue,
		},
	}

	// Transfer tokens.
	return fc.ts.Transfer(stub, transfers...)
}
