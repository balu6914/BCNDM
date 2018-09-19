package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	feeChaincode = "fee"
	transfer     = "transfer"
)

// TransferService defines token transfer service API.
type TransferService interface {
	// Transfer given amount from callers account to specified account. Returns
	// error only if transaction can't be executed.
	Transfer(shim.ChaincodeStubInterface, string, uint64, ...Transfer) error
}

// Transfer contains transfer data.
type Transfer struct {
	To    string
	Value uint64
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

func (ts transferService) Transfer(stub shim.ChaincodeStubInterface, owner string, value uint64, transfers ...Transfer) error {
	trs := []transferReq{}
	for _, tr := range transfers {
		trs = append(trs, transferReq{
			To:    tr.To,
			Value: tr.Value,
		})
	}
	req := transferStatusReq{
		Owner:     owner,
		Value:     value,
		Transfers: trs,
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}

	args := [][]byte{[]byte(transfer), payload}
	res := stub.InvokeChaincode(feeChaincode, args, stub.GetChannelID())
	if res.GetStatus() != shim.OK {
		return errors.New(res.GetMessage())
	}

	return nil
}

type transferStatusReq struct {
	Owner     string        `json:"owner"`
	Value     uint64        `json:"value"`
	Transfers []transferReq `json:"transfers"`
}
