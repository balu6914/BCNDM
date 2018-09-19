package fabric

import (
	"encoding/json"
	"fmt"
	log "monetasa/logger"
	"monetasa/transactions"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	createFcn = "createContracts"
	signFcn   = "signContract"
)

var _ transactions.ContractLedger = (*contractLedger)(nil)

type contractLedger struct {
	sdk       *fabsdk.FabricSDK
	admin     string
	org       string
	chaincode string
	logger    log.Logger
}

// NewContractLedger returns Fabric instance of contract ledger.
func NewContractLedger(sdk *fabsdk.FabricSDK, admin, org, chaincode string, logger log.Logger) transactions.ContractLedger {
	return contractLedger{
		sdk:       sdk,
		admin:     admin,
		org:       org,
		chaincode: chaincode,
		logger:    logger,
	}
}

func (cl contractLedger) Create(contracts ...transactions.Contract) error {
	if len(contracts) == 0 {
		return nil
	}

	ctx := cl.sdk.ChannelContext(
		chanID,
		fabsdk.WithUser(contracts[0].OwnerID),
		fabsdk.WithOrg(cl.org),
	)

	client, err := channel.New(ctx)
	if err != nil {
		cl.logger.Warn(fmt.Sprintf("failed to create channel client: %s", err))
		return err
	}

	req := []createContractReq{}
	for _, contract := range contracts {
		req = append(req, createContractReq{
			StreamID:  contract.StreamID,
			StartTime: contract.StartTime,
			EndTime:   contract.EndTime,
			OwnerID:   contract.OwnerID,
			PartnerID: contract.PartnerID,
			Share:     contract.Share,
		})
	}

	data, err := json.Marshal(req)
	if err != nil {
		cl.logger.Warn(fmt.Sprintf("failed to serialize create contract request: %s", err))
		return err
	}

	_, err = client.Execute(channel.Request{
		ChaincodeID: cl.chaincode,
		Fcn:         createFcn,
		Args:        [][]byte{data},
	})
	if err != nil {
		cl.logger.Warn(fmt.Sprintf("failed to execute create contract chaincode: %s", err))
		return err
	}

	return nil
}

func (cl contractLedger) Sign(contract transactions.Contract) error {
	ctx := cl.sdk.ChannelContext(
		chanID,
		fabsdk.WithUser(contract.PartnerID),
		fabsdk.WithOrg(cl.org),
	)

	client, err := channel.New(ctx)
	if err != nil {
		cl.logger.Warn(fmt.Sprintf("failed to create channel client: %s", err))
		return err
	}

	req := signContractReq{
		StreamID: contract.StreamID,
		EndTime:  contract.EndTime,
	}

	data, err := json.Marshal(req)
	if err != nil {
		cl.logger.Warn(fmt.Sprintf("failed to serialize sign contract request: %s", err))
		return err
	}

	_, err = client.Execute(channel.Request{
		ChaincodeID: cl.chaincode,
		Fcn:         signFcn,
		Args:        [][]byte{data},
	})
	if err != nil {
		cl.logger.Warn(fmt.Sprintf("failed to execute create contract chaincode: %s", err))
		return err
	}

	return nil
}
