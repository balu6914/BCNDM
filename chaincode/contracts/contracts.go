package main

import (
	"encoding/json"
	"math"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	objectType     = "contract"
	format         = "2006-01-02T15:04:05"
	decimals       = 5
	dateTimeFormat = "02-01-2006 15:04:05"
)

var _ Service = (*contractChaincode)(nil)

type contractChaincode struct {
	ts TransferService
}

// NewService returns contract chaincode implementation.
func NewService(ts TransferService) Service {
	return contractChaincode{ts: ts}
}

func (cc contractChaincode) CreateContracts(stub shim.ChaincodeStubInterface, contracts ...Contract) error {
	for _, contract := range contracts {
		dbContract := Contract{
			OwnerID:   contract.OwnerID,
			StartTime: contract.StartTime,
			PartnerID: contract.PartnerID,
			Share:     contract.Share,
			Signed:    false,
		}
		payload, err := json.Marshal(dbContract)
		if err != nil {
			return ErrInvalidArgument
		}

		endTime := contract.EndTime.Format(format)
		key, err := stub.CreateCompositeKey(objectType, []string{contract.StreamID, endTime, contract.PartnerID})
		if err != nil {
			return ErrFailedKeyCreation
		}

		if err := stub.PutState(key, payload); err != nil {
			return ErrSettingState
		}
	}

	return nil
}

func (cc contractChaincode) SignContract(stub shim.ChaincodeStubInterface, signedContract Contract) error {
	endTime := signedContract.EndTime.Format(format)
	key, err := stub.CreateCompositeKey(objectType, []string{signedContract.StreamID, endTime, signedContract.PartnerID})
	if err != nil {
		return ErrFailedKeyCreation
	}

	data, err := stub.GetState(key)
	if err != nil {
		return ErrGettingState
	}
	if data == nil {
		return ErrNotFound
	}

	var contract Contract
	if err := json.Unmarshal(data, &contract); err != nil {
		return ErrInvalidStateData
	}
	contract.Signed = true

	payload, err := json.Marshal(contract)
	if err != nil {
		return ErrGettingState
	}

	if err := stub.PutState(key, payload); err != nil {
		return ErrSettingState
	}

	return nil
}

func (cc contractChaincode) Transfer(stub shim.ChaincodeStubInterface, stream string, owner string, dateTime string, currentTime time.Time, value uint64) error {
	iter, err := stub.GetStateByPartialCompositeKey(objectType, []string{stream})
	if err != nil {
		return ErrGettingState
	}
	defer iter.Close()

	wholeValue := uint64(math.Pow(10, decimals))
	ownerValue := value

	transfers := []Transfer{}
	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			return ErrGettingState
		}

		var contract Contract
		if err := json.Unmarshal(kv.GetValue(), &contract); err != nil {
			return ErrGettingState
		}

		_, keys, err := stub.SplitCompositeKey(kv.GetKey())
		if err != nil {
			return ErrGettingState
		}

		contract.StreamID = keys[0]
		contract.EndTime, err = time.Parse(format, keys[1])
		if err != nil {
			return ErrGettingState
		}
		contract.PartnerID = keys[2]

		if contract.EndTime.Before(currentTime) || !contract.Signed {
			continue
		}

		transfer := Transfer{
			To:       contract.PartnerID,
			Value:    contract.Share * value / wholeValue,
			DateTime: dateTime,
		}
		transfers = append(transfers, transfer)
		ownerValue -= transfer.Value
	}

	return cc.ts.Transfer(stub, owner, dateTime, ownerValue, transfers...)
}
