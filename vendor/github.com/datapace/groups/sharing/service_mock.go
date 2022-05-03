package sharing

type (
	serviceMock struct {
	}
)

func NewServiceMock() Service {
	return serviceMock{}
}

func (svc serviceMock) DeleteReceivers(_ string) error {
	return nil
}
