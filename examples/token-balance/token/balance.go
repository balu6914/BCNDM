package token

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type Balance struct {
	User  string `json:"user"`
	Value uint64 `json:"value"`
}

// Returns the account balance of another account with address user.
func (bc *BcNetwork) Balance(name string) (b []byte, err error) {

	// ClientContext allows creation of transactions using the supplied identity as the credential.
	adminChannelContext := bc.Fabric.Sdk.ChannelContext(
		bc.Fabric.ChannelID,
		fabsdk.WithUser(bc.Fabric.OrgAdmin), fabsdk.WithOrg(bc.Fabric.OrgName))

	client, err := channel.New(adminChannelContext)

	if err != nil {
		fmt.Println("Error init SDK")
		return nil, err
	}

	// Build chaincode arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "balance") // Chaincode function name
	args = append(args, "user")    // Chaincode fn params
	args = append(args, name)      // Chaincode fn params

	balance, err := client.Query(channel.Request{
		ChaincodeID: bc.Fabric.ChannelID,
		Fcn:         args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4])},
	})

	if err != nil {
		fmt.Println("Error fetching balance")
		return nil, err

	}

	return balance.Payload, nil

}
