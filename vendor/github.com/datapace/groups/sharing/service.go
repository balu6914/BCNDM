package sharing

import (
	"context"
	"fmt"
	log "github.com/datapace/datapace/logger"
	sharingProto "github.com/datapace/sharing/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"sync"
	"time"
)

var (
	logger = log.New(os.Stdout)
)

type (
	// Service is the sharing service
	Service interface {

		// DeleteReceivers deletes the specified receiver group id from all sharings
		DeleteReceivers(string) error
	}

	sharingService struct {
		client sharingProto.SharingServiceClient
	}
)

func NewService(client sharingProto.SharingServiceClient) Service {
	return sharingService{client: client}
}

func (svc sharingService) DeleteReceivers(groupId string) error {
	for {
		sharings, err := svc.getAllSharingsToGroup(groupId)
		st, _ := status.FromError(err)
		if st != nil && st.Code() == codes.Unavailable {
			logger.Warn("Sharing service is unavailable")
			return nil
		}
		if err != nil {
			return err
		}
		if len(sharings) == 0 {
			break
		}
		var latch sync.WaitGroup
		for _, s := range sharings {
			sharingUpdate := cleanReceiverGroup(s, groupId)
			latch.Add(1)
			go svc.update(&sharingUpdate, &latch)
		}
		latch.Wait()
	}
	return nil
}

func (svc sharingService) getAllSharingsToGroup(groupId string) ([]*sharingProto.Sharing, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := sharingProto.GetSharingsToGroupsRequest{ReceiverGroupIds: []string{groupId}}
	resp, err := svc.client.GetSharingsToGroups(ctx, &req)
	if err != nil {
		return []*sharingProto.Sharing{}, err
	}
	return resp.Sharings, nil
}

func (svc sharingService) update(s *sharingProto.Sharing, wg *sync.WaitGroup) {
	defer wg.Done()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := svc.client.UpdateReceivers(ctx, s)
	switch err {
	case nil:
		logger.Warn(fmt.Sprintf("Update sharing success: %v", s))
	case status.Error(codes.NotFound, "entity not found or version mismatch"):
		logger.Warn(fmt.Sprintf("Update sharing failed due to version mismatch/deletion, going to retry: %v", s))
	default:
		logger.Error(fmt.Sprintf("Update sharing failed: %v, error: %s", s, err))
	}
}

func cleanReceiverGroup(s *sharingProto.Sharing, groupId string) sharingProto.Sharing {
	r := s.Receivers
	var rcvGroupIds []string
	for _, g := range r.GroupIds {
		if g != groupId {
			rcvGroupIds = append(rcvGroupIds, g)
		}
	}
	receiversUppdate := sharingProto.Receivers{
		VersionOption: r.VersionOption,
		GroupIds:      rcvGroupIds,
		UserIds:       r.UserIds,
	}
	sharingUpdate := sharingProto.Sharing{
		SourceUserId: s.SourceUserId,
		StreamId:     s.StreamId,
		Receivers:    &receiversUppdate,
	}
	return sharingUpdate
}
