package http

import "net/http"

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
