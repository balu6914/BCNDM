package http

import (
	"net/http"
	"time"
)

type apiRes interface {
	code() int
	headers() map[string]string
	empty() bool
}

type balanceRes struct {
	Balance uint64 `json:"balance"`
}

func (res balanceRes) headers() map[string]string {
	return map[string]string{}
}

func (res balanceRes) code() int {
	return http.StatusOK
}

func (res balanceRes) empty() bool {
	return false
}

type buyRes struct{}

func (res buyRes) headers() map[string]string {
	return map[string]string{}
}

func (res buyRes) code() int {
	return http.StatusOK
}

func (res buyRes) empty() bool {
	return true
}

type withdrawRes struct{}

func (res withdrawRes) headers() map[string]string {
	return map[string]string{}
}

func (res withdrawRes) code() int {
	return http.StatusOK
}

func (res withdrawRes) empty() bool {
	return true
}

type createContractsRes struct{}

func (res createContractsRes) headers() map[string]string {
	return map[string]string{}
}

func (res createContractsRes) code() int {
	return http.StatusCreated
}

func (res createContractsRes) empty() bool {
	return true
}

type signContractRes struct{}

func (res signContractRes) headers() map[string]string {
	return map[string]string{}
}

func (res signContractRes) code() int {
	return http.StatusOK
}

func (res signContractRes) empty() bool {
	return true
}

type listContractsRes struct {
	Page      uint64         `json:"page"`
	Limit     uint64         `json:"limit"`
	Total     uint64         `json:"total"`
	Contracts []contractView `json:"contracts"`
}

type contractView struct {
	StreamID  string    `json:"stream_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	OwnerID   string    `json:"owner_id"`
	PartnerID string    `json:"partner_id"`
	Share     uint64    `json:"share"`
	Signed    bool      `json:"signed"`
}

func (res listContractsRes) headers() map[string]string {
	return map[string]string{}
}

func (res listContractsRes) code() int {
	return http.StatusOK
}

func (res listContractsRes) empty() bool {
	return false
}
