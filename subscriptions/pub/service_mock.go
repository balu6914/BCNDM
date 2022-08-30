package pub

type (
	serviceMock struct {
	}
)

func NewServiceMock() Service {
	return serviceMock{}
}

func (svc serviceMock) Publish(evt SubscriptionCreateEvent, toUserId string) (uint64, error) {
	return 0, nil
}
