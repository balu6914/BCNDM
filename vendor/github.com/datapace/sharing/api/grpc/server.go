package grpc

import (
	"context"
	"github.com/datapace/sharing"
	"github.com/datapace/sharing/proto"
	"github.com/golang/protobuf/ptypes/empty"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcServer struct {
	proto.UnimplementedSharingServiceServer
	updateReceivers     kitgrpc.Handler
	getSharingsToGroups kitgrpc.Handler
	getSharings         kitgrpc.Handler
}

// NewServer instantiates new groups gRPC server.
func NewServer(svc sharing.Service) proto.SharingServiceServer {
	updateReceivers := kitgrpc.NewServer(
		updateReceiversEndpoint(svc),
		decodeUpdateReceiversRequest,
		encodeUpdateReceiversResponse,
	)
	getSharingsToGroups := kitgrpc.NewServer(
		getSharingsToGroupsEndpoint(svc),
		decodeGetSharingsToGroupsRequest,
		encodeGetSharingsResponse,
	)
	getSharings := kitgrpc.NewServer(
		getSharingsEndpoint(svc),
		decodeGetSharingsRequest,
		encodeGetSharingsResponse,
	)
	return &grpcServer{
		updateReceivers:     updateReceivers,
		getSharingsToGroups: getSharingsToGroups,
		getSharings:         getSharings,
	}
}

func (srv *grpcServer) UpdateReceivers(ctx context.Context, req *proto.Sharing) (*empty.Empty, error) {
	_, _, err := srv.updateReceivers.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}
	return &empty.Empty{}, nil
}

func (srv *grpcServer) GetSharingsToGroups(ctx context.Context, req *proto.GetSharingsToGroupsRequest) (*proto.GetSharingsResponse, error) {
	_, resp, err := srv.getSharingsToGroups.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}
	return resp.(*proto.GetSharingsResponse), nil
}

func (srv *grpcServer) GetSharings(ctx context.Context, req *proto.GetSharingsRequest) (*proto.GetSharingsResponse, error) {
	_, resp, err := srv.getSharings.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}
	return resp.(*proto.GetSharingsResponse), nil
}

func decodeUpdateReceiversRequest(ctx context.Context, grpcReq interface{}) (request interface{}, err error) {
	req := grpcReq.(*proto.Sharing)
	var versionRef *uint64
	switch versionOpt := req.Receivers.VersionOption.(type) {
	case *proto.Receivers_VersionEmpty:
		versionRef = nil
	case *proto.Receivers_Version:
		versionRef = &versionOpt.Version
	}
	sp := sharingPayload{
		sourceUserId: req.SourceUserId,
		streamId:     req.StreamId,
		receivers: receivers{
			versionRef: versionRef,
			groupIds:   req.Receivers.GroupIds,
			userIds:    req.Receivers.UserIds,
		},
	}
	return sp, nil
}

func encodeUpdateReceiversResponse(ctx context.Context, grpcResp interface{}) (response interface{}, err error) {
	return &empty.Empty{}, nil
}

func decodeGetSharingsToGroupsRequest(ctx context.Context, grpcReq interface{}) (request interface{}, err error) {
	req := grpcReq.(*proto.GetSharingsToGroupsRequest)
	return getSharingsToGroupsRequest{receiverGroupIds: req.ReceiverGroupIds}, nil
}

func encodeGetSharingsResponse(ctx context.Context, grpcResp interface{}) (response interface{}, err error) {
	resp := grpcResp.(getSharingsResponse)
	var sharings []*proto.Sharing
	for _, sp := range resp.sharings {
		var rs proto.Receivers
		versionRef := sp.receivers.versionRef
		switch versionRef {
		case nil:
			rs = proto.Receivers{
				VersionOption: emptyVersion,
				GroupIds:      sp.receivers.groupIds,
				UserIds:       sp.receivers.userIds,
			}
		default:
			rs = proto.Receivers{
				VersionOption: &proto.Receivers_Version{Version: *versionRef},
				GroupIds:      sp.receivers.groupIds,
				UserIds:       sp.receivers.userIds,
			}
		}
		s := proto.Sharing{
			SourceUserId: sp.sourceUserId,
			StreamId:     sp.streamId,
			Receivers:    &rs,
		}
		sharings = append(sharings, &s)
	}
	return &proto.GetSharingsResponse{Sharings: sharings}, nil
}

func decodeGetSharingsRequest(ctx context.Context, grpcReq interface{}) (request interface{}, err error) {
	req := grpcReq.(*proto.GetSharingsRequest)
	return getSharingsRequest{
		getSharingsToGroupsRequest: getSharingsToGroupsRequest{receiverGroupIds: req.ReceiverGroupIds},
		receiverUserId:             req.ReceiverUserId,
	}, nil
}

func encodeError(err error) error {
	switch err {
	case nil:
		return nil
	case sharing.ErrBadRequest:
		return status.Error(codes.InvalidArgument, "received invalid request")
	case sharing.ErrNotFound:
		return status.Error(codes.NotFound, "entity not found or version mismatch")
	case sharing.ErrConflict:
		return status.Error(codes.AlreadyExists, "entity already exists")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
