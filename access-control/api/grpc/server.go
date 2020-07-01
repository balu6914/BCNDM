package grpc

import (
	"context"

	access "github.com/datapace/datapace/access-control"
	accessproto "github.com/datapace/datapace/proto/access"
	commonproto "github.com/datapace/datapace/proto/common"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ accessproto.AccessServiceServer = (*grpcServer)(nil)

type grpcServer struct {
	accessproto.UnimplementedAccessServiceServer
	partners  kitgrpc.Handler
	potential kitgrpc.Handler
}

// NewServer instantiates new Access Control gRPC server.
func NewServer(svc access.Service) accessproto.AccessServiceServer {
	partners := kitgrpc.NewServer(
		partnersEndpoint(svc),
		decodePartnersRequest,
		encodePartnersResponse,
	)

	potential := kitgrpc.NewServer(
		potentialEndpoint(svc),
		decodePartnersRequest,
		encodePartnersResponse,
	)

	return &grpcServer{
		partners:  partners,
		potential: potential,
	}
}

func (s *grpcServer) Partners(ctx context.Context, id *commonproto.ID) (*accessproto.PartnersList, error) {
	_, res, err := s.partners.ServeGRPC(ctx, id)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*accessproto.PartnersList), nil
}

func (s *grpcServer) PotentialPartners(ctx context.Context, id *commonproto.ID) (*accessproto.PartnersList, error) {
	_, res, err := s.potential.ServeGRPC(ctx, id)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*accessproto.PartnersList), nil
}

func decodePartnersRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*commonproto.ID)
	return partnersReq{req.GetValue()}, nil
}

func encodePartnersResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(partnersRes)
	return &accessproto.PartnersList{Value: res.partners}, encodeError(res.err)
}

func encodeError(err error) error {
	switch err {
	case nil:
		return nil
	case access.ErrMalformedEntity:
		return status.Error(codes.InvalidArgument, "received invalid request")
	case access.ErrUnauthorizedAccess:
		return status.Error(codes.Unauthenticated, "failed to identify user from key")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
