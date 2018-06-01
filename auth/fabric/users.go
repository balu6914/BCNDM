package fabric

import (
	"fmt"
	"monetasa/auth/fabric/blockchain"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
)

// Add new user to fabric network
func CreateUser(name, secret string) (mspctx.SigningIdentity, error) {

	fSetup := blockchain.FabricSetup{
		OrgAdmin:   "admin",
		OrgName:    "org1",
		ConfigFile: os.Getenv("GOPATH") + "/src/monetasa/examples/config/config.yaml",
		ChannelID:  "myc",
	}
	// Initialization of the Fabric SDK from the previously set properties
	if err := fSetup.Initialize(); err != nil {
		return nil, fmt.Errorf("Unable to initialize the Fabric SDK: %v\n", err)
	}

	bc := BcNetwork{Fabric: &fSetup}
	sdk := bc.Fabric.Sdk
	ctxProvider := sdk.Context()
	mspClient, err := msp.New(ctxProvider)
	if err != nil {
		return nil, fmt.Errorf("MSP client init failed: %v\n", err)
	}

	// Register the new user
	enrollmentSecret, err := mspClient.Register(&msp.RegistrationRequest{
		Name:        name,
		Affiliation: "org1",
		Secret:      secret,
	})
	if err != nil {
		return nil, fmt.Errorf("Registration failed: %v\n", err)
	}

	// Enroll the new user
	err = mspClient.Enroll(name, msp.WithSecret(enrollmentSecret))
	if err != nil {
		return nil, fmt.Errorf("Enroll failed: %v\n", err)
	}

	// Get the new user's signing identity
	si, err := mspClient.GetSigningIdentity(name)
	if err != nil {
		return nil, fmt.Errorf("GetSigningIdentity failed: %v\n", err)
	}

	return si, nil
}
