package groups

import (
	"context"
	"fmt"
	log "github.com/datapace/datapace/logger"
	groupsProto "github.com/datapace/groups/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"time"
)

type (

	// Service is the groups service level API
	Service interface {

		// GetUserGroups returns all groups those the specified user belongs to
		GetUserGroups(userId string) (groupIds []string, err error)
	}

	service struct {
		client groupsProto.GroupsServiceClient
	}
)

var (
	logger = log.New(os.Stdout)
)

func NewService(client groupsProto.GroupsServiceClient) Service {
	return service{client: client}
}

func (svc service) GetUserGroups(userId string) (groupIds []string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := groupsProto.GetUserGroupsRequest{Uid: userId}
	resp, err := svc.client.GetUserGroups(ctx, &req)
	st, _ := status.FromError(err)
	if st != nil && st.Code() == codes.Unavailable {
		logger.Warn(fmt.Sprintf("Failed to resolve user %s groups, groups service is unavailable (not supported).", userId))
		return []string{}, nil
	}
	if st != nil && st.Code() == codes.NotFound {
		return []string{}, nil // not a failure if a user is not in any group
	}
	if err != nil {
		return []string{}, err
	}
	return resp.Gids, nil
}
