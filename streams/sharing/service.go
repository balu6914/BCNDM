package sharing

import (
	"context"
	"fmt"
	log "github.com/datapace/datapace/logger"
	"github.com/datapace/sharing"
	sharingProto "github.com/datapace/sharing/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"time"
)

var (
	logger = log.New(os.Stdout)
)

type (
	// Service is the sharing service
	Service interface {

		// DeleteSharing deletes the sharing of the specified stream for the specified source user
		DeleteSharing(userId string, streamId string) error

		// GetSharings queries for the all Sharing entities where any receiver is either:
		//	(*) specified UserId,
		//	(*) any of the specified groupIds.
		GetSharings(rcvUserId string, rcvGroupIds []string) ([]sharing.Sharing, error)
	}

	sharingService struct {
		client sharingProto.SharingServiceClient
	}
)

func NewService(client sharingProto.SharingServiceClient) Service {
	return sharingService{client: client}
}

func (svc sharingService) DeleteSharing(userId string, streamId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := sharingProto.Sharing{
		SourceUserId: userId,
		StreamId:     streamId,
		Receivers: &sharingProto.Receivers{
			VersionOption: &sharingProto.Receivers_VersionEmpty{VersionEmpty: true},
			GroupIds:      nil,
			UserIds:       nil,
		},
	}
	_, err := svc.client.UpdateReceivers(ctx, &req)
	st, _ := status.FromError(err)
	if st != nil && st.Code() == codes.Unavailable {
		logger.Warn(fmt.Sprintf("Failed to unshare stream %s: sharing service is unavailable (not supported).", streamId))
		return nil
	}
	return err
}

func (svc sharingService) GetSharings(rcvUserId string, rcvGroupIds []string) ([]sharing.Sharing, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := sharingProto.GetSharingsRequest{
		ReceiverGroupIds: rcvGroupIds,
		ReceiverUserId:   rcvUserId,
	}
	resp, err := svc.client.GetSharings(ctx, &req)
	st, _ := status.FromError(err)
	if st != nil && st.Code() == codes.Unavailable {
		logger.Warn("Failed to resolve shared streams: sharing service is unavailable (not supported).")
		return []sharing.Sharing{}, nil
	}
	var results []sharing.Sharing
	for _, s := range resp.Sharings {
		result := convertSharing(s)
		results = append(results, result)
	}
	return results, nil
}

func convertSharing(sp *sharingProto.Sharing) sharing.Sharing {
	rcvs := convertReceivers(sp.Receivers)
	result := sharing.Sharing{
		SourceUserId: sharing.UserId(sp.SourceUserId),
		StreamId:     sharing.StreamId(sp.StreamId),
		Receivers:    *rcvs,
	}
	return result
}

func convertReceivers(rp *sharingProto.Receivers) *sharing.Receivers {
	if rp == nil {
		return nil
	}
	var versionRef *uint64
	switch versionOpt := rp.VersionOption.(type) {
	case *sharingProto.Receivers_VersionEmpty:
		versionRef = nil
	case *sharingProto.Receivers_Version:
		versionRef = &versionOpt.Version
	}
	var rcvUserIds []sharing.UserId
	for _, u := range rp.UserIds {
		rcvUserIds = append(rcvUserIds, sharing.UserId(u))
	}
	var rcvGroupIds []sharing.GroupId
	for _, g := range rp.GroupIds {
		rcvGroupIds = append(rcvGroupIds, sharing.GroupId(g))
	}
	return &sharing.Receivers{
		Version:  versionRef,
		UserIds:  rcvUserIds,
		GroupIds: rcvGroupIds,
	}
}
