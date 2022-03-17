package groups

type (
	serviceMock struct{}
)

func NewServiceMock() Service {
	return serviceMock{}
}

func (svc serviceMock) GetUserGroups(userId string) (groupIds []string, err error) {
	if userId == "userInSomeGroups" {
		return []string{"group1"}, nil
	}
	if userId == "sharingReceiverUserInSomeGroups" {
		return []string{"group1"}, nil
	}
	return []string{}, nil
}
