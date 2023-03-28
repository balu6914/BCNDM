package fabric

import (
	"encoding/json"
	"fmt"

	access "github.com/datapace/datapace/access-control"
	log "github.com/datapace/datapace/logger"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	requestFcn = "requestAccess"
	approveFcn = "approveAccess"
	revokeFcn  = "revokeAccess"
	grantFcn   = "grantAccess"
	chanID     = "datapacechannel"
)

var _ access.RequestLedger = (*accessRequestLedger)(nil)

type accessRequestLedger struct {
	sdk       *fabsdk.FabricSDK
	admin     string
	org       string
	chaincode string
	logger    log.Logger
}

// NewRequestLedger returns Fabric instance of access control ledger.
func NewRequestLedger(sdk *fabsdk.FabricSDK, admin, org, chaincode string, logger log.Logger) access.RequestLedger {
	return accessRequestLedger{
		sdk:       sdk,
		admin:     admin,
		org:       org,
		chaincode: chaincode,
		logger:    logger,
	}
}

func (arl accessRequestLedger) RequestAccess(sender, receiver string) error {
	ctx := arl.sdk.ChannelContext(
		chanID,
		fabsdk.WithUser(sender),
		fabsdk.WithOrg(arl.org),
	)

	client, err := channel.New(ctx)
	if err != nil {
		arl.logger.Warn(fmt.Sprintf("failed to create channel client: %s", err))
		return err
	}

	req := accessReq{
		Receiver: receiver,
	}

	data, err := json.Marshal(req)
	if err != nil {
		arl.logger.Warn(fmt.Sprintf("failed to serialize access request: %s", err))
		return err
	}

	_, err = client.Execute(channel.Request{
		ChaincodeID: arl.chaincode,
		Fcn:         requestFcn,
		Args:        [][]byte{data},
	})
	if err != nil {
		arl.logger.Warn(fmt.Sprintf("failed to execute access control chaincode: %s", err))
		return err
	}

	return nil
}

func (arl accessRequestLedger) Approve(approver, sender string) error {
	ctx := arl.sdk.ChannelContext(
		chanID,
		fabsdk.WithUser(approver),
		fabsdk.WithOrg(arl.org),
	)

	client, err := channel.New(ctx)
	if err != nil {
		arl.logger.Warn(fmt.Sprintf("failed to create channel client: %s", err))
		return err
	}

	req := approveReq{
		Requester: sender,
	}

	data, err := json.Marshal(req)
	if err != nil {
		arl.logger.Warn(fmt.Sprintf("failed to serialize approve request: %s", err))
		return err
	}

	_, err = client.Execute(channel.Request{
		ChaincodeID: arl.chaincode,
		Fcn:         approveFcn,
		Args:        [][]byte{data},
	})
	if err != nil {
		arl.logger.Warn(fmt.Sprintf("failed to execute access control chaincode: %s", err))
		return err
	}

	return nil
}

func (arl accessRequestLedger) Revoke(revoker, sender string) error {
	ctx := arl.sdk.ChannelContext(
		chanID,
		fabsdk.WithUser(revoker),
		fabsdk.WithOrg(arl.org),
	)

	client, err := channel.New(ctx)
	if err != nil {
		arl.logger.Warn(fmt.Sprintf("failed to create channel client: %s", err))
		return err
	}

	req := revokeReq{
		Requester: sender,
	}

	data, err := json.Marshal(req)
	if err != nil {
		arl.logger.Warn(fmt.Sprintf("failed to serialize revoke request: %s", err))
		return err
	}

	_, err = client.Execute(channel.Request{
		ChaincodeID: arl.chaincode,
		Fcn:         revokeFcn,
		Args:        [][]byte{data},
	})
	if err != nil {
		arl.logger.Warn(fmt.Sprintf("failed to execute access control chaincode: %s", err))
		return err
	}

	return nil
}

func (arl accessRequestLedger) Grant(src, dst string) error {
	ctx := arl.sdk.ChannelContext(
		chanID,
		fabsdk.WithUser(src),
		fabsdk.WithOrg(arl.org),
	)

	client, err := channel.New(ctx)
	if err != nil {
		arl.logger.Warn(fmt.Sprintf("failed to create channel client: %s", err))
		return err
	}

	req := grantReq{
		Destination: dst,
	}

	data, err := json.Marshal(req)
	if err != nil {
		arl.logger.Warn(fmt.Sprintf("failed to serialize grant request: %s", err))
		return err
	}

	_, err = client.Execute(channel.Request{
		ChaincodeID: arl.chaincode,
		Fcn:         grantFcn,
		Args:        [][]byte{data},
	})
	if err != nil {
		arl.logger.Warn(fmt.Sprintf("failed to execute access control chaincode: %s", err))
		return err
	}

	return nil
}
