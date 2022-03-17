package sharing

import (
	"github.com/datapace/sharing"
)

var (
	version1 = uint64(1)
)

type (
	serviceMock struct {
	}
)

func NewServiceMock() Service {
	return serviceMock{}
}

func (svc serviceMock) DeleteSharing(userId string, streamId string) error {
	return nil
}

func (svc serviceMock) GetSharings(rcvUserId string, rcvGroupIds []string) ([]sharing.Sharing, error) {
	var sharings []sharing.Sharing
	if rcvUserId == "sharingReceiverUser" || rcvUserId == "sharingReceiverUserInSomeGroups" {
		s := sharing.Sharing{
			SourceUserId: "user0",
			StreamId:     "stream0",
			Receivers: sharing.Receivers{
				Version: &version1,
				UserIds: []sharing.UserId{
					"sharingReceiverUser",
					"sharingReceiverUserInSomeGroups",
				},
			},
		}
		sharings = append(sharings, s)
	}
	for _, rcvGroupId := range rcvGroupIds {
		if rcvGroupId == "group1" {
			s := sharing.Sharing{
				SourceUserId: "user1",
				StreamId:     "stream1",
				Receivers: sharing.Receivers{
					Version:  &version1,
					GroupIds: []sharing.GroupId{"group1"},
				},
			}
			sharings = append(sharings, s)
		}
	}
	return sharings, nil
}
