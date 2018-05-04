package main

import (
	"fmt"
	"monetasa/examples/blockchain"
	"monetasa/examples/fabric-ca/users"
)

func main() {

	fSetup := blockchain.FabricSetup{
		OrgAdmin:   "admin",
		OrgName:    "org1",
		ConfigFile: "../config.yaml",
		ChannelID:  "myc",
	}

	// Initialization of the Fabric SDK from the previously set properties
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
	}
	fmt.Println("Successfully connected to Fabric network")

	b := users.BcNetwork{Fabric: &fSetup}

	// Create New user in Fabric network calling fabric-ca
	newUser, err := b.CreateUser()
	if err != nil {
		fmt.Println("Unable to create a user in the fabric-ca %v\n", err)
	}

	fmt.Println("User created!: %v\n", newUser)

}
