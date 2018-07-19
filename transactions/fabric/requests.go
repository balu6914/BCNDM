package fabric

type balanceReq struct {
	Owner string `json:"user"`
}

type transferFromReq struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint64 `json:"value"`
}
