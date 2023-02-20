package fabric

import "time"

type balanceReq struct {
	Owner string `json:"owner"`
}

type txHistoryReq struct {
	Owner        string `json:"owner"`
	FromDateTime string `json:"fromDateTime"`
	ToDateTime   string `json:"toDateTime"`
	TxType       string `json:"txType"`
}

type transferReq struct {
	StreamID string    `json:"stream_id"`
	To       string    `json:"to"`
	Time     time.Time `json:"time"`
	Value    uint64    `json:"value"`
	DateTime string    `json:"dateTime"`
}

type createContractReq struct {
	StreamID  string    `json:"stream_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	OwnerID   string    `json:"owner_id"`
	PartnerID string    `json:"partner_id"`
	Share     uint64    `json:"share"`
}

type signContractReq struct {
	StreamID string    `json:"stream_id"`
	EndTime  time.Time `json:"end_time"`
}
