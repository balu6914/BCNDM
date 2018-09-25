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
	identify kitgrpc.Handler
	email    kitgrpc.Handler
}

// NewServer instantiates new Auth gRPC server.
func NewServer(svc auth.Service) monetasa.AuthServiceServer {
	identify := kitgrpc.NewServer(
		identifyEndpoint(svc),
		decodeIdentifyRequest,
		encodeIdentifyResponse,
	)

	email := kitgrpc.NewServer(
		emailEndpoint(svc),
		decodeEmailRequest,
		encodeEmailResponse,
	)

	return &grpcServer{identify, email}
}

func (s *grpcServer) Identify(ctx context.Context, token *monetasa.Token) (*monetasa.UserID, error) {
	_, res, err := s.identify.ServeGRPC(ctx, token)
	if err != nil {
		return nil, encodeError(err)
	}
	return res.(*monetasa.UserID), nil
}

func (s *grpcServer) Email(ctx context.Context, id *monetasa.UserID) (*monetasa.UserEmail, error) {
	_, res, err := s.email.ServeGRPC(ctx, id)
	if err != nil {
		return nil, encodeError(err)
	}
	return res.(*monetasa.UserEmail), nil
}

func decodeIdentifyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*monetasa.Token)
	return identityReq{req.GetValue()}, nil
}

func encodeIdentifyResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(identityRes)
	return &monetasa.UserID{Value: res.id}, encodeError(res.err)
}

func decodeEmailRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*monetasa.UserID)
	return emailReq{req.GetValue()}, nil
}

func encodeEmailResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(emailRes)
	return &monetasa.UserEmail{Value: res.email}, encodeError(res.err)
}

func encodeError(err error) error {
	switch err {
	case nil:
		return nil
	case auth.ErrMalformedEntity:
		return status.Error(codes.InvalidArgument, "received invalid request")
	case auth.ErrUnauthorizedAccess:
		return status.Error(codes.Unauthenticated, "failed to identify user from key")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
