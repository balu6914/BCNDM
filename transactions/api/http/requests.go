package http

import (
	"time"

	"github.com/datapace/datapace/transactions"
)

type balanceReq struct {
	userID string
	ID     string
}

func (req balanceReq) validate() error {
	if req.userID == "" {
		return errMalformedEntity
	}

	return nil
}

type txHistoryReq struct {
	userID       string
	fromDateTime string
	toDateTime   string
	txType       string
}

func (req txHistoryReq) validate() error {
	if req.userID == "" {
		return errMalformedEntity
	}

	return nil
}

type buyReq struct {
	userID string
	FundID string `json:"fund_id"`
	Amount uint64 `json:"amount"`
}

func (req buyReq) validate() error {
	if req.Amount == 0 {
		return errMalformedEntity
	}

	return nil
}

type withdrawReq struct {
	userID string
	FundID string `json:"fund_id"`
	Amount uint64 `json:"amount"`
}

func (req withdrawReq) validate() error {
	if req.Amount == 0 {
		return errMalformedEntity
	}

	return nil
}

type createContractsReq struct {
	ownerID     string
	StreamID    string         `json:"stream_id"`
	EndTime     time.Time      `json:"end_time"`
	Description string         `json:"description"`
	Items       []contractItem `json:"items"`
}

func (req createContractsReq) validate() error {
	if req.ownerID == "" || req.StreamID == "" ||
		req.EndTime.Before(time.Now()) || len(req.Items) == 0 {
		return errMalformedEntity
	}

	if len(req.Description) > maxDescLength {
		return errInvalidDescription
	}

	for _, item := range req.Items {
		if req.ownerID == item.PartnerID {
			return errMalformedEntity
		}

		if err := item.validate(); err != nil {
			return err
		}
	}

	return nil
}

type contractItem struct {
	PartnerID string `json:"partner_id"`
	Share     uint64 `json:"share"`
}

func (req contractItem) validate() error {
	if req.PartnerID == "" {
		return errMalformedEntity
	}

	return nil
}

type signContractReq struct {
	partnerID string
	StreamID  string    `json:"stream_id"`
	EndTime   time.Time `json:"end_time"`
}

func (req signContractReq) validate() error {
	if req.EndTime.Before(time.Now()) || req.partnerID == "" || req.StreamID == "" {
		return errMalformedEntity
	}

	return nil
}

type listContractsReq struct {
	userID string
	page   uint64
	limit  uint64
	role   transactions.Role
}

func (req listContractsReq) validate() error {
	if req.userID == "" {
		return errMalformedEntity
	}

	return nil
}
