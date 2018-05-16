package api

import (
	"net/http"

	"monetasa/dapp"
)

const contentType = "application/json"

type apiRes interface {
	code() int
	headers() map[string]string
	empty() bool
}

type createStreamRes struct {
	ID string
}

func (res createStreamRes) headers() map[string]string {
	return map[string]string{
		"location": res.ID,
	}
}

func (res createStreamRes) code() int {
	return http.StatusCreated
}

func (res createStreamRes) empty() bool {
	return false
}

type modifyStreamRes struct{}

func (res modifyStreamRes) headers() map[string]string {
	return map[string]string{}
}

func (res modifyStreamRes) code() int {
	return http.StatusOK
}

func (res modifyStreamRes) empty() bool {
	return true
}

type readStreamRes struct {
	Stream dapp.Stream
}

func (res readStreamRes) headers() map[string]string {
	return map[string]string{}
}

func (res readStreamRes) code() int {
	return http.StatusOK
}

func (res readStreamRes) empty() bool {
	return false
}

type searchStreamRes struct {
	Streams []dapp.Stream
}

func (res searchStreamRes) headers() map[string]string {
	return map[string]string{}
}

func (res searchStreamRes) code() int {
	return http.StatusOK
}

func (res searchStreamRes) empty() bool {
	return false
}
