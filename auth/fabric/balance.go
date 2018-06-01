package fabric

import (
	"encoding/json"
	"fmt"
	"monetasa/auth/fabric/blockchain"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric/common/util"
)

type userBalance struct {
	User  string `json:"user"`
	Value uint64 `json:"value"`
}

// Returns the account balance of another account with address user.
func Balance(name string) (uint64, error) {

	fSetup := blockchain.FabricSetup{
		OrgAdmin:    "admin",
		OrgName:     "Org1",
		ConfigFile:  os.Getenv("GOPATH") + "/src/monetasa/examples/config/config.yaml",
		ChannelID:   "myc",
		ChaincodeID: "token",
	}

	// Initialization of the Fabric SDK from the previously set properties
	if err := fSetup.Initialize(); err != nil {
		return 0, fmt.Errorf("Unable to initialize the Fabric SDK: %v\n", err)
	}

	bc := BcNetwork{Fabric: &fSetup}
	// ClientContext allows creation of transactions using the supplied identity as the credential.
	adminChannelContext := bc.Fabric.Sdk.ChannelContext(
		bc.Fabric.ChannelID,
		fabsdk.WithUser(bc.Fabric.OrgAdmin), fabsdk.WithOrg(bc.Fabric.OrgName))

	client, err := channel.New(adminChannelContext)
	if err != nil {
		fmt.Errorf("Error init SDK")
		return 0, err
	}

	ub := userBalance{
		User: name,
	}

	balanceRqBytes, _ := json.Marshal(ub)
	balance, err := client.Query(channel.Request{
		ChaincodeID: bc.Fabric.ChaincodeID,
		Fcn:         "balance",
		Args:        util.ToChaincodeArgs(string(balanceRqBytes)),
	})
	if err != nil {
		return 0, fmt.Errorf("Error fetching balance: %v\n", err)

	}

	err = json.Unmarshal(balance.Payload, &ub)
	if err != nil {
		return 0, fmt.Errorf("Error Unmarshaling balance: %v\n", err)

	}

	return ub.Value, nil

}
