package api

import (
	"context"
	"datapace/executions"

	"github.com/go-kit/kit/endpoint"
)

func startExecutionEndpoint(svc executions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(startExecutionReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		res := startExecutionRes{
			Results: []startExecutionResult{},
		}
		for _, exec := range req.Executions {
			id, err := svc.Start(req.owner, req.Algo, exec.Data, exec.Mode)
			result := startExecutionResult{
				ID: id,
			}
			if err != nil {
				result.Error = err.Error()
			}

			res.Results = append(res.Results, result)
		}

		return res, nil
	}
}

func viewEndpoint(svc executions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(viewReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		exec, err := svc.Execution(req.owner, req.id)
		if err != nil {
			return nil, err
		}

		res := viewRes{
			ID:    exec.ID,
			State: exec.State,
			Algo:  exec.Algo,
			Data:  exec.Data,
			Mode:  exec.Mode,
		}

		return res, nil
	}
}

func listEndpoint(svc executions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(listReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		execs, err := svc.List(req.owner)
		if err != nil {
			return nil, err
		}

		res := listRes{
			Executions: []viewRes{},
		}
		for _, exec := range execs {
			res.Executions = append(res.Executions, viewRes{
				ID:    exec.ID,
				State: exec.State,
				Algo:  exec.Algo,
				Data:  exec.Data,
				Mode:  exec.Mode,
			})
		}

		return res, nil
	}
}
