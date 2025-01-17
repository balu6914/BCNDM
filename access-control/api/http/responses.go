package http

import (
	"fmt"
	"net/http"

	access "github.com/datapace/datapace/access-control"
)

const contentType = "application/json"

type apiRes interface {
	code() int
	headers() map[string]string
	empty() bool
}

type requestAccessRes struct {
	id string
}

func (res requestAccessRes) code() int {
	return http.StatusCreated
}

func (res requestAccessRes) headers() map[string]string {
	return map[string]string{
		"Location": fmt.Sprintf("/access-requests/%s", res.id),
	}
}

func (res requestAccessRes) empty() bool {
	return true
}

type approveAccessRes struct{}

func (res approveAccessRes) code() int {
	return http.StatusOK
}

func (res approveAccessRes) headers() map[string]string {
	return map[string]string{}
}

func (res approveAccessRes) empty() bool {
	return true
}

type revokeAccessRes struct{}

func (res revokeAccessRes) code() int {
	return http.StatusOK
}

func (res revokeAccessRes) headers() map[string]string {
	return map[string]string{}
}

func (res revokeAccessRes) empty() bool {
	return true
}

type viewAccessRequestRes struct {
	ID       string       `json:"id"`
	Sender   string       `json:"sender"`
	Receiver string       `json:"receiver"`
	State    access.State `json:"state"`
}

type listAccessRequestsRes struct {
	Requests []viewAccessRequestRes
}

func (res listAccessRequestsRes) code() int {
	return http.StatusOK
}

func (res listAccessRequestsRes) headers() map[string]string {
	return map[string]string{}
}

func (res listAccessRequestsRes) empty() bool {
	return false
}

type grantAccessResp struct {
	id string
}

func (res grantAccessResp) code() int {
	return http.StatusCreated
}

func (res grantAccessResp) headers() map[string]string {
	return map[string]string{
		"Location": fmt.Sprintf("/access-requests/%s", res.id),
	}
}

func (res grantAccessResp) empty() bool {
	return true
}
