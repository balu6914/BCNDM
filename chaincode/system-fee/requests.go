package main

type feeReq struct {
	Value uint64 `json:"value"`
}

type transferReq struct {
	To    string `json:"to"`
	Value uint64 `json:"value"`
}
