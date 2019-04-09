package http

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
	ID                       string             `json:"id"`
	Algo                     string             `json:"algo"`
	Data                     string             `json:"data"`
	AdditionalLocalJobArgs   []string           `json:"local_args"`
	Type                     string             `json:"type"`
	GlobalTimeout            uint64             `json:"global_timeout"`
	LocalTimeout             uint64             `json:"local_timeout"`
	AdditionalPreprocessArgs []string           `json:"preprocess_args"`
	Mode                     executions.JobMode `json:"mode"`
	AdditionalGlobalJobArgs  []string           `json:"global_args"`
	AdditionalFiles          []string           `json:"files"`
	TrainedModelToken        string             `json:"model_token,omitempty"`
	TrainedModelName         string             `json:"model_name,omitempty"`
	State                    executions.State   `json:"state"`
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

type resultRes struct {
	Result map[string]interface{} `json:"result"`
}

func (res resultRes) code() int {
	return http.StatusOK
}

func (res resultRes) headers() map[string]string {
	return map[string]string{}
}

func (res resultRes) empty() bool {
	return false
}
