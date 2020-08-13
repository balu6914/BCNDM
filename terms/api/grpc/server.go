package grpc

import (
	"context"
	termsproto "github.com/datapace/datapace/proto/terms"
	"github.com/datapace/datapace/terms"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ termsproto.TermsServiceServer = (*grpcServer)(nil)

type grpcServer struct {
	termsproto.UnimplementedTermsServiceServer
	createTerms kitgrpc.Handler
}

// NewServer instantiates new Terms gRPC server.
func NewServer(svc terms.Service) termsproto.TermsServiceServer {
	createTerms := kitgrpc.NewServer(
		createTermsEndpoint(svc),
		decodeCreateTermsReq,
		encodeCreateRes,
	)

	return &grpcServer{
		createTerms: createTerms,
	}
}

func (s *grpcServer) CreateTerms(ctx context.Context, terms *termsproto.Terms) (*empty.Empty, error) {
	_, res, err := s.createTerms.ServeGRPC(ctx, terms)
	if err != nil {
		return nil, encodeError(err)
	}
	return res.(*empty.Empty), nil
}

func decodeCreateTermsReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*termsproto.Terms)
	return termsReq{
		termsUrl: req.GetUrl(),
		streamId: req.GetStreamId(),
	}, nil
}

func encodeCreateRes(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(createRes)
	return &empty.Empty{}, encodeError(res.err)
}

func encodeError(err error) error {
	switch err {
	case nil:
		return nil
	case terms.ErrMalformedEntity:
		return status.Error(codes.InvalidArgument, "received invalid request")
	case terms.ErrConflict:
		return status.Error(codes.AlreadyExists, "entity already exists")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
