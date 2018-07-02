package auth

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// FabricSetup with configuration parameters.
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

// FabricNetwork specifies an fabric account persistence API.
type FabricNetwork interface {
	Initialize() error
	CreateUser(*User) error
	Balance(string) (uint64, error)
}
