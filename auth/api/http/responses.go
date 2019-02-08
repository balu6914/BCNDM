package http

import (
	"datapace/auth"
	"net/http"
)

const contentType = "application/json"

type apiRes interface {
	code() int
	headers() map[string]string
	empty() bool
}

type createRes struct{}

func (res createRes) headers() map[string]string {
	return map[string]string{}
}

func (res createRes) code() int {
	return http.StatusCreated
}

func (res createRes) empty() bool {
	return true
}

type tokenRes struct {
	Token string `json:"token,omitempty"`
}

func (res tokenRes) code() int {
	return http.StatusCreated
}

func (res tokenRes) headers() map[string]string {
	return map[string]string{}
}

func (res tokenRes) empty() bool {
	return res.Token == ""
}

type updateRes struct{}

type updatePasswordRes struct{}

func (res updateRes) code() int {
	return http.StatusOK
}

func (res updateRes) headers() map[string]string {
	return map[string]string{}
}

func (res updateRes) empty() bool {
	return true
}

type viewRes struct {
	ID           string `json:"id"`
	Email        string `json:"email,omitempty"`
	ContactEmail string `json:"contact_email,omitempty"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

func (res viewRes) code() int {
	return http.StatusOK
}

func (res viewRes) headers() map[string]string {
	return map[string]string{}
}

func (res viewRes) empty() bool {
	return false
}

type listRes struct {
	Users []viewRes `json:"users"`
}

func (res listRes) code() int {
	return http.StatusOK
}

func (res listRes) headers() map[string]string {
	return map[string]string{}
}

func (res listRes) empty() bool {
	return false
}

type requestAccessRes struct {
	id string
}

func (res requestAccessRes) code() int {
	return http.StatusCreated
}

func (res requestAccessRes) headers() map[string]string {
	return map[string]string{
		"Location": res.id,
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

type rejectAccessRes struct{}

func (res rejectAccessRes) code() int {
	return http.StatusOK
}

func (res rejectAccessRes) headers() map[string]string {
	return map[string]string{}
}

func (res rejectAccessRes) empty() bool {
	return true
}

type viewAccessRequestRes struct {
	ID       string     `json:"id"`
	Sender   string     `json:"sender"`
	Receiver string     `json:"receiver"`
	State    auth.State `json:"state"`
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
