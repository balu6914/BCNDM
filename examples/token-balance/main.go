package main

import (
	"fmt"
	"monetasa/examples/blockchain"
	"monetasa/examples/token-balance/token"
	"os"
)

func main() {

	fSetup := blockchain.FabricSetup{
		OrgAdmin:    "admin",
		OrgName:     "Org1",
		ConfigFile:  "../../config/fabric/config.yaml",
		ChannelID:   "myc",
		ChaincodeID: "token",
	}

	// Initialization of the Fabric SDK from the previously set properties
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		os.Exit(3)
	}
	fmt.Println("Successfully connected to Fabric network")

	b := token.BcNetwork{Fabric: &fSetup}
	// Get balance for user Nikola
	balance, err := b.Balance("Nikola")

	if err != nil {
		fmt.Println("Error fetching balance!!!", err)
		os.Exit(3)
	}

	fmt.Println("Great here is a balance", balance)

}
