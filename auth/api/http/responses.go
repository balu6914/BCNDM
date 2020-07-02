package http

import (
	"net/http"
)

const contentType = "application/json"

type apiRes interface {
	code() int
	headers() map[string]string
	empty() bool
}

type createRes struct {
	ID string
}

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
	Company      string `json:"company,omitempty"`
	Address      string `json:"address,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Disabled     bool   `json:"disabled,omitempty"`
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
