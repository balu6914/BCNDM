package fabric

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric/common/util"
)

type userBalance struct {
	User  string `json:"user"`
	Value uint64 `json:"value"`
}

// Returns the account balance of another account with address user.
func Balance(name string, bcn BcNetwork) (uint64, error) {
	// ClientContext allows creation of transactions using the supplied identity as the credential.
	adminChannelContext := bcn.Fabric.Sdk.ChannelContext(
		bcn.Fabric.ChannelID,
		fabsdk.WithUser(bcn.Fabric.OrgAdmin), fabsdk.WithOrg(bcn.Fabric.OrgName))

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
		ChaincodeID: bcn.Fabric.ChaincodeID,
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
