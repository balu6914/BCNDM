package api

import (
	"datapace/executions"
	"net/http"
)

type apiRes interface {
	code() int
	headers() map[string]string
	empty() bool
}

type startExecutionResult struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

type startExecutionRes struct {
	Results []startExecutionResult `json:"results"`
}

func (res startExecutionRes) code() int {
	return http.StatusCreated
}

func (res startExecutionRes) headers() map[string]string {
	return map[string]string{}
}

func (res startExecutionRes) empty() bool {
	return false
}

type viewRes struct {
	ID    string             `json:"id"`
	State executions.State   `json:"state"`
	Algo  string             `json:"algo"`
	Data  string             `json:"data"`
	Mode  executions.JobMode `json:"mode"`
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
	Executions []viewRes `json:"executions"`
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
