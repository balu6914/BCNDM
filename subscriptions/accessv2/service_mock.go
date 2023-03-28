package accessv2

import (
	"context"
	"time"
)

type serviceMock struct {
}

func NewServiceMock() Service {
	return serviceMock{}
}

func (sm serviceMock) Get(ctx context.Context, k Key) (a Access, err error) {
	switch k.ConsumerId {
	case "unavailable":
		err = ErrNotAvailable
	case "fail":
		err = ErrInternal
	case "pending":
		a = Access{
			Key:   k,
			State: StatePending,
			Time:  time.Date(2022, time.December, 16, 11, 41, 25, 0, time.UTC),
		}
	case "approved":
		a = Access{
			Key:   k,
			State: StateApproved,
			Time:  time.Date(2022, time.December, 16, 11, 41, 25, 0, time.UTC),
		}
	case "cancelled":
		a = Access{
			Key:   k,
			State: StateCancelled,
			Time:  time.Date(2022, time.December, 16, 11, 41, 25, 0, time.UTC),
		}
	default:
		err = ErrNotFound
	}
	return
}
