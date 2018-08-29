package main

type feeReq struct {
	Owner string `json:"owner"`
	Value uint64 `json:"value"`
}

type transferReq struct {
	To    string `json:"to"`
	Value uint64 `json:"value"`
}
