package http

import (
	"net/http"
	"time"
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
	return map[string]string{
		"Location": res.ID,
	}
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

type okRes struct{}

func (res okRes) code() int {
	return http.StatusOK
}

func (res okRes) headers() map[string]string {
	return map[string]string{}
}

func (res okRes) empty() bool {
	return true
}

type viewUserRes struct {
	ID           string                 `json:"id"`
	Email        string                 `json:"email,omitempty"`
	ContactEmail string                 `json:"contact_email,omitempty"`
	FirstName    string                 `json:"first_name"`
	LastName     string                 `json:"last_name"`
	Company      string                 `json:"company,omitempty"`
	Address      string                 `json:"address,omitempty"`
	Country      string                 `json:"country,omitempty"`
	Mobile       string                 `json:"mobile,omitempty"`
	Phone        string                 `json:"phone,omitempty"`
	Disabled     bool                   `json:"disabled,omitempty"`
	Locked       bool                   `json:"locked,omitempty"`
	Role         string                 `json:"role,omitempty"`
	CreatedDate  *time.Time             `json:"created_date,omitempty"`
	Policies     []interface{}          `json:"policies,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

func (res viewUserRes) code() int {
	return http.StatusOK
}

func (res viewUserRes) headers() map[string]string {
	return map[string]string{}
}

func (res viewUserRes) empty() bool {
	return false
}

type listUsersRes struct {
	Users []viewUserRes `json:"users"`
}

func (res listUsersRes) code() int {
	return http.StatusOK
}

func (res listUsersRes) headers() map[string]string {
	return map[string]string{}
}

func (res listUsersRes) empty() bool {
	return false
}

type viewPolicyRes struct {
	ID      string `json:"id"`
	Version string `json:"version,omitempty"`
	Owner   string `json:"owner"`
	Name    string `json:"name,omitempty"`
	Rules   []rule `json:"rules,omitempty"`
}

func (res viewPolicyRes) code() int {
	return http.StatusOK
}

func (res viewPolicyRes) headers() map[string]string {
	return map[string]string{}
}

func (res viewPolicyRes) empty() bool {
	return false
}

type listPoliciesRes struct {
	Policies []viewPolicyRes `json:"policies"`
}

func (res listPoliciesRes) code() int {
	return http.StatusOK
}

func (res listPoliciesRes) headers() map[string]string {
	return map[string]string{}
}

func (res listPoliciesRes) empty() bool {
	return false
}

type removeRes struct{}

func (res removeRes) Code() int {
	return http.StatusNoContent
}

func (res removeRes) Headers() map[string]string {
	return map[string]string{}
}

func (res removeRes) Empty() bool {
	return true
}
