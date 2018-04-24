package main

import (
	"fmt"
	"monetasa/examples/blockchain"
	"monetasa/examples/fabric-ca/users"
)

func main() {

	// Definition of the Fabric SDK properties
	fSetup := blockchain.FabricSetup{
		OrgAdmin:   "Admin",
		OrgName:    "org1",
		ConfigFile: "../config.yaml",
		// Channel parameters
		ChannelID: "myc",
	}

	/**
	 * Initialization of the Fabric SDK from the previously set properties
	 */
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
	}
	fmt.Println("Successfully connected to Fabric network")

	b := users.BcNetwork{Fabric: &fSetup}

	/**
	 * Create New user in Fabric network calling fabric-ca
	 * [err fabric-ca err response]
	 * user user Object
	 */
	newUser, err := b.CreateUser()
	if err != nil {
		fmt.Println("Unable to create a user in the fabric-ca %v\n", err)
	}

	fmt.Println("User created!: %v\n", newUser)

	/**
	 * Set the user in context
	 */

	//userContext := fSetup.Sdk.Context(fabsdk.WithUser("test1"), fabsdk.WithOrg(fSetup.OrgName))

}
