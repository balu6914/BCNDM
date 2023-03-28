package offers

import (
	"context"
)

const price uint64 = 15

type serviceMock struct {
}

func NewServiceMock() Service {
	return serviceMock{}
}

func (sm serviceMock) GetPrice(ctx context.Context, k OfferKey) (p OfferPrice, err error) {
	switch k.StreamId {
	case "offer_accepted":
		p.Price = price
	case "unavailable":
		err = ErrNotAvailable
	default:
		err = ErrNotAvailable
	}
	return
}
