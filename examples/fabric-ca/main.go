package main

import (
	"fmt"
	"monetasa/examples/blockchain"
	"monetasa/examples/fabric-ca/users"
)

// BC network instance
type bcNetwork struct {
	Fabric *blockchain.FabricSetup
}

func main() {

	// Definition of the Fabric SDK properties
	fSetup := blockchain.FabricSetup{
		OrgAdmin:   "Admin",
		OrgName:    "Org1",
		ConfigFile: "config.yaml",
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

	/**
	 * Create New user in Fabric network calling fabric-ca
	 * [err fabric-ca err response]
	 * user user Object
	 */
	user, err := users.Create(fSetup)
	if err != nil {
		fmt.Println("Unable to create a user in the fabric-ca %v\n", err)
	}

	fmt.Printf(user)
	fmt.Println("User created!: %v\n", err)

	/**
	 * Fetch this user from Fabric Network by passing a  email
	 */

	/**
	 * Set this user in context so we can sign transaction with his keys
	 */

}
