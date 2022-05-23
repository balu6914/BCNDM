package grpc

import (
	"context"
	"github.com/datapace/sharing"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/go-kit/kit/endpoint"
)

func updateReceiversEndpoint(svc sharing.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(sharingPayload)
		if err := req.validate(); err != nil {
			return nil, err
		}
		var rcvUserIds []sharing.UserId
		for _, u := range req.receivers.userIds {
			rcvUserIds = append(rcvUserIds, sharing.UserId(u))
		}
		var rcvGroupIds []sharing.GroupId
		for _, g := range req.receivers.groupIds {
			rcvGroupIds = append(rcvGroupIds, sharing.GroupId(g))
		}

		s := sharing.Sharing{
			SourceUserId: sharing.UserId(req.sourceUserId),
			StreamId:     sharing.StreamId(req.streamId),
			Receivers: sharing.Receivers{
				Version:  req.receivers.versionRef,
				GroupIds: rcvGroupIds,
				UserIds:  rcvUserIds,
			},
		}
		if err := svc.UpdateReceivers(s); err != nil {
			return nil, err
		}
		return empty.Empty{}, nil
	}
}

func getSharingsToGroupsEndpoint(svc sharing.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getSharingsToGroupsRequest)
		if err := req.validate(); err != nil {
			return nil, err
		}
		var rcvGroupIds []sharing.GroupId
		for _, g := range req.receiverGroupIds {
			rcvGroupIds = append(rcvGroupIds, sharing.GroupId(g))
		}
		result, err := svc.GetSharingsToGroups(rcvGroupIds)
		if err != nil {
			return nil, err
		}
		resp := convertToGetSharingsResponse(result)
		return resp, nil
	}
}

func getSharingsEndpoint(svc sharing.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getSharingsRequest)
		if err := req.validate(); err != nil {
			return nil, err
		}
		var rcvGroupIds []sharing.GroupId
		for _, g := range req.receiverGroupIds {
			rcvGroupIds = append(rcvGroupIds, sharing.GroupId(g))
		}
		result, err := svc.GetSharings(sharing.UserId(req.receiverUserId), rcvGroupIds)
		if err != nil {
			return nil, err
		}
		resp := convertToGetSharingsResponse(result)
		return resp, nil
	}
}

func convertToGetSharingsResponse(sharings []sharing.Sharing) getSharingsResponse {
	var sharingPayloads []sharingPayload
	for _, s := range sharings {
		var rcvrGroupIds []string
		for _, g := range s.Receivers.GroupIds {
			rcvrGroupIds = append(rcvrGroupIds, string(g))
		}
		var rcvrUserIds []string
		for _, u := range s.Receivers.UserIds {
			rcvrUserIds = append(rcvrUserIds, string(u))
		}
		sp := sharingPayload{
			sourceUserId: string(s.SourceUserId),
			streamId:     string(s.StreamId),
			receivers: receivers{
				versionRef: s.Receivers.Version,
				groupIds:   rcvrGroupIds,
				userIds:    rcvrUserIds,
			},
		}
		sharingPayloads = append(sharingPayloads, sp)
	}
	return getSharingsResponse{sharings: sharingPayloads}
}
