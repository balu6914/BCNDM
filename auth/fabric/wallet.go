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

const balanceFcn = "balance"

// Returns the account balance of another account with address user.
func (fn *fabricNetwork) Balance(name string) (uint64, error) {
	// ClientContext allows creation of transactions using the supplied identity as the credential.
	adminChannelContext := fn.setup.Sdk.ChannelContext(fn.setup.ChannelID,
		fabsdk.WithUser(fn.setup.OrgAdmin), fabsdk.WithOrg(fn.setup.OrgName))

	client, err := channel.New(adminChannelContext)
	if err != nil {
		return 0, fmt.Errorf("Error on admin channel context init: %v\n", err)
	}

	ub := userBalance{
		User: name,
	}

	balanceRqBytes, _ := json.Marshal(ub)
	balance, err := client.Query(channel.Request{
		ChaincodeID: fn.setup.ChaincodeID,
		Fcn:         balanceFcn,
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
