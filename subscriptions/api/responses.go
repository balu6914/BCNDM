package api

import (
	"monetasa/subscriptions"
	"net/http"
)

const contentType = "application/json"

type apiRes interface {
	code() int
	headers() map[string]string
	empty() bool
}

type subscriptionRes struct{}

func (res subscriptionRes) headers() map[string]string {
	return map[string]string{}
}

func (res subscriptionRes) code() int {
	return http.StatusCreated
}

func (res subscriptionRes) empty() bool {
	return false
}

type listSubsRes struct {
	Subscriptions []subscriptions.Subscription
}

func (res listSubsRes) headers() map[string]string {
	return map[string]string{}
}

func (res listSubsRes) code() int {
	return http.StatusOK
}

func (res listSubsRes) empty() bool {
	return false
}
