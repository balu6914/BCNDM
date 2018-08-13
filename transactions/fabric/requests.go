package fabric

type balanceReq struct {
	Owner string `json:"owner"`
}

type transferReq struct {
	To    string `json:"to"`
	Value uint64 `json:"value"`
}
