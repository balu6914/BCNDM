package pub

type (
	serviceMock struct {
	}
)

func NewServiceMock() Service {
	return serviceMock{}
}

func (svc serviceMock) PublishSubscriptionCreated(evt interface{}, toUserId string) (uint64, error) {
	return 0, nil
}
