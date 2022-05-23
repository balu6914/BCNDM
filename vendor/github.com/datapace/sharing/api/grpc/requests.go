package grpc

import "github.com/datapace/sharing"

type receivers struct {
	versionRef *uint64
	groupIds   []string
	userIds    []string
}

func (r receivers) validate() error {
	if r.groupIds == nil {
		return sharing.ErrBadRequest
	}
	if r.userIds == nil {
		return sharing.ErrBadRequest
	}
	return nil
}

type sharingPayload struct {
	sourceUserId string
	streamId     string
	receivers    receivers
}

func (s sharingPayload) validate() error {
	if s.sourceUserId == "" {
		return sharing.ErrBadRequest
	}
	if s.streamId == "" {
		return sharing.ErrBadRequest
	}
	return nil
}

type getSharingsToGroupsRequest struct {
	receiverGroupIds []string
}

func (req getSharingsToGroupsRequest) validate() error {
	if req.receiverGroupIds == nil {
		return sharing.ErrBadRequest
	}
	if len(req.receiverGroupIds) == 0 {
		return sharing.ErrBadRequest
	}
	for _, rcvGroupId := range req.receiverGroupIds {
		if rcvGroupId == "" {
			return sharing.ErrBadRequest
		}
	}
	return nil
}

type getSharingsRequest struct {
	getSharingsToGroupsRequest
	receiverUserId string
}

func (req getSharingsRequest) validate() error {
	if req.receiverUserId == "" {
		return sharing.ErrBadRequest
	}
	return nil
}
