package main

import (
	"fmt"

	"github.com/datapace/datapace/examples/blockchain"
	"github.com/datapace/datapace/examples/fabric-ca/users"
)

func main() {

	fSetup := blockchain.FabricSetup{
		OrgAdmin:   "Admin@org1.datapace.com",
		OrgName:    "Org1",
		ConfigFile: "../../config/fabric/config.yaml",
		ChannelID:  "datapacechannel",
	}

	// Initialization of the Fabric SDK from the previously set properties
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
	}
	fmt.Println("Successfully connected to Fabric network")

	bc := users.BcNetwork{Fabric: &fSetup}

	// Create New user in Fabric network calling fabric-ca
	newUser, err := bc.CreateUser("Nikola", "12345")
	if err != nil {
		fmt.Printf("Unable to create a user in the fabric-ca %v\n", err)
	}

	fmt.Printf("User created!: %v\n", newUser)
}
