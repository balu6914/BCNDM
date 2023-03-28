package main

type balanceRes struct {
	Value uint64 `json:"value"`
}

type allowanceRes struct {
	Value uint64 `json:"value"`
}

type txHistoryRes struct {
	TInfo  TokenInfo      `json:"tokenInfo"`
	TxList []TransferFrom `json:"txList"`
}
