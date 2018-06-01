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
func Balance(name string, fabric Fabric) (uint64, error) {
	// ClientContext allows creation of transactions using the supplied identity as the credential.
	adminChannelContext := fabric.Sdk.ChannelContext(fabric.ChannelID,
		fabsdk.WithUser(fabric.OrgAdmin), fabsdk.WithOrg(fabric.OrgName))

	client, err := channel.New(adminChannelContext)
	if err != nil {
		return 0, fmt.Errorf("Error on admin channel context init")
	}

	ub := userBalance{
		User: name,
	}

	balanceRqBytes, _ := json.Marshal(ub)
	balance, err := client.Query(channel.Request{
		ChaincodeID: fabric.ChaincodeID,
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
