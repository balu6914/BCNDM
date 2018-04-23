package users

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
)

// Add new user to fabric network
func (bc *BcNetwork) CreateUser() (usr mspctx.SigningIdentity, err error) {

	sdk := bc.Fabric.Sdk

	ctxProvider := sdk.Context()

	mspClient, err := msp.New(ctxProvider)

	if err != nil {
		return nil, err
		fmt.Println("MSP client init failed: %v", err)
	}

	// Register the new user
	enrollmentSecret, err := mspClient.Register(&msp.RegistrationRequest{
		Name:        "test1",
		Affiliation: "org1",
		Secret:      "12345",
	})

	if err != nil {
		return nil, err
		fmt.Println("Registration failed: %v", err)
	}

	// Enroll the new user
	err = mspClient.Enroll("test", msp.WithSecret(enrollmentSecret))

	if err != nil {
		fmt.Println("Enroll failed: %v", err)
		return nil, err
	}

	// Get the new user's signing identity
	si, err := mspClient.GetSigningIdentity("test")
	if err != nil {
		fmt.Println("GetSigningIdentity failed: %v", err)
		return nil, err
	}

	return si, nil
}
