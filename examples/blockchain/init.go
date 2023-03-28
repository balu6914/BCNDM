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
	ChaincodeID string
	Initialized bool
	OrgAdmin    string
	OrgName     string
	admin       resmgmt.Client
	Sdk         *fabsdk.FabricSDK
}

// Initialize reads the configuration file and sets up the client, chain and event hub
func (setup *FabricSetup) Initialize() error {

	// Add parameters for the initialization
	if setup.Initialized {
		return fmt.Errorf("sdk already initialized")
	}

	// Initialize the SDK with the configuration file
	sdk, err := fabsdk.New(config.FromFile(setup.ConfigFile))
	if err != nil {
		return fmt.Errorf("failed to create sdk: %v", err)
	}

	setup.Sdk = sdk
	setup.Initialized = true
	return nil
}
