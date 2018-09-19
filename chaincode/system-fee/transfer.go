package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	erc20Chaincode = "token"
	transfer       = "groupTransfer"
)

// TransferService defines token transfer service API.
type TransferService interface {
	// Transfer given amount from callers account to specified account. Returns
	// error only if transaction can't be executed.
	Transfer(shim.ChaincodeStubInterface, ...Transfer) error
}

var _ TransferService = (*transferService)(nil)

type transferService struct {
	accounts map[string]uint64
}

// NewTransferService returns token transfer service instance.
func NewTransferService() TransferService {
	return transferService{
		accounts: map[string]uint64{},
	}
}

func (ts transferService) Transfer(stub shim.ChaincodeStubInterface, transfers ...Transfer) error {
	req := []transferReq{}
	for _, tr := range transfers {
		req = append(req, transferReq{
			To:    tr.To,
			Value: tr.Value,
		})
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}

	args := [][]byte{[]byte(transfer), payload}
	res := stub.InvokeChaincode(erc20Chaincode, args, stub.GetChannelID())
	if res.GetStatus() != shim.OK {
		return errors.New(res.GetMessage())
	}

	return nil
}
