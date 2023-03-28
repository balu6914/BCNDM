package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func main() {
	ts := NewTransferService()
	svc := NewService(ts)
	if err := shim.Start(NewChaincode(svc)); err != nil {
		fmt.Printf("Error starting Token chaincode: %s\n", err)
	}
}
