package grpc

import (
	"context"
	"github.com/datapace/sharing/proto"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

var (
	emptyVersion = &proto.Receivers_VersionEmpty{VersionEmpty: true}
)

type grpcClient struct {
	updateReceivers     endpoint.Endpoint
	getSharingsToGroups endpoint.Endpoint
	getSharings         endpoint.Endpoint
}

// NewClient returns new gRPC client instance.<
func NewClient(conn *grpc.ClientConn) proto.SharingServiceClient {
	updateReceivers := kitgrpc.NewClient(
		conn,
		"sharing.SharingService",
		"UpdateReceivers",
		encodeUpdateReceiversRequest,
		decodeUpdateReceiversResponse,
		empty.Empty{},
	).Endpoint()
	getSharingsToGroups := kitgrpc.NewClient(
		conn,
		"sharing.SharingService",
		"GetSharingsToGroups",
		encodeGetSharingsToGroupsRequest,
		decodeGetSharingsResponse,
		proto.GetSharingsResponse{},
	).Endpoint()
	getSharings := kitgrpc.NewClient(
		conn,
		"sharing.SharingService",
		"GetSharings",
		encodeGetSharingsRequest,
		decodeGetSharingsResponse,
		proto.GetSharingsResponse{},
	).Endpoint()
	return &grpcClient{
		updateReceivers:     updateReceivers,
		getSharingsToGroups: getSharingsToGroups,
		getSharings:         getSharings,
	}
}

func (client grpcClient) UpdateReceivers(ctx context.Context, in *proto.Sharing, opts ...grpc.CallOption) (*empty.Empty, error) {
	var versionRef *uint64
	switch versionOpt := in.Receivers.VersionOption.(type) {
	case *proto.Receivers_VersionEmpty:
		versionRef = nil
	case *proto.Receivers_Version:
		versionRef = &versionOpt.Version
	}
	sp := sharingPayload{
		sourceUserId: in.SourceUserId,
		streamId:     in.StreamId,
		receivers: receivers{
			versionRef: versionRef,
			groupIds:   in.Receivers.GroupIds,
			userIds:    in.Receivers.UserIds,
		},
	}
	_, err := client.updateReceivers(ctx, sp)
	return &empty.Empty{}, err
}

func (client grpcClient) GetSharingsToGroups(ctx context.Context, in *proto.GetSharingsToGroupsRequest, opts ...grpc.CallOption) (*proto.GetSharingsResponse, error) {
	req := getSharingsToGroupsRequest{receiverGroupIds: in.ReceiverGroupIds}
	resp, err := client.getSharingsToGroups(ctx, req)
	if err != nil {
		return nil, err
	}
	getSharingsResp := resp.(getSharingsResponse)
	var sharings []*proto.Sharing
	for _, s := range getSharingsResp.sharings {
		rcvrs := proto.Receivers{
			VersionOption: &proto.Receivers_Version{Version: *s.receivers.versionRef},
			GroupIds:      s.receivers.groupIds,
			UserIds:       s.receivers.userIds,
		}
		sharing := proto.Sharing{
			SourceUserId: s.sourceUserId,
			StreamId:     s.streamId,
			Receivers:    &rcvrs,
		}
		sharings = append(sharings, &sharing)
	}
	return &proto.GetSharingsResponse{Sharings: sharings}, nil
}

func (client grpcClient) GetSharings(ctx context.Context, in *proto.GetSharingsRequest, opts ...grpc.CallOption) (*proto.GetSharingsResponse, error) {
	req := getSharingsRequest{
		getSharingsToGroupsRequest: getSharingsToGroupsRequest{receiverGroupIds: in.ReceiverGroupIds},
		receiverUserId:             in.ReceiverUserId,
	}
	resp, err := client.getSharings(ctx, req)
	if err != nil {
		return nil, err
	}
	getSharingsResp := resp.(getSharingsResponse)
	var sharings []*proto.Sharing
	for _, s := range getSharingsResp.sharings {
		rcvrs := proto.Receivers{
			VersionOption: &proto.Receivers_Version{Version: *s.receivers.versionRef},
			GroupIds:      s.receivers.groupIds,
			UserIds:       s.receivers.userIds,
		}
		sharing := proto.Sharing{
			SourceUserId: s.sourceUserId,
			StreamId:     s.streamId,
			Receivers:    &rcvrs,
		}
		sharings = append(sharings, &sharing)
	}
	return &proto.GetSharingsResponse{Sharings: sharings}, nil
}

func encodeUpdateReceiversRequest(ctx context.Context, grpcReq interface{}) (request interface{}, err error) {
	req := grpcReq.(sharingPayload)
	var rcvrs proto.Receivers
	versionRef := req.receivers.versionRef
	switch versionRef {
	case nil:
		rcvrs = proto.Receivers{
			VersionOption: emptyVersion,
			GroupIds:      req.receivers.groupIds,
			UserIds:       req.receivers.userIds,
		}
	default:
		rcvrs = proto.Receivers{
			VersionOption: &proto.Receivers_Version{Version: *versionRef},
			GroupIds:      req.receivers.groupIds,
			UserIds:       req.receivers.userIds,
		}
	}
	return &proto.Sharing{
		SourceUserId: req.sourceUserId,
		StreamId:     req.streamId,
		Receivers:    &rcvrs,
	}, nil
}

func decodeUpdateReceiversResponse(ctx context.Context, grpcResp interface{}) (response interface{}, err error) {
	return &updateReceiversResponse{err: nil}, nil
}

func encodeGetSharingsToGroupsRequest(ctx context.Context, grpcReq interface{}) (request interface{}, err error) {
	req := grpcReq.(getSharingsToGroupsRequest)
	return &proto.GetSharingsToGroupsRequest{ReceiverGroupIds: req.receiverGroupIds}, nil
}

func decodeGetSharingsResponse(ctx context.Context, grpcResp interface{}) (response interface{}, err error) {
	resp := grpcResp.(*proto.GetSharingsResponse)
	var sharings []sharingPayload
	for _, s := range resp.Sharings {
		var versionRef *uint64
		switch vestionOpt := s.Receivers.VersionOption.(type) {
		case *proto.Receivers_VersionEmpty:
			versionRef = nil
		case *proto.Receivers_Version:
			versionRef = &vestionOpt.Version
		}
		rs := receivers{
			versionRef: versionRef,
			groupIds:   s.Receivers.GroupIds,
			userIds:    s.Receivers.UserIds,
		}
		sp := sharingPayload{
			sourceUserId: s.SourceUserId,
			streamId:     s.StreamId,
			receivers:    rs,
		}
		sharings = append(sharings, sp)
	}
	return getSharingsResponse{sharings: sharings}, nil
}

func encodeGetSharingsRequest(ctx context.Context, grpcReq interface{}) (request interface{}, err error) {
	req := grpcReq.(getSharingsRequest)
	return &proto.GetSharingsRequest{
		ReceiverGroupIds: req.receiverGroupIds,
		ReceiverUserId:   req.receiverUserId,
	}, nil
}
