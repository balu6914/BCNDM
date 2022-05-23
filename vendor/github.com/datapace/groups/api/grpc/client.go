package grpc

import (
	"context"
	"github.com/datapace/groups/proto"

	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

var _ proto.GroupsServiceClient = (*grpcClient)(nil)

type grpcClient struct {
	getUserGroupsEndpoint endpoint.Endpoint
}

// NewClient returns new gRPC client instance.<
func NewClient(conn *grpc.ClientConn) proto.GroupsServiceClient {
	getUserGroupsEndpoint := kitgrpc.NewClient(
		conn,
		"groups.GroupsService",
		"GetUserGroups",
		encodeGetUserGroupsRequest,
		decodeGetUserGroupsResponse,
		proto.GetUserGroupsResponse{},
	).Endpoint()

	return &grpcClient{
		getUserGroupsEndpoint: getUserGroupsEndpoint,
	}
}

func (client grpcClient) GetUserGroups(ctx context.Context, r *proto.GetUserGroupsRequest, _ ...grpc.CallOption) (*proto.GetUserGroupsResponse, error) {
	resp, err := client.getUserGroupsEndpoint(ctx, getUserGroupsRequest{r.GetUid()})
	if err != nil {
		return nil, err
	}
	getUserGroupsResp := resp.(getUserGroupsResponse)
	return &proto.GetUserGroupsResponse{Gids: getUserGroupsResp.gids}, getUserGroupsResp.err
}

func encodeGetUserGroupsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(getUserGroupsRequest)
	return &proto.GetUserGroupsRequest{Uid: req.uid}, nil
}

func decodeGetUserGroupsResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	resp := grpcRes.(*proto.GetUserGroupsResponse)
	return getUserGroupsResponse{resp.Gids, nil}, nil
}
