package users

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
)

// Add new user to fabric network
func (bc *BcNetwork) CreateUser(name, secret string) (mspctx.SigningIdentity, error) {

	sdk := bc.Fabric.Sdk
	ctxProvider := sdk.Context()
	mspClient, err := msp.New(ctxProvider)
	if err != nil {
		fmt.Printf("MSP client init failed: %v\n", err)
		return nil, err
	}

	// Register the new user
	enrollmentSecret, err := mspClient.Register(&msp.RegistrationRequest{
		Name:        name,
		Affiliation: "org1",
		Secret:      secret,
	})
	if err != nil {
		fmt.Printf("Registration failed: %v\n", err)
		return nil, err
	}

	// Enroll the new user
	err = mspClient.Enroll(name, msp.WithSecret(enrollmentSecret))
	if err != nil {
		fmt.Printf("Enroll failed: %v\n", err)
		return nil, err
	}

	// Get the new user's signing identity
	si, err := mspClient.GetSigningIdentity(name)
	if err != nil {
		fmt.Printf("GetSigningIdentity failed: %v\n", err)
		return nil, err
	}

	return si, nil
}
