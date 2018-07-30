package main

type balanceReq struct {
	Owner string `json:"owner"`
}

type transferReq struct {
	To    string `json:"to"`
	Value uint64 `json:"value"`
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
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint64 `json:"value"`
}
