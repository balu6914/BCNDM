package blockchain

import (
	"fmt"

	resmgmt "github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// FabricSetup implementation
type FabricSetup struct {
	ConfigFile  string
	ChannelID   string
	initialized bool
	OrgAdmin    string
	OrgName     string
	admin       resmgmt.Client
	sdk         *fabsdk.FabricSDK
}

// Initialize reads the configuration file and sets up the client, chain and event hub
func (setup *FabricSetup) Initialize() error {

	// Add parameters for the initialization
	if setup.initialized {
		return fmt.Errorf("sdk already initialized")
	}

	// Initialize the SDK with the configuration file
	sdk, err := fabsdk.New(config.FromFile(setup.ConfigFile))
	if err != nil {
		return fmt.Errorf("failed to create sdk: %v", err)
	}

	setup.sdk = sdk

	// ClientContext allows creation of transactions using the supplied identity as the credential.
	// We will need this to set specific user to context (e.g make transactions in his name etc...).
	//clientContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(ordererOrgName))
	setup.initialized = true
	return nil
}
