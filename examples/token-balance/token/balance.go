package token

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric/common/util"
)

type Balance struct {
	User  string `json:"user"`
	Value uint64 `json:"value"`
}

// Returns the account balance of another account with address user.
func (bc *BcNetwork) Balance(name string) (b Balance, err error) {

	// ClientContext allows creation of transactions using the supplied identity as the credential.
	adminChannelContext := bc.Fabric.Sdk.ChannelContext(
		bc.Fabric.ChannelID,
		fabsdk.WithUser(bc.Fabric.OrgAdmin), fabsdk.WithOrg(bc.Fabric.OrgName))

	client, err := channel.New(adminChannelContext)

	if err != nil {
		fmt.Println("Error init SDK")
		return b, err
	}

	balanceRq := Balance{User: name}
	balanceRqBytes, _ := json.Marshal(balanceRq)

	balance, err := client.Query(channel.Request{
		ChaincodeID: bc.Fabric.ChaincodeID,
		Fcn:         "balance",
		Args:        util.ToChaincodeArgs(string(balanceRqBytes))})

	if err != nil {
		fmt.Println("Error fetching balance!!!")
		return b, err

	}

	result := Balance{}
	err = json.Unmarshal(balance.Payload, &result)

	return result, nil

}
