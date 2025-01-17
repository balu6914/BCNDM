package api

import (
	"net/http"

	"github.com/datapace/datapace/subscriptions"
)

const contentType = "application/json"

type apiRes interface {
	code() int
	headers() map[string]string
	empty() bool
}
type addSubResp struct {
	StreamID       string `json:"streamID,omitempty"`
	SubscriptionID string `json:"subscriptionID,omitempty"`
	ErrorMessage   string `json:"errorMessage,omitempty"`
}

type addSubsRes struct {
	Responses []addSubResp
}

func (res addSubsRes) headers() map[string]string {
	return map[string]string{}
}

func (res addSubsRes) code() int {
	return http.StatusCreated
}

func (res addSubsRes) empty() bool {
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

type reportResponse []byte

func (res reportResponse) headers() map[string]string {
	return map[string]string{
		"Content-Type": "application/pdf",
	}
}

func (res reportResponse) code() int {
	return http.StatusOK
}

func (res reportResponse) empty() bool {
	return false
}
