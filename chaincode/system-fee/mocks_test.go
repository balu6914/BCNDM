package main_test

import (
	fee "github.com/datapace/chaincode/system-fee"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var _ fee.TransferService = (*mockTransferService)(nil)

type mockTransferService struct{}

func NewMockTransferService() fee.TransferService {
	return mockTransferService{}
}

func (mts mockTransferService) Transfer(stub shim.ChaincodeStubInterface, transfers ...fee.Transfer) error {
	return nil
}
