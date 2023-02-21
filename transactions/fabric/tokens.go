package fabric

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	log "github.com/datapace/datapace/logger"
	"github.com/datapace/datapace/transactions"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	affiliation    = "org1"
	balanceFcn     = "balanceOf"
	transferFcn    = "transfer"
	chanID         = "datapacechannel"
	txHistoryFcn   = "txHistory"
	dateTimeFormat = "02-01-2006 15:04:05"
)

var _ transactions.TokenLedger = (*tokenLedger)(nil)

type tokenLedger struct {
	sdk               *fabsdk.FabricSDK
	admin             string
	org               string
	tokenChaincode    string
	contractChaincode string
	logger            log.Logger
}

// NewTokenLedger returns Fabric instance of token ledger.
func NewTokenLedger(sdk *fabsdk.FabricSDK, admin, org, token, contract string, logger log.Logger) transactions.TokenLedger {
	return tokenLedger{
		sdk:               sdk,
		admin:             admin,
		org:               org,
		tokenChaincode:    token,
		contractChaincode: contract,
		logger:            logger,
	}
}

func (tl tokenLedger) CreateUser(id, secret string) error {
	ctx := tl.sdk.Context()
	mspClient, err := msp.New(ctx)
	if err != nil {
		tl.logger.Warn(fmt.Sprintf("failed to create msp client: %s", err))
		return err
	}

	es, err := mspClient.Register(&msp.RegistrationRequest{
		Name:        id,
		Affiliation: affiliation,
		Secret:      secret,
	})
	if err != nil {
		tl.logger.Warn(fmt.Sprintf("failed to register user: %s", err))
		return err
	}

	if err := mspClient.Enroll(id, msp.WithSecret(es)); err != nil {
		tl.logger.Warn(fmt.Sprintf("failed to enroll user: %s", err))
		return err
	}

	return nil
}

func (tl tokenLedger) Balance(userID string) (uint64, error) {
	ctx := tl.sdk.ChannelContext(
		chanID,
		fabsdk.WithUser(tl.admin),
		fabsdk.WithOrg(tl.org),
	)

	client, err := channel.New(ctx)
	if err != nil {
		tl.logger.Warn(fmt.Sprintf("failed to create channel client: %s", err))
		return 0, err
	}

	req := balanceReq{Owner: userID}

	data, err := json.Marshal(req)
	if err != nil {
		tl.logger.Warn(fmt.Sprintf("failed to serialize balance request: %s", err))
		return 0, err
	}

	balance, err := client.Query(channel.Request{
		ChaincodeID: tl.tokenChaincode,
		Fcn:         balanceFcn,
		Args:        [][]byte{data},
	})
	if err != nil {
		tl.logger.Warn(fmt.Sprintf("failed to query blockchain for balance: %s", err))
		return 0, err
	}

	var res balanceRes
	if err := json.Unmarshal(balance.Payload, &res); err != nil {
		tl.logger.Warn(fmt.Sprintf("failed to deserialize balance payload: %s", err))
		return 0, err
	}

	return res.Value, nil
}

func (tl tokenLedger) Transfer(stream, from, to string, value uint64) error {

	ctx := tl.sdk.ChannelContext(
		chanID,
		fabsdk.WithUser(from),
		fabsdk.WithOrg(tl.org),
	)

	client, err := channel.New(ctx)
	if err != nil {
		tl.logger.Warn(fmt.Sprintf("failed to create channel client: %s", err))
		return err
	}

	req := transferReq{
		StreamID: stream,
		To:       to,
		Time:     time.Now(),
		Value:    value,
		DateTime: time.Now().Format(dateTimeFormat),
	}

	data, err := json.Marshal(req)
	if err != nil {
		tl.logger.Warn(fmt.Sprintf("failed to serialize transfer request: %s", err))
		return err
	}

	_, err = client.Execute(channel.Request{
		ChaincodeID: tl.contractChaincode,
		Fcn:         transferFcn,
		Args:        [][]byte{data},
	})
	if err != nil {
		tl.logger.Warn(fmt.Sprintf("failed to execute transfer chaincode: %s", err))
		e, ok := status.FromError(err)
		if ok && strings.Contains(e.Message, transactions.ErrNotEnoughTokens.Error()) {
			return transactions.ErrNotEnoughTokens
		}
		return err
	}

	return nil
}

func (tl tokenLedger) BuyTokens(account string, value uint64) error {
	return tl.transfer(tl.admin, account, value)
}

func (tl tokenLedger) WithdrawTokens(account string, value uint64) error {
	return tl.transfer(account, tl.admin, value)
}

func (tl tokenLedger) transfer(from, to string, value uint64) error {

	ctx := tl.sdk.ChannelContext(
		chanID,
		fabsdk.WithUser(from),
		fabsdk.WithOrg(tl.org),
	)

	client, err := channel.New(ctx)
	if err != nil {
		tl.logger.Warn(fmt.Sprintf("failed to create channel client: %s", err))
		return err
	}

	req := transferReq{
		To:       to,
		Value:    value,
		DateTime: time.Now().Format(dateTimeFormat),
	}

	data, err := json.Marshal(req)
	if err != nil {
		tl.logger.Warn(fmt.Sprintf("failed to serialize transfer request: %s", err))
		return err
	}

	_, err = client.Execute(channel.Request{
		ChaincodeID: tl.tokenChaincode,
		Fcn:         transferFcn,
		Args:        [][]byte{data},
	})
	if err != nil {
		tl.logger.Warn(fmt.Sprintf("failed to execute transfer chaincode: %s", err))
		e, ok := status.FromError(err)
		if ok && strings.Contains(e.Message, transactions.ErrNotEnoughTokens.Error()) {
			return transactions.ErrNotEnoughTokens
		}
		return err
	}

	return nil
}

func (tl tokenLedger) TxHistory(userID, fromDateTime, toDateTime, txType string) (transactions.TokenTxHistory, error) {
	txHis := new(transactions.TokenTxHistory)
	ctx := tl.sdk.ChannelContext(
		chanID,
		fabsdk.WithUser(tl.admin),
		fabsdk.WithOrg(tl.org),
	)

	client, err := channel.New(ctx)
	if err != nil {
		tl.logger.Warn(fmt.Sprintf("failed to create channel client: %s", err))
		return *txHis, err
	}

	req := txHistoryReq{
		Owner:        userID,
		FromDateTime: fromDateTime,
		ToDateTime:   toDateTime,
		TxType:       txType,
	}

	data, err := json.Marshal(req)
	if err != nil {
		tl.logger.Warn(fmt.Sprintf("failed to serialize request: %s", err))
		return *txHis, err
	}

	txHistoryBytes, err := client.Query(channel.Request{
		ChaincodeID: tl.tokenChaincode,
		Fcn:         txHistoryFcn,
		Args:        [][]byte{data},
	})

	if err != nil {
		tl.logger.Warn(fmt.Sprintf("failed to query blockchain for tx history: %s", err))
		return *txHis, err
	}

	if err := json.Unmarshal(txHistoryBytes.Payload, txHis); err != nil {
		tl.logger.Warn(fmt.Sprintf("failed to deserialize balance payload: %s", err))
		return *txHis, err
	}

	return *txHis, nil
}
