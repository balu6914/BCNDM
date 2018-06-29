package grpc

import (
	"context"
	"errors"
	"monetasa"
	"monetasa/transactions"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errMalformedEntity = errors.New("malformed entity")

var _ monetasa.TransactionsServiceServer = (*grpcServer)(nil)

type grpcServer struct {
	handler kitgrpc.Handler
}

// NewServer instantiates new Auth gRPC server.
func NewServer(svc transactions.Service) monetasa.TransactionsServiceServer {
	handler := kitgrpc.NewServer(
		createUserEndpoint(svc),
		decodeCreateUserRequest,
		encodeCreateUserResponse,
	)

	return &grpcServer{handler}
}

func (s grpcServer) CreateUser(ctx context.Context, user *monetasa.User) (*monetasa.Key, error) {
	_, res, err := s.handler.ServeGRPC(ctx, user)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*monetasa.Key), nil
}

func decodeCreateUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*monetasa.User)
	return createUserReq{id: req.GetId(), secret: req.GetSecret()}, nil
}

func encodeCreateUserResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(createUserRes)
	return &monetasa.Key{Value: res.key}, encodeError(res.err)
}

func encodeError(err error) error {
	switch err {
	case nil:
		return nil
	case errMalformedEntity:
		return status.Error(codes.InvalidArgument, "received invalid user request")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
