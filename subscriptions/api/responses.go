package api

import (
	"datapace/subscriptions"
	"net/http"
)

const contentType = "application/json"

type apiRes interface {
	code() int
	headers() map[string]string
	empty() bool
}

type addSubRes struct {
	ID string `json:"id"`
}

func (res addSubRes) headers() map[string]string {
	return map[string]string{}
}

func (res addSubRes) code() int {
	return http.StatusCreated
}

func (res addSubRes) empty() bool {
	return false
}

type viewSubRes struct {
	subscriptions.Subscription
}

func (res viewSubRes) headers() map[string]string {
	return map[string]string{}
}

func (res viewSubRes) code() int {
	return http.StatusOK
}

func (res viewSubRes) empty() bool {
	return false
}

type searchSubsRes struct {
	subscriptions.Page
}

func (res searchSubsRes) headers() map[string]string {
	return map[string]string{}
}

func (res searchSubsRes) code() int {
	return http.StatusOK
}

func (res searchSubsRes) empty() bool {
	return false
}
