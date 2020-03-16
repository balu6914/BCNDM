package http

import (
	"net/http"

	"datapace/streams"
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

type addStreamRes struct {
	ID string `json:"id"`
}

func (res addStreamRes) headers() map[string]string {
	return map[string]string{
		"location": res.ID,
	}
}

func (res addStreamRes) code() int {
	return http.StatusCreated
}

func (res addStreamRes) empty() bool {
	return false
}

type addBulkStreamsRes struct{}

func (res addBulkStreamsRes) headers() map[string]string {
	return map[string]string{}
}

func (res addBulkStreamsRes) code() int {
	return http.StatusCreated
}

func (res addBulkStreamsRes) empty() bool {
	return true
}

type conflictBulkStreamsRes struct {
	streams.ErrBulkConflict
}

func (res conflictBulkStreamsRes) headers() map[string]string {
	return map[string]string{}
}

func (res conflictBulkStreamsRes) code() int {
	return http.StatusConflict
}

func (res conflictBulkStreamsRes) empty() bool {
	return false
}

type searchStreamsRes struct {
	streams.Page
}

func (res searchStreamsRes) headers() map[string]string {
	return map[string]string{}
}

func (res searchStreamsRes) code() int {
	return http.StatusOK
}

func (res searchStreamsRes) empty() bool {
	return false
}

type updateStreamRes struct{}

func (res updateStreamRes) headers() map[string]string {
	return map[string]string{}
}

func (res updateStreamRes) code() int {
	return http.StatusOK
}

func (res updateStreamRes) empty() bool {
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

type errorRes struct {
	Err string `json:"error"`
}
