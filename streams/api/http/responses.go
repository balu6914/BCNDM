package http

import (
	"net/http"

	"monetasa/streams"
)

const (
	contentType     = "application/json"
	fileContentType = "multipart/form-data"
)

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

type createBulkStreamRes struct{}

func (res createBulkStreamRes) headers() map[string]string {
	return map[string]string{}
}

func (res createBulkStreamRes) code() int {
	return http.StatusCreated
}

func (res createBulkStreamRes) empty() bool {
	return true
}

type editStreamRes struct{}

func (res editStreamRes) headers() map[string]string {
	return map[string]string{}
}

func (res editStreamRes) code() int {
	return http.StatusOK
}

func (res editStreamRes) empty() bool {
	return true
}

type viewStreamRes struct {
	streams.Stream
}

func (res viewStreamRes) headers() map[string]string {
	return map[string]string{}
}

func (res viewStreamRes) code() int {
	return http.StatusOK
}

func (res viewStreamRes) empty() bool {
	return false
}

type searchStreamRes struct {
	streams.Page
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

type removeStreamRes struct{}

func (res removeStreamRes) code() int {
	return http.StatusNoContent
}

func (res removeStreamRes) headers() map[string]string {
	return map[string]string{}
}

func (res removeStreamRes) empty() bool {
	return true
}
