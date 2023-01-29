package main

type tokenInfo struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals uint8  `json:"decimals"`
}

type balanceReq struct {
	Owner string `json:"owner"`
}

type txHistoryReq struct {
	Owner string `json:"owner"`
}

type transferReq struct {
	To       string `json:"to"`
	Value    uint64 `json:"value"`
	DateTime string `json:"dateTime"` // dateTime should be added at middleware level in format: DD-MM-YYYY hh:mm:ss
	TxType   string `json:"TxType"`
}

type approveReq struct {
	Spender string `json:"spender"`
	Value   uint64 `json:"value"`
}

type allowanceReq struct {
	Owner   string `json:"owner"`
	Spender string `json:"spender"`
}

type transferFromReq struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Value    uint64 `json:"value"`
	DateTime string `json:"dateTime"` // dateTime should be added at middleware level in format: DD-MM-YYYY hh:mm:ss
}
