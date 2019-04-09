package grpc

import (
	"datapace/executions"
	"fmt"
)

type algoReq struct {
	id         string
	name       string
	path       string
	modelToken string
	modelName  string
}

func (req algoReq) validate() error {
	if req.id == "" || req.name == "" || req.path == "" {
		return executions.ErrMalformedData
	}

	return nil
}

type dataReq struct {
	id   string
	path string
}

func (req dataReq) validate() error {
	fmt.Println(req)
	if req.id == "" || req.path == "" {
		return executions.ErrMalformedData
	}

	return nil
}
