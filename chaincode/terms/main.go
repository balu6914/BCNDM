package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func main() {
	svc := NewService()
	if err := shim.Start(NewChaincode(svc)); err != nil {
		fmt.Printf("Error starting Terms chaincode: %s\n", err)
	}
}
