package grpc

import (
	"context"
	"github.com/datapace/groups"
	"github.com/datapace/groups/proto"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ proto.GroupsServiceServer = (*grpcServer)(nil)

type grpcServer struct {
	proto.UnimplementedGroupsServiceServer
	getUserGroups kitgrpc.Handler
}

// NewServer instantiates new groups gRPC server.
func NewServer(svc groups.Service) proto.GroupsServiceServer {
	getUserGroups := kitgrpc.NewServer(
		getUserGroupsEndpoint(svc),
		decodeGetUserGroupsRequest,
		encodeGetUserGroupsResponse,
	)
	return &grpcServer{
		getUserGroups: getUserGroups,
	}
}

func (s *grpcServer) GetUserGroups(ctx context.Context, r *proto.GetUserGroupsRequest) (*proto.GetUserGroupsResponse, error) {
	_, resp, err := s.getUserGroups.ServeGRPC(ctx, r)
	if err != nil {
		return nil, encodeError(err)
	}
	return resp.(*proto.GetUserGroupsResponse), nil
}

func decodeGetUserGroupsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.GetUserGroupsRequest)
	return getUserGroupsRequest{req.GetUid()}, nil
}

func encodeGetUserGroupsResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	resp := grpcRes.(getUserGroupsResponse)
	return &proto.GetUserGroupsResponse{Gids: resp.gids}, encodeError(resp.err)
}

func encodeError(err error) error {
	switch err {
	case nil:
		return nil
	case groups.ErrBadRequest:
		return status.Error(codes.InvalidArgument, "received invalid request")
	case groups.ErrNotFound:
		return status.Error(codes.NotFound, "groups were not found")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
