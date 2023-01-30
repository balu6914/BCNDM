package main

import "time"

type contractReq struct {
	StreamID  string    `json:"stream_id,omitempty"`
	StartTime time.Time `json:"start_time,omitempty"`
	EndTime   time.Time `json:"end_time,omitempty"`
	OwnerID   string    `json:"owner_id,omitempty"`
	PartnerID string    `json:"partner_id,omitempty"`
	Share     uint64    `json:"share,omitempty"`
	Signed    bool      `json:"signed,omitempty"`
}

type signReq struct {
	StreamID string    `json:"stream_id,omitempty"`
	EndTime  time.Time `json:"end_time,omitempty"`
}

type transferReq struct {
	StreamID string    `json:"stream_id,omitempty"`
	Time     time.Time `json:"time,omitempty"`
	To       string    `json:"to"`
	Value    uint64    `json:"value"`
	DateTime string    `json:"dateTime"`
}
