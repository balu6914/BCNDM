package grpc

import (
	"context"
	"github.com/datapace/groups"

	"github.com/go-kit/kit/endpoint"
)

func getUserGroupsEndpoint(svc groups.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getUserGroupsRequest)
		if err := req.validate(); err != nil {
			return nil, err
		}
		gids, err := svc.GetUserGroups(groups.Uid(req.uid))
		if err != nil {
			return getUserGroupsResponse{}, err
		}
		results := []string{}
		for _, gid := range gids {
			results = append(results, string(gid))
		}
		res := getUserGroupsResponse{
			gids: results,
		}
		return res, nil
	}
}
