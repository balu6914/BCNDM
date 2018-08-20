package main

import (
	"encoding/binary"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	keyToken       = "__token"
	indexBalance   = "cn~balance"
	indexAllowance = "cn~allowance"
)

var _ Service = (*tokenChaincode)(nil)

type tokenChaincode struct{}

// NewService returns new ERC20 implementation instance.
func NewService() Service {
	return tokenChaincode{}
}

func (tc tokenChaincode) Init(stub shim.ChaincodeStubInterface) error {
	function, args := stub.GetFunctionAndParameters()

	if function != "init" {
		return ErrInvalidFuncCall
	}

	if len(args) != 1 {
		return ErrInvalidNumOfArgs
	}

	// get token data from JSON
	var ti tokenInfo
	if err := json.Unmarshal([]byte(args[0]), &ti); err != nil {
		return ErrInvalidArgument
	}

	if err := stub.PutState(keyToken, []byte(args[0])); err != nil {
		return ErrSettingState
	}

	// get caller CN from his certificate
	caller, err := callerCN(stub)
	if err != nil {
		return ErrUnauthorized
	}

	// set the balance using a helper function
	if err := setBalance(stub, caller, ti.TotalSupply); err != nil {
		return ErrFailedBalanceSet
	}

	return nil
}

func (tc tokenChaincode) TotalSupply(stub shim.ChaincodeStubInterface) (uint64, error) {
	data, err := stub.GetState(keyToken)
	if err != nil {
		return 0, ErrGettingState
	}

	var ti tokenInfo
	if err := json.Unmarshal(data, &ti); err != nil {
		return 0, ErrInvalidStateData
	}

	return ti.TotalSupply, nil
}

func (tc tokenChaincode) BalanceOf(stub shim.ChaincodeStubInterface, owner string) (uint64, error) {
	key, err := stub.CreateCompositeKey(indexBalance, []string{owner})
	if err != nil {
		return 0, ErrFailedKeyCreation
	}

	// if the user cn is not in the state, then the balance is 0
	data, err := stub.GetState(key)
	if err != nil {
		return 0, ErrGettingState
	}

	if data == nil {
		return 0, nil
	}

	return binary.LittleEndian.Uint64(data), nil
}

func (tc tokenChaincode) Transfer(stub shim.ChaincodeStubInterface, to string, value uint64) bool {
	from, err := callerCN(stub)
	if err != nil {
		return false
	}

	return tc.transfer(stub, from, to, value)
}

func (tc tokenChaincode) TransferFrom(stub shim.ChaincodeStubInterface, from, to string, value uint64) bool {
	spender, err := callerCN(stub)
	if err != nil {
		return false
	}

	allowance, err := tc.Allowance(stub, from, spender)
	if err != nil {
		allowance = 0
	}

	if allowance < value {
		return false
	}

	// allowance[_from][msg.sender] -= _value
	if err := setAllowance(stub, from, spender, allowance-value); err != nil {
		return false
	}

	return tc.transfer(stub, from, to, value)
}

func (tc tokenChaincode) Approve(stub shim.ChaincodeStubInterface, spender string, value uint64) bool {
	from, err := callerCN(stub)
	if err != nil {
		return false
	}

	// allowance[msg.sender][_spender] = _value
	if err := setAllowance(stub, from, spender, value); err != nil {
		return false
	}

	a := approve{
		Spender: spender,
		Value:   value,
	}
	approveData, err := json.Marshal(a)
	if err != nil {
		return false
	}

	err = stub.SetEvent("Approval", approveData)
	return err == nil
}

func (tc tokenChaincode) Allowance(stub shim.ChaincodeStubInterface, owner, spender string) (uint64, error) {
	key, err := stub.CreateCompositeKey(indexAllowance, []string{owner, spender})
	if err != nil {
		return 0, ErrFailedKeyCreation
	}

	data, err := stub.GetState(key)
	if err != nil {
		return 0, ErrGettingState
	}

	// if the key is not in the state, then the value is 0
	if data == nil {
		return 0, nil
	}

	return binary.LittleEndian.Uint64(data), nil
}

func (tc tokenChaincode) GroupTransfer(stub shim.ChaincodeStubInterface, transfers ...Transfer) error {
	from, err := callerCN(stub)
	if err != nil {
		return err
	}

	fromBalance, err := tc.BalanceOf(stub, from)
	if err != nil {
		return err
	}

	events := []transfer{}

	for _, tr := range transfers {
		if from == tr.To {
			continue
		}

		toBalance, err := tc.BalanceOf(stub, tr.To)
		if err != nil {
			return err
		}

		if fromBalance < tr.Value {
			return ErrNotEnoughTokens
		}

		if toBalance+tr.Value < toBalance {
			return ErrOverflow
		}

		fromBalance -= tr.Value
		if err := setBalance(stub, from, fromBalance); err != nil {
			return err
		}

		toBalance += tr.Value
		if err := setBalance(stub, tr.To, toBalance); err != nil {
			return err
		}

		events = append(events, transfer{
			From:  from,
			To:    tr.To,
			Value: tr.Value,
		})
	}

	payload, err := json.Marshal(events)
	if err != nil {
		return ErrFailedSerialization
	}

	if err := stub.SetEvent("Transfers", payload); err != nil {
		return ErrSettingState
	}

	return nil
}

func (tc tokenChaincode) transfer(stub shim.ChaincodeStubInterface, from, to string, value uint64) bool {
	if from == to {
		return true
	}

	// retrieving balances
	fromBalance, err := tc.BalanceOf(stub, from)
	if err != nil {
		return false
	}

	toBalance, err := tc.BalanceOf(stub, to)
	if err != nil {
		toBalance = 0
	}

	// if (balanceOf[_from] < _value) throw
	if fromBalance < value {
		return false
	}

	// if (balanceOf[_to] + _value < balanceOf[_to]) throw
	if toBalance+value < toBalance {
		return false
	}

	// balanceOf[_from] -= _value
	if err := setBalance(stub, from, fromBalance-value); err != nil {
		return false
	}

	// balanceOf[_to] += _value
	if err := setBalance(stub, to, toBalance+value); err != nil {
		return false
	}

	t := transfer{
		From:  from,
		To:    to,
		Value: value,
	}
	transferData, err := json.Marshal(t)
	if err != nil {
		return false
	}

	// Transfer(msg.sender, _to, _value)
	err = stub.SetEvent("Transfer", transferData)
	return err == nil
}

func setBalance(stub shim.ChaincodeStubInterface, cn string, balance uint64) error {
	key, err := stub.CreateCompositeKey(indexBalance, []string{cn})
	if err != nil {
		return ErrFailedKeyCreation
	}

	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, balance)
	if err := stub.PutState(key, data); err != nil {
		return ErrSettingState
	}

	return nil
}

func setAllowance(stub shim.ChaincodeStubInterface, from, spender string, value uint64) error {
	key, err := stub.CreateCompositeKey(indexAllowance, []string{from, spender})
	if err != nil {
		return ErrFailedKeyCreation
	}

	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, value)
	if err := stub.PutState(key, data); err != nil {
		return ErrSettingState
	}

	return nil
}
