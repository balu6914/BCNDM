package grpc

import (
	"context"
	"datapace"
	"datapace/auth"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ datapace.AuthServiceServer = (*grpcServer)(nil)

type grpcServer struct {
	identify kitgrpc.Handler
	email    kitgrpc.Handler
	partners kitgrpc.Handler
}

// NewServer instantiates new Auth gRPC server.
func NewServer(svc auth.Service) datapace.AuthServiceServer {
	identify := kitgrpc.NewServer(
		identifyEndpoint(svc),
		decodeIdentifyRequest,
		encodeIdentifyResponse,
	)

	email := kitgrpc.NewServer(
		emailEndpoint(svc),
		decodeIdentifyRequest,
		encodeEmailResponse,
	)

	partners := kitgrpc.NewServer(
		partnersEndpoint(svc),
		decodePartnersRequest,
		encodePartnersResponse,
	)

	return &grpcServer{
		identify: identify,
		email:    email,
		partners: partners,
	}
}

func (s *grpcServer) Identify(ctx context.Context, token *datapace.Token) (*datapace.UserID, error) {
	_, res, err := s.identify.ServeGRPC(ctx, token)
	if err != nil {
		return nil, encodeError(err)
	}
	return res.(*datapace.UserID), nil
}

func (s *grpcServer) Email(ctx context.Context, id *datapace.Token) (*datapace.UserEmail, error) {
	_, res, err := s.email.ServeGRPC(ctx, id)
	if err != nil {
		return nil, encodeError(err)
	}
	return res.(*datapace.UserEmail), nil
}

func (s *grpcServer) Partners(ctx context.Context, id *datapace.UserID) (*datapace.PartnersList, error) {
	_, res, err := s.partners.ServeGRPC(ctx, id)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*datapace.PartnersList), nil
}

func decodeIdentifyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*datapace.Token)
	return identityReq{req.GetValue()}, nil
}

func encodeIdentifyResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(identityRes)
	return &datapace.UserID{Value: res.id}, encodeError(res.err)
}

func encodeEmailResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(emailRes)
	return &datapace.UserEmail{Email: res.email, ContactEmail: res.contactEmail}, encodeError(res.err)
}

func decodePartnersRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*datapace.UserID)
	return identityReq{req.GetValue()}, nil
}

func encodePartnersResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(partnersRes)
	return &datapace.PartnersList{Value: res.partners}, encodeError(res.err)
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
