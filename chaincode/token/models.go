package main

type tokenInfo struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Decimals    uint8  `json:"decimals"`
	TotalSupply uint64 `json:"totalSupply"`
}

type balance struct {
	User  string `json:"user"`
	Value uint64 `json:"value"`
}

type transfer struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint64 `json:"value"`
}

type approve struct {
	Spender string `json:"spender"`
	Value   uint64 `json:"value"`
}
