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
	createUser kitgrpc.Handler
	transfer   kitgrpc.Handler
}

// NewServer instantiates new Auth gRPC server.
func NewServer(svc transactions.Service) monetasa.TransactionsServiceServer {
	createUser := kitgrpc.NewServer(
		createUserEndpoint(svc),
		decodeCreateUserRequest,
		encodeCreateUserResponse,
	)

	transfer := kitgrpc.NewServer(
		transferEndpoint(svc),
		decodeTransferRequest,
		encodeTransferResponse,
	)

	return &grpcServer{
		createUser: createUser,
		transfer:   transfer,
	}
}

func (s grpcServer) CreateUser(ctx context.Context, user *monetasa.ID) (*empty.Empty, error) {
	_, res, err := s.createUser.ServeGRPC(ctx, user)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*empty.Empty), nil
}

func (s grpcServer) Transfer(ctx context.Context, td *monetasa.TransferData) (*empty.Empty, error) {
	_, res, err := s.transfer.ServeGRPC(ctx, td)
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

func decodeTransferRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*monetasa.TransferData)
	return transferReq{
		from:  req.GetFrom(),
		to:    req.GetTo(),
		value: req.GetValue(),
	}, nil
}

func encodeTransferResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(transferRes)
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
