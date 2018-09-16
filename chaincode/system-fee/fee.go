package main

import (
	"encoding/json"
	"math"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	feeKey  = "__fee"
	adminCN = "Admin@org1.monetasa.com"

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

func (fc feeChaincode) Init(stub shim.ChaincodeStubInterface, fee Fee) error {
	return fc.SetFee(stub, fee)
}

func (fc feeChaincode) Fee(stub shim.ChaincodeStubInterface) Fee {
	data, err := stub.GetState(feeKey)
	if err != nil || data == nil {
		return Fee{
			Owner: adminCN,
			Value: 0,
		}
	}

	var fee Fee
	if err := json.Unmarshal(data, &fee); err != nil {
		return Fee{
			Owner: adminCN,
			Value: 0,
		}
	}

	return fee
}

func (fc feeChaincode) SetFee(stub shim.ChaincodeStubInterface, fee Fee) error {
	if fee.Value > uint64(math.Pow(10, decimals)) {
		return ErrInvalidArgument
	}

	cn, err := callerCN(stub)
	if err != nil || cn != adminCN {
		return ErrUnauthorized
	}

	data, err := json.Marshal(fee)
	if err != nil {
		return ErrInvalidArgument
	}

	if err := stub.PutState(feeKey, data); err != nil {
		return ErrSettingState
	}

	return nil
}

func (fc feeChaincode) Transfer(stub shim.ChaincodeStubInterface, owner string, value uint64, transfers ...Transfer) error {
	sum := value
	for _, t := range transfers {
		sum += t.Value
	}

	fee := fc.Fee(stub)
	wholeValue := uint64(math.Pow(10, decimals))

	feeValue := fee.Value * sum / wholeValue

	ts := []Transfer{
		{
			To:    fee.Owner,
			Value: feeValue,
		},
		{
			To:    owner,
			Value: value - feeValue,
		},
	}
	for _, t := range ts {
		ts = append(ts, Transfer{
			To:    t.To,
			Value: t.Value,
		})
	}

	// Transfer tokens.
	return fc.ts.Transfer(stub, ts...)
}
