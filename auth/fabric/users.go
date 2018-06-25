package fabric

import (
	"fmt"
	"monetasa/auth"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
)

const affiliation = "org1"

var _ auth.FabricNetwork = (*fabricNetwork)(nil)

type fabricNetwork struct {
	setup *auth.FabricSetup
}

// NewUserRepository instantiates a PostgreSQL implementation of user
// repository.
func NewFabricNetwork(fs *auth.FabricSetup) auth.FabricNetwork {
	return &fabricNetwork{fs}
}

// Add new user to fabric network
func (fn *fabricNetwork) CreateUser(name, secret string) error {
	sdk := fn.setup.Sdk
	ctxProvider := sdk.Context()
	mspClient, err := msp.New(ctxProvider)
	if err != nil {
		return fmt.Errorf("MSP client init failed: %v\n", err)
	}

	// Register the new user
	enrollmentSecret, err := mspClient.Register(&msp.RegistrationRequest{
		Name:        name,
		Affiliation: affiliation,
		Secret:      secret,
	})
	if err != nil {
		return fmt.Errorf("Registration failed: %v\n", err)
	}

	// Enroll the new user
	if err := mspClient.Enroll(name, msp.WithSecret(enrollmentSecret)); err != nil {
		return fmt.Errorf("Enroll failed: %v\n", err)
	}

	// Get the new user's signing identity
	si, err := mspClient.GetSigningIdentity(name)
	if err != nil {
		return fmt.Errorf("Unable to create a user in the fabric-ca, GetSigningIdentity failed: %v\n", err)
	}
	fmt.Printf("User created: %v\n", si)

	return nil
}
