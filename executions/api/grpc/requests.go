package grpc

import (
	"datapace/executions"
)

type algoReq struct {
	id       string
	name     string
	metadata map[string]string
}

func (req algoReq) validate() error {
	if req.id == "" || req.name == "" || len(req.metadata) == 0 {
		return executions.ErrMalformedData
	}

	return nil
}

type dataReq struct {
	id       string
	metadata map[string]string
}

func (req dataReq) validate() error {
	if req.id == "" || len(req.metadata) == 0 {
		return executions.ErrMalformedData
	}

	return nil
}
