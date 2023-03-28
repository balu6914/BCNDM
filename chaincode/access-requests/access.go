package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const indexAccess = "cn~access"

var _ Service = (*accessChaincode)(nil)

type accessChaincode struct{}

// NewService returns access request chaincode implementation.
func NewService() Service {
	return accessChaincode{}
}

func (ac accessChaincode) RequestAccess(stub shim.ChaincodeStubInterface, receiver string) error {
	requester, err := callerCN(stub)
	if err != nil {
		return err
	}

	key, err := stub.CreateCompositeKey(indexAccess, []string{requester, receiver})
	if err != nil {
		return ErrFailedKeyCreation
	}

	val, err := stub.GetState(key)
	if err != nil {
		return ErrGettingState
	}

	if val != nil && State(val) != Revoked {
		return ErrConflict
	}

	if err := stub.PutState(key, []byte(Pending)); err != nil {
		return ErrSettingState
	}

	return nil
}

func (ac accessChaincode) ApproveAccess(stub shim.ChaincodeStubInterface, requester string) error {
	receiver, err := callerCN(stub)
	if err != nil {
		return err
	}

	key, err := stub.CreateCompositeKey(indexAccess, []string{requester, receiver})
	if err != nil {
		return ErrFailedKeyCreation
	}

	val, err := stub.GetState(key)
	if err != nil {
		return ErrGettingState
	}

	if val == nil {
		return ErrNotFound
	}

	if State(val) != Pending {
		return ErrInvalidStateTransition
	}

	if err := stub.PutState(key, []byte(Approved)); err != nil {
		return ErrSettingState
	}

	return nil
}

func (ac accessChaincode) RevokeAccess(stub shim.ChaincodeStubInterface, requester string) error {
	receiver, err := callerCN(stub)
	if err != nil {
		return err
	}

	key, err := stub.CreateCompositeKey(indexAccess, []string{requester, receiver})
	if err != nil {
		return ErrFailedKeyCreation
	}

	val, err := stub.GetState(key)
	if err != nil {
		return ErrGettingState
	}

	if val == nil {
		return ErrNotFound
	}

	state := State(val)
	if state != Approved && state != Pending {
		return fmt.Errorf("%w: %s", ErrInvalidStateTransition, state)
	}

	if err := stub.PutState(key, []byte(Revoked)); err != nil {
		return ErrSettingState
	}

	return nil
}

func (ac accessChaincode) ListAccess(stub shim.ChaincodeStubInterface) ([]Access, error) {
	receiver, err := callerCN(stub)
	if err != nil {
		return []Access{}, err
	}

	iter, err := stub.GetStateByPartialCompositeKey(indexAccess, []string{receiver})
	if err != nil {
		return []Access{}, ErrGettingState
	}
	defer iter.Close()

	res := []Access{}
	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			continue
		}

		_, participants, err := stub.SplitCompositeKey(kv.GetKey())
		if err != nil {
			continue
		}

		if len(participants) != 2 {
			continue
		}

		if participants[1] != receiver {
			continue
		}

		val := string(kv.GetValue())
		res = append(res, Access{
			requester: participants[0],
			receiver:  participants[1],
			state:     State(val),
		})
	}

	return res, nil
}

func (ac accessChaincode) GrantAccess(stub shim.ChaincodeStubInterface, dst string) error {
	src, err := callerCN(stub)
	if err != nil {
		return err
	}

	key, err := stub.CreateCompositeKey(indexAccess, []string{dst, src})
	if err != nil {
		return ErrFailedKeyCreation
	}

	if err := stub.PutState(key, []byte(Approved)); err != nil {
		return ErrSettingState
	}

	return nil
}
