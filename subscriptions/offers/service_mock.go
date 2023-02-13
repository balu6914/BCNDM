package offers

import (
	"context"
)

const price uint64 = 63

type serviceMock struct {
}

func NewServiceMock() Service {
	return serviceMock{}
}

func (sm serviceMock) GetPrice(ctx context.Context, k OfferKey) (p OfferPrice, err error) {
	p.Price = price
	return
}
