package accessv2

import (
	"context"
	"errors"
	"fmt"
	accessProtoV2 "github.com/datapace/datapace/proto/accessv2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Access struct {
	Key   Key
	State State
	Time  time.Time
}

type Key struct {
	ConsumerId string
	ProviderId string
	ProductId  string
}

type State int

const (
	StatePending = iota
	StateApproved
	StateCancelled
)

func (s State) String() string {
	return [...]string{
		"pending",
		"approved",
		"cancelled",
	}[s]
}

type Service interface {
	Get(ctx context.Context, k Key) (a Access, err error)
}

var ErrInternal = errors.New("internal failure")

var ErrNotFound = errors.New("not found")

var ErrNotAvailable = errors.New("service not available")

type service struct {
	client accessProtoV2.ServiceClient
}

func NewService(client accessProtoV2.ServiceClient) Service {
	return service{
		client: client,
	}
}

func (svc service) Get(ctx context.Context, k Key) (a Access, err error) {
	req := &accessProtoV2.Key{
		ConsumerId: k.ConsumerId,
		ProviderId: k.ProviderId,
		ProductId:  k.ProductId,
	}
	var resp *accessProtoV2.Access
	resp, err = svc.client.Get(ctx, req)
	if err != nil {
		err = decodeError(err)
	} else {
		decodeAccess(resp, &a)
	}
	return
}

func decodeAccess(resp *accessProtoV2.Access, a *Access) {
	respKey := resp.GetKey()
	a.Key = Key{
		ConsumerId: respKey.GetConsumerId(),
		ProviderId: respKey.GetProviderId(),
		ProductId:  respKey.GetProductId(),
	}
	a.State = State(resp.GetState().Number())
	a.Time = resp.GetTime().AsTime()
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
