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
	affiliation = "org1"
	balanceFcn  = "balance"
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

func (fn fabricNetwork) CreateUser(id, secret string) ([]byte, error) {
	ctx := fn.sdk.Context()
	mspClient, err := msp.New(ctx)
	if err != nil {
		fn.logger.Warn(fmt.Sprintf("failed to create msp client: %s", err))
		return []byte{}, err
	}

	es, err := mspClient.Register(&msp.RegistrationRequest{
		Name:        id,
		Affiliation: affiliation,
		Secret:      secret,
	})
	if err != nil {
		fn.logger.Warn(fmt.Sprintf("failed to register user: %s", err))
		return []byte{}, err
	}

	if err := mspClient.Enroll(id, msp.WithSecret(es)); err != nil {
		fn.logger.Warn(fmt.Sprintf("failed to enroll user: %s", err))
		return []byte{}, err
	}

	si, err := mspClient.GetSigningIdentity(id)
	if err != nil {
		fn.logger.Warn(fmt.Sprintf("failed to get signing identity for user: %s", err))
		return []byte{}, err
	}

	return si.EnrollmentCertificate(), nil
}

func (fn fabricNetwork) Balance(userID, chanID string) (uint64, error) {
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

	ub := userBalance{User: userID}

	balanceReq, err := json.Marshal(ub)
	if err != nil {
		fn.logger.Warn(fmt.Sprintf("failed to serialize balance request: %s", err))
		return 0, err
	}

	balance, err := client.Query(channel.Request{
		ChaincodeID: fn.chaincodeID,
		Fcn:         balanceFcn,
		Args:        [][]byte{balanceReq},
	})
	if err != nil {
		fn.logger.Warn(fmt.Sprintf("failed to query blockchain for balance: %s", err))
		return 0, err
	}

	if err := json.Unmarshal(balance.Payload, &ub); err != nil {
		fn.logger.Warn(fmt.Sprintf("failed to deserialize balance payload: %s", err))
		return 0, err
	}

	return ub.Value, nil
}

type userBalance struct {
	User  string `json:"user"`
	Value uint64 `json:"value"`
}
