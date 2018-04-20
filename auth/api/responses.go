package api

import (
	"net/http"

	"monetasa/auth"
)

const contentType = "application/json; charset=utf-8"

type apiRes interface {
	code() int
	headers() map[string]string
	empty() bool
}

type identityRes struct {
	Email string `json:"email,omitempty"`
}

func (res identityRes) code() int {
	return http.StatusOK
}

func (res identityRes) headers() map[string]string {
	return map[string]string{}
}

func (res identityRes) empty() bool {
	return res.Email == ""
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

type removeRes struct{}

func (res removeRes) code() int {
	return http.StatusNoContent
}

func (res removeRes) headers() map[string]string {
	return map[string]string{}
}

func (res removeRes) empty() bool {
	return true
}

type userRes struct {
	created bool
}

func (res userRes) code() int {
	if res.created {
		return http.StatusCreated
	}

	return http.StatusOK
}

func (res userRes) headers() map[string]string {
	return map[string]string{}
}

func (res userRes) empty() bool {
	return true
}

type viewRes struct {
	auth.User
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
