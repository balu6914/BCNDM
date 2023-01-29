package main

type feeReq struct {
	Owner string `json:"owner"`
	Value uint64 `json:"value"`
}

type transferReq struct {
	To       string `json:"to"`
	Value    uint64 `json:"value"`
	DateTime string `json:"dateTime"`
	TxType   string `json:"txType"`
}

type transferStatusReq struct {
	Owner     string        `json:"owner"`
	Value     uint64        `json:"value"`
	DateTime  string        `json:"dateTime"`
	Transfers []transferReq `json:"transfers"`
}
