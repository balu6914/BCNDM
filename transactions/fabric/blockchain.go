package fabric

import (
	"encoding/json"
	"fmt"
	log "monetasa/logger"
	"monetasa/transactions"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	affiliation     = "org1"
	balanceFcn      = "balanceOf"
	transferFromFcn = "transferFrom"
	transferFcn     = "transfer"
	chanID          = "myc"
)

var _ transactions.BlockchainNetwork = (*fabricNetwork)(nil)

type fabricNetwork struct {
	sdk         *fabsdk.FabricSDK
	admin       string
	org         string
	chaincodeID string
	logger      log.Logger
}

// NewNetwork returns Fabric instance of blockchain network.
func NewNetwork(sdk *fabsdk.FabricSDK, admin, org, chaincodeID string, logger log.Logger) transactions.BlockchainNetwork {
	return fabricNetwork{
		sdk:         sdk,
		admin:       admin,
		org:         org,
		chaincodeID: chaincodeID,
		logger:      logger,
	}
}

func (fn fabricNetwork) CreateUser(id, secret string) error {
	ctx := fn.sdk.Context()
	mspClient, err := msp.New(ctx)
	if err != nil {
		fn.logger.Warn(fmt.Sprintf("failed to create msp client: %s", err))
		return err
	}

	es, err := mspClient.Register(&msp.RegistrationRequest{
		Name:        id,
		Affiliation: affiliation,
		Secret:      secret,
	})
	if err != nil {
		fn.logger.Warn(fmt.Sprintf("failed to register user: %s", err))
		return err
	}

	if err := mspClient.Enroll(id, msp.WithSecret(es)); err != nil {
		fn.logger.Warn(fmt.Sprintf("failed to enroll user: %s", err))
		return err
	}

	return nil
}

func (fn fabricNetwork) Balance(userID string) (uint64, error) {
	ctx := fn.sdk.ChannelContext(
		chanID,
		fabsdk.WithUser(fn.admin),
		fabsdk.WithOrg(fn.org),
	)

	client, err := channel.New(ctx)
	if err != nil {
		fn.logger.Warn(fmt.Sprintf("failed to create channel client: %s", err))
		return 0, err
	}

	req := balanceReq{Owner: userID}

	data, err := json.Marshal(req)
	if err != nil {
		fn.logger.Warn(fmt.Sprintf("failed to serialize balance request: %s", err))
		return 0, err
	}

	balance, err := client.Query(channel.Request{
		ChaincodeID: fn.chaincodeID,
		Fcn:         balanceFcn,
		Args:        [][]byte{data},
	})
	if err != nil {
		fn.logger.Warn(fmt.Sprintf("failed to query blockchain for balance: %s", err))
		return 0, err
	}

	var res balanceRes
	if err := json.Unmarshal(balance.Payload, &res); err != nil {
		fn.logger.Warn(fmt.Sprintf("failed to deserialize balance payload: %s", err))
		return 0, err
	}

	return res.Value, nil
}

func (fn fabricNetwork) Transfer(from, to string, value uint64) error {
	ctx := fn.sdk.ChannelContext(
		chanID,
		fabsdk.WithUser(fn.admin),
		fabsdk.WithOrg(fn.org),
	)

	client, err := channel.New(ctx)
	if err != nil {
		fn.logger.Warn(fmt.Sprintf("failed to create channel client: %s", err))
		return err
	}

	req := transferFromReq{
		From:  from,
		To:    to,
		Value: value,
	}

	data, err := json.Marshal(req)
	if err != nil {
		fn.logger.Warn(fmt.Sprintf("failed to serialize transfer_from request: %s", err))
		return err
	}

	_, err = client.Execute(channel.Request{
		ChaincodeID: fn.chaincodeID,
		Fcn:         transferFromFcn,
		Args:        [][]byte{data},
	})
	if err != nil {
		fn.logger.Warn(fmt.Sprintf("failed to execute transfer_from chaincode: %s", err))
		return err
	}

	return nil
}

func (fn fabricNetwork) BuyTokens(account string, value uint64) error {
	return fn.transfer(fn.admin, account, value)
}

func (fn fabricNetwork) WithdrawTokens(account string, value uint64) error {
	return fn.transfer(account, fn.admin, value)
}

func (fn fabricNetwork) transfer(from, to string, value uint64) error {
	ctx := fn.sdk.ChannelContext(
		chanID,
		fabsdk.WithUser(from),
		fabsdk.WithOrg(fn.org),
	)

	client, err := channel.New(ctx)
	if err != nil {
		fn.logger.Warn(fmt.Sprintf("failed to create channel client: %s", err))
		return err
	}

	req := transferReq{
		To:    to,
		Value: value,
	}

	data, err := json.Marshal(req)
	if err != nil {
		fn.logger.Warn(fmt.Sprintf("failed to serialize transfer request: %s", err))
		return err
	}

	_, err = client.Execute(channel.Request{
		ChaincodeID: fn.chaincodeID,
		Fcn:         transferFcn,
		Args:        [][]byte{data},
	})
	if err != nil {
		fn.logger.Warn(fmt.Sprintf("failed to execute transfer chaincode: %s", err))
		return err
	}

	return nil
}
