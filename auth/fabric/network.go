package fabric

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// Initialize reads the configuration file and sets up the client, chain and event hub
func (fn *fabricNetwork) Initialize() error {
	// Add parameters for the initialization
	if fn.setup.Initialized {
		return fmt.Errorf("sdk already initialized")
	}

	// Initialize the SDK with the configuration file
	sdk, err := fabsdk.New(config.FromFile(fn.setup.ConfigFile))
	if err != nil {
		return fmt.Errorf("failed to create sdk: %v", err)
	}

	fn.setup.Sdk = sdk
	fn.setup.Initialized = true
	return nil
}
