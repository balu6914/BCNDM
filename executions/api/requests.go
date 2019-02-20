package api

import "datapace/executions"

type apiReq interface {
	validate() error
}

type executionReq struct {
	Data string             `json:"data"`
	Mode executions.JobMode `json:"mode"`
}

func (req executionReq) validate() error {
	if req.Data == "" {
		return executions.ErrMalformedData
	}

	if req.Mode != executions.Centralized &&
		req.Mode != executions.Distributed &&
		req.Mode != executions.Federated {
		return executions.ErrMalformedData
	}

	return nil
}

type startExecutionReq struct {
	owner      string
	Algo       string         `json:"algo"`
	Executions []executionReq `json:"executions"`
}

func (req startExecutionReq) validate() error {
	if req.Algo == "" {
		return executions.ErrMalformedData
	}

	for _, exec := range req.Executions {
		if err := exec.validate(); err != nil {
			return err
		}
	}

	return nil
}

type viewReq struct {
	owner string
	id    string
}

func (req viewReq) validate() error {
	if req.id == "" {
		return executions.ErrMalformedData
	}

	return nil
}

type listReq struct {
	owner string
}

func (req listReq) validate() error {
	return nil
}
