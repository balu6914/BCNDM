package fabric

type balanceReq struct {
	Owner string `json:"owner"`
}

type transferFromReq struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint64 `json:"value"`
}

type transferReq struct {
	To    string `json:"to"`
	Value uint64 `json:"value"`
}
