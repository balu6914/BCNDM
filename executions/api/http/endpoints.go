package http

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
			e := executions.Execution{
				Owner:    req.owner,
				Algo:     req.Algo,
				Data:     exec.Data,
				Metadata: exec.Metadata,
			}
			id, err := svc.Start(e)
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
			ID:         exec.ID,
			ExternalID: exec.ExternalID,
			Algo:       exec.Algo,
			Data:       exec.Data,
			Metadata:   exec.Metadata,
			State:      exec.State,
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
				ID:         exec.ID,
				ExternalID: exec.ExternalID,
				Algo:       exec.Algo,
				Data:       exec.Data,
				Metadata:   exec.Metadata,
				State:      exec.State,
			})
		}

		return res, nil
	}
}

func resultEndpoint(svc executions.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(viewReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		result, err := svc.Result(req.owner, req.id)
		if err != nil {
			return nil, err
		}

		res := resultRes{
			Result: result,
		}

		return res, nil
	}
}
