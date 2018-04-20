package blockchain

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
)

// Add new user to fabric network
func (bc *FabricSetup) CreateUser() (response string, err error) {

	sdk := bc.sdk

	ctxProvider := sdk.Context()

	mspClient, err := msp.New(ctxProvider)

	if err != nil {
		fmt.Println("MSP client init failed: %v", err)
	}

	// Register the new user
	enrollmentSecret, err := mspClient.Register(&msp.RegistrationRequest{
		Name: "test",
		// Affiliation is mandatory. "org1" and "org2" are hardcoded as CA defaults
		// See https://github.com/hyperledger/fabric-ca/blob/release/cmd/fabric-ca-server/config.go
		Affiliation: "org1",
		Secret:      "12345",
	})

	if err != nil {
		return "", err
		fmt.Println("Registration failed: %v", err)
	}

	// Enroll the new user
	err = mspClient.Enroll("test", msp.WithSecret(enrollmentSecret))

	if err != nil {
		fmt.Println("Enroll failed: %v", err)
		return "", err
	}

	// Get the new user's signing identity
	si, err := mspClient.GetSigningIdentity("test")
	if err != nil {
		fmt.Println("GetSigningIdentity failed: %v", err)
		return "", err
	}

	fmt.Println("YESSSS Here is si", si)

	return "Cool", nil
}
