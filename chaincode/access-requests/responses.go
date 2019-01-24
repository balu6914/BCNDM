package main

type accessRes struct {
	Requester string `json:"requester"`
	Receiver  string `json:"receiver"`
	State     State  `json:"state"`
}
