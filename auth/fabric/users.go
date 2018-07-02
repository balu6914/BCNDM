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
func (fn *fabricNetwork) CreateUser(user *auth.User) error {

	sdk := fn.setup.Sdk
	ctxProvider := sdk.Context()
	mspClient, err := msp.New(ctxProvider)
	if err != nil {
		return fmt.Errorf("MSP client init failed: %v\n", err)
	}

	// Register the new user
	enrollmentSecret, err := mspClient.Register(&msp.RegistrationRequest{
		Name:        user.ID.Hex(),
		Affiliation: affiliation,
		Secret:      user.Password,
	})
	if err != nil {
		return fmt.Errorf("Registration failed: %v\n", err)
	}

	// Enroll the new user
	err = mspClient.Enroll(user.ID.Hex(), msp.WithSecret(enrollmentSecret))
	if err != nil {
		return fmt.Errorf("Enroll failed: %v\n", err)
	}

	// Get the new user's signing identity
	si, err := mspClient.GetSigningIdentity(user.ID.Hex())
	if err != nil {
		return fmt.Errorf("Unable to create a user in the fabric-ca, GetSigningIdentity failed: %v\n", err)
	}

	// TODO: Private key is "not supported". Understand why.
	user.PubCert = si.EnrollmentCertificate()
	user.PrivCert, _ = si.PrivateKey().Bytes()

	return nil
}
