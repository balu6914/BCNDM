package offers

import (
	"context"
	"errors"
	"fmt"

	offersProto "github.com/datapace/datapace/proto/offers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OfferPrice struct {
	Price uint64
}

type OfferKey struct {
	StreamId string
	BuyerId  string
}

type Service interface {
	GetPrice(context.Context, OfferKey) (OfferPrice, error)
}

var ErrInternal = errors.New("internal failure")

var ErrNotFound = errors.New("not found")

var ErrNotAvailable = errors.New("service not available")

type service struct {
	client offersProto.OffersServiceClient
}

func NewService(client offersProto.OffersServiceClient) Service {
	return service{
		client: client,
	}
}

func (svc service) GetPrice(ctx context.Context, k OfferKey) (p OfferPrice, err error) {
	req := &offersProto.GetOfferPriceRequest{
		StreamId: k.StreamId,
		BuyerId:  k.BuyerId,
	}
	var resp *offersProto.GetOfferPriceResponse
	resp, err = svc.client.GetOfferPrice(ctx, req)
	if err != nil {
		err = decodeError(err)
	} else {
		p.Price = resp.GetOfferPrice()
	}
	return
}

func decodeError(src error) (dst error) {
	st, _ := status.FromError(src)
	switch st.Code() {
	case codes.Unavailable:
		dst = fmt.Errorf("%w: %s", ErrNotAvailable, src)
	case codes.NotFound:
		dst = fmt.Errorf("%w: %s", ErrNotFound, src)
	default:
		dst = fmt.Errorf("%w: %s", ErrInternal, src)
	}
	return
}
