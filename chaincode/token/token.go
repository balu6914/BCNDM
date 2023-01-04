package main

import (
	"encoding/binary"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
)

const (
	keyToken       = "token"
	indexBalance   = "cn~balance"
	indexAllowance = "cn~allowance"
	txKeyIndex     = "cn~txId"
	txHistoryIndex = "tx~cn~txId"
)

var _ Service = (*tokenChaincode)(nil)

type tokenChaincode struct{}

// NewService returns new ERC20 implementation instance.
func NewService() Service {
	return tokenChaincode{}
}

func (tc tokenChaincode) Init(stub shim.ChaincodeStubInterface, ti TokenInfo) error {
	// get caller CN from his certificate
	caller, err := callerCN(stub)
	if err != nil {
		return ErrUnauthorized
	}

	ti.ContractOwner = caller
	data, err := json.Marshal(ti)
	if err != nil {
		return ErrInvalidArgument
	}

	if err := stub.PutState(keyToken, data); err != nil {
		return ErrSettingState
	}

	// set the balance using a helper function
	if err := setBalance(stub, caller, 0); err != nil {
		return ErrFailedBalanceSet
	}

	return nil
}

func (tc tokenChaincode) TotalSupply(stub shim.ChaincodeStubInterface) (uint64, error) {
	ti, err := getTokenInfo(stub)
	if err != nil {
		return 0, err
	}

	balance, _, err := tc.BalanceOf(stub, ti.ContractOwner)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (tc tokenChaincode) BalanceOf(stub shim.ChaincodeStubInterface, owner string) (uint64, []string, error) {
	var balance uint64
	txDeltas := []string{}
	balance = 0
	key, err := stub.CreateCompositeKey(indexBalance, []string{owner})
	if err != nil {
		return balance, txDeltas, ErrFailedKeyCreation
	}

	// if the user cn is not in the state, then the balance is 0
	data, err := stub.GetState(key)
	if err != nil {
		return balance, txDeltas, ErrGettingState
	}

	if data == nil {
		balance = 0
	} else {
		balance += binary.LittleEndian.Uint64(data)
	}

	// check if from (address) is same to contract owner i.e. treasury account.
	// if yes, no need to check treasury account balance as there is unlimited supply.
	ti, err := getTokenInfo(stub)
	if err != nil {
		return balance, txDeltas, err
	}

	// check delats too
	resultsIterator, err := stub.GetStateByPartialCompositeKey(txKeyIndex, []string{owner})
	if err != nil {
		return balance, txDeltas, ErrFailedRichQuery
	}
	defer resultsIterator.Close()

	var responseRange *queryresult.KV
	// Check Delta's
	for i := 0; resultsIterator.HasNext(); i++ {
		if responseRange, err = resultsIterator.Next(); err != nil {
			return balance, txDeltas, err
		}
		txDeltas = append(txDeltas, responseRange.Key)
		jsonBytesJSON := responseRange.Value
		tokenDelta := new(TokenDelta)
		if err = json.Unmarshal([]byte(jsonBytesJSON), tokenDelta); err != nil {
			return balance, txDeltas, ErrGettingState
		}

		if tokenDelta.Operation == "PLUS" {
			if owner == ti.ContractOwner {
				balance -= tokenDelta.Value
			} else {
				balance += tokenDelta.Value
			}
		} else {
			if owner == ti.ContractOwner {
				balance += tokenDelta.Value
			} else {
				balance -= tokenDelta.Value
			}
		}
	}

	return balance, txDeltas, nil
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

	a := Approve{
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

	// check if from (address) is same to contract owner i.e. treasury account.
	// if yes, no need to check treasury account balance as there is unlimited supply.
	ti, err := getTokenInfo(stub)
	if err != nil {
		return err
	}

	var totalAmountToTransfer uint64
	totalAmountToTransfer = 0
	for _, tr := range transfers {
		if from == tr.To {
			continue
		}
		totalAmountToTransfer += tr.Value
	}

	// withdraw tokens from fromAccount
	if from != ti.ContractOwner {
		// if from is not treasury account, need to check from account balance first.
		fromBalance, txDeltas, err := tc.BalanceOf(stub, from)
		if err != nil {
			return err
		}

		if totalAmountToTransfer > fromBalance {
			return ErrNotEnoughTokens
		}

		// balanceOf[_from] -= totalAmountToTransfer
		if err := setBalance(stub, from, fromBalance-totalAmountToTransfer); err != nil {
			return err
		}

		// delete all deltas
		if err := deleteAllDeltas(stub, txDeltas); err != nil {
			return err
		}
	} else {
		// transfer is from treasury account. Just put minus deltas
		if err := newDelta(stub, from, "MINUS", totalAmountToTransfer); err != nil {
			return err
		}
	}

	events := []TransferFrom{}

	for _, tr := range transfers {
		if from == tr.To {
			continue
		}

		if err := newDelta(stub, tr.To, "PLUS", tr.Value); err != nil {
			return err
		}

		t := TransferFrom{
			From:  from,
			To:    tr.To,
			Value: tr.Value,
		}
		transferData, err := json.Marshal(t)
		if err != nil {
			return ErrFailedSerialization
		}

		if err := recordTxForHistory(stub, from, tr.To, transferData); err != nil {
			return err
		}

		events = append(events, t)
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

	// check if from (address) is same to contract owner i.e. treasury account.
	// if yes, no need to check treasury account balance as there is unlimited supply.
	ti, err := getTokenInfo(stub)
	if err != nil {
		return false
	}

	if from == ti.ContractOwner {
		// transaction is initiated from treasury account
		if err := newDelta(stub, from, "MINUS", value); err != nil {
			return false
		}

		if err := newDelta(stub, to, "PLUS", value); err != nil {
			return false
		}
	} else {
		// transaction is initiated from a user's account
		// retrieving balances
		fromBalance, txDeltas, err := tc.BalanceOf(stub, from)
		if err != nil {
			return false
		}

		// if (balanceOf[_from] < _value) throw
		// insufficient balance
		if fromBalance < value {
			return false
		}

		// balanceOf[_from] -= _value
		if err := setBalance(stub, from, fromBalance-value); err != nil {
			return false
		}

		// delete all deltas
		if err := deleteAllDeltas(stub, txDeltas); err != nil {
			return false
		}

		// write new delta for to account
		if err := newDelta(stub, to, "PLUS", value); err != nil {
			return false
		}
	}

	t := TransferFrom{
		From:  from,
		To:    to,
		Value: value,
	}
	transferData, err := json.Marshal(t)
	if err != nil {
		return false
	}

	if err := recordTxForHistory(stub, from, to, transferData); err != nil {
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

func newDelta(stub shim.ChaincodeStubInterface, owner, operation string, value uint64) error {
	txId := stub.GetTxID()
	key, err := getCompositeKey(stub, txKeyIndex, owner, txId)
	if err != nil {
		return ErrFailedKeyCreation
	}

	delta := new(TokenDelta)
	delta.Value = value
	delta.Operation = operation

	delatAsBytes, _ := json.Marshal(delta)
	if err := stub.PutState(key, delatAsBytes); err != nil {
		return ErrSettingState
	}
	return nil
}

func getCompositeKey(stub shim.ChaincodeStubInterface, compositeKeyIndex string, args ...string) (string, error) {
	key, err := stub.CreateCompositeKey(compositeKeyIndex, args)
	if err != nil {
		return "", ErrFailedKeyCreation
	}
	return key, nil
}

func deleteAllDeltas(stub shim.ChaincodeStubInterface, txDeltas []string) error {
	for _, txDelta := range txDeltas {
		if err := stub.DelState(txDelta); err != nil {
			return ErrDeletingState
		}
	}
	return nil
}

func recordTxForHistory(stub shim.ChaincodeStubInterface, from, to string, txBytes []byte) error {
	txId := stub.GetTxID()
	fromKey, err := getCompositeKey(stub, txHistoryIndex, from, txId)
	if err != nil {
		return ErrFailedKeyCreation
	}

	if err := stub.PutState(fromKey, txBytes); err != nil {
		return ErrSettingState
	}

	toKey, err := getCompositeKey(stub, txHistoryIndex, to, txId)
	if err != nil {
		return ErrFailedKeyCreation
	}

	if err := stub.PutState(toKey, txBytes); err != nil {
		return ErrSettingState
	}

	return nil
}

func getTokenInfo(stub shim.ChaincodeStubInterface) (*TokenInfo, error) {
	ti := new(TokenInfo)
	data, err := stub.GetState(keyToken)
	if err != nil {
		return ti, ErrGettingState
	}

	if err := json.Unmarshal(data, ti); err != nil {
		return ti, ErrFailedSerialization
	}

	return ti, nil
}

func (tc tokenChaincode) TxHistory(stub shim.ChaincodeStubInterface) ([]TransferFrom, error) {
	txList := []TransferFrom{}

	// get caller CN from his certificate
	caller, err := callerCN(stub)
	if err != nil {
		return txList, ErrUnauthorized
	}

	resultsIterator, err := stub.GetStateByPartialCompositeKey(txHistoryIndex, []string{caller})
	if err != nil {
		return txList, ErrFailedRichQuery
	}
	defer resultsIterator.Close()

	var responseRange *queryresult.KV
	for i := 0; resultsIterator.HasNext(); i++ {
		if responseRange, err = resultsIterator.Next(); err != nil {
			return txList, err
		}

		jsonBytesJSON := responseRange.Value
		tx := new(TransferFrom)
		if err = json.Unmarshal([]byte(jsonBytesJSON), tx); err != nil {
			return txList, ErrFailedSerialization
		}

		txList = append(txList, *tx)
	}

	return txList, nil
}

func (tc tokenChaincode) CollectDeltasForTreasury(stub shim.ChaincodeStubInterface) error {
	// get caller CN from his certificate
	caller, err := callerCN(stub)
	if err != nil {
		return ErrUnauthorized
	}

	ti, err := getTokenInfo(stub)
	if err != nil {
		return ErrGettingState
	}

	if ti.ContractOwner != caller {
		return ErrUnauthorized
	}

	balance, txDeltas, err := tc.BalanceOf(stub, ti.ContractOwner)
	if err != nil {
		return err
	}

	// set the balance using a helper function
	if err := setBalance(stub, ti.ContractOwner, balance); err != nil {
		return ErrFailedBalanceSet
	}

	// delete all deltas
	if err := deleteAllDeltas(stub, txDeltas); err != nil {
		return err
	}

	return nil
}
