package grpc

import (
	"context"
	"errors"
	"monetasa"
	"monetasa/transactions"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/golang/protobuf/ptypes/empty"
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

func (s grpcServer) CreateUser(ctx context.Context, user *monetasa.ID) (*empty.Empty, error) {
	_, res, err := s.handler.ServeGRPC(ctx, user)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*empty.Empty), nil
}

func decodeCreateUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*monetasa.ID)
	return createUserReq{id: req.GetValue()}, nil
}

func encodeCreateUserResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(createUserRes)
	return &empty.Empty{}, encodeError(res.err)
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
