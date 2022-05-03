package grpc

import "github.com/datapace/groups"

type getUserGroupsRequest struct {
	uid string
}

func (r getUserGroupsRequest) validate() error {
	if r.uid == "" {
		return groups.ErrBadRequest
	}
	return nil
}
