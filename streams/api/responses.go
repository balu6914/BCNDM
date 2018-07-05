package api

import (
	"net/http"

	"monetasa/streams"
)

const contentType = "application/json"

type apiRes interface {
	code() int
	headers() map[string]string
	empty() bool
}

type createStreamRes struct {
	ID string `json:"id"`
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

type createBulkStreamResponse struct{}

func (res createBulkStreamResponse) headers() map[string]string {
	return map[string]string{}
}

func (res createBulkStreamResponse) code() int {
	return http.StatusCreated
}

func (res createBulkStreamResponse) empty() bool {
	return true
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
	Stream streams.Stream
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
	Streams []streams.Stream
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
