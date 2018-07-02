package grpc

import (
	"context"
	"monetasa"
	"monetasa/auth"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ monetasa.AuthServiceServer = (*grpcServer)(nil)

type grpcServer struct {
	handler kitgrpc.Handler
}

// NewServer instantiates new Auth gRPC server.
func NewServer(svc auth.Service) monetasa.AuthServiceServer {
	handler := kitgrpc.NewServer(
		identifyEndpoint(svc),
		decodeIdentifyRequest,
		encodeIdentifyResponse,
	)

	return &grpcServer{handler}
}

func (s *grpcServer) Identify(ctx context.Context, token *monetasa.Token) (*monetasa.UserID, error) {
	_, res, err := s.handler.ServeGRPC(ctx, token)
	if err != nil {
		return nil, encodeError(err)
	}
	return res.(*monetasa.UserID), nil
}

func decodeIdentifyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*monetasa.Token)
	return identityReq{req.GetValue()}, nil
}

func encodeIdentifyResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(identityRes)
	return &monetasa.UserID{Value: res.id}, encodeError(res.err)
}

func encodeError(err error) error {
	switch err {
	case nil:
		return nil
	case auth.ErrMalformedEntity:
		return status.Error(codes.InvalidArgument, "received invalid key request")
	case auth.ErrUnauthorizedAccess:
		return status.Error(codes.Unauthenticated, "failed to identify user from key")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
