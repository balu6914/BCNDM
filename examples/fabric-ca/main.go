package main

import (
	"fmt"

	"monetasa/examples/fabric-ca/blockchain"
)

func main() {

	// Definition of the Fabric SDK properties
	fSetup := blockchain.FabricSetup{
		OrgAdmin:   "Admin",
		OrgName:    "Org1",
		ConfigFile: "config.yaml",
		// Channel parameters
		ChannelID: "myc",
	}

	// Initialization of the Fabric SDK from the previously set properties
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
	}
	fmt.Printf("Successfully connected to Fabric network")
}
