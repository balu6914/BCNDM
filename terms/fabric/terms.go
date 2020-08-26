package fabric

import (
	"encoding/json"
	"fmt"
	log "github.com/datapace/datapace/logger"
	t "github.com/datapace/datapace/terms"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	createTermsFcn   = "storeTerms"
	validateTermsFcn = "validateTerms"
	chanID           = "datapacechannel"
)

var _ t.TermsLedger = (*createTermsLedger)(nil)

type createTermsLedger struct {
	sdk       *fabsdk.FabricSDK
	admin     string
	org       string
	chaincode string
	logger    log.Logger
}

// NewTermsLedger returns Fabric instance of terms ledger.
func NewTermsLedger(sdk *fabsdk.FabricSDK, admin, org, chaincode string, logger log.Logger) t.TermsLedger {
	return createTermsLedger{
		sdk:       sdk,
		admin:     admin,
		org:       org,
		chaincode: chaincode,
		logger:    logger,
	}
}

func (ctl createTermsLedger) CreateTerms(terms t.Terms) error {
	ctx := ctl.sdk.ChannelContext(
		chanID,
		fabsdk.WithUser(ctl.admin),
		fabsdk.WithOrg(ctl.org),
	)

	client, err := channel.New(ctx)
	if err != nil {
		ctl.logger.Warn(fmt.Sprintf("failed to create channel client: %s", err))
		return err
	}

	req := termsReq{
		StreamID:  terms.StreamID,
		TermsURL:  terms.TermsHash,
		TermsHash: terms.TermsHash,
	}

	data, err := json.Marshal(req)
	if err != nil {
		ctl.logger.Warn(fmt.Sprintf("failed to serialize create terms request: %s", err))
		return err
	}

	_, err = client.Execute(channel.Request{
		ChaincodeID: ctl.chaincode,
		Fcn:         createTermsFcn,
		Args:        [][]byte{data},
	})
	if err != nil {
		ctl.logger.Warn(fmt.Sprintf("failed to execute terms chaincode: %s", err))
		return err
	}
	return nil
}

func (ctl createTermsLedger) ValidateTerms(terms t.Terms) (bool, error) {
	ctx := ctl.sdk.ChannelContext(
		chanID,
		fabsdk.WithUser(ctl.admin),
		fabsdk.WithOrg(ctl.org),
	)

	client, err := channel.New(ctx)
	if err != nil {
		ctl.logger.Warn(fmt.Sprintf("failed to create channel client: %s", err))
		return false, err
	}

	req := termsReq{
		StreamID:  terms.StreamID,
		TermsURL:  terms.TermsHash,
		TermsHash: terms.TermsHash,
	}

	data, err := json.Marshal(req)
	if err != nil {
		ctl.logger.Warn(fmt.Sprintf("failed to serialize validate terms request: %s", err))
		return false, err
	}
	resp, err := client.Execute(channel.Request{
		ChaincodeID: ctl.chaincode,
		Fcn:         validateTermsFcn,
		Args:        [][]byte{data},
	})
	if err != nil {
		ctl.logger.Warn(fmt.Sprintf("failed to validate terms chaincode: %s", err))
		return false, err
	}
	var res validationRes
	if err := json.Unmarshal(resp.Payload, &res); err != nil {
		ctl.logger.Warn(fmt.Sprintf("failed to deserialize terms validation payload: %s", err))
		return false, err
	}
	return res.Valid, nil
}
