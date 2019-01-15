package grpc

import (
	"context"
	"datapace"
	"datapace/streams"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ datapace.StreamsServiceServer = (*grpcServer)(nil)

type grpcServer struct {
	handler kitgrpc.Handler
}

// NewServer instantiates new Auth gRPC server.
func NewServer(svc streams.Service) datapace.StreamsServiceServer {
	return &grpcServer{
		handler: kitgrpc.NewServer(
			oneEndpoint(svc),
			decodeOneRequest,
			encodeOneResponse,
		),
	}
}

func (s grpcServer) One(ctx context.Context, id *datapace.ID) (*datapace.Stream, error) {
	_, res, err := s.handler.ServeGRPC(ctx, id)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*datapace.Stream), nil
}

func decodeOneRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*datapace.ID)
	return oneReq{id: req.GetValue()}, nil
}

func encodeOneResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(oneRes)
	stream := datapace.Stream{
		Id:       res.id,
		Name:     res.name,
		Owner:    res.owner,
		Url:      res.url,
		Price:    res.price,
		External: res.external,
		Project:  res.project,
		Dataset:  res.dataset,
		Table:    res.table,
		Fields:   res.fields,
	}

	return &stream, nil
}

func encodeError(err error) error {
	switch err {
	case nil:
		return nil
	case streams.ErrMalformedData:
		return status.Error(codes.InvalidArgument, "received invalid id")
	case streams.ErrNotFound:
		return status.Error(codes.NotFound, "stream doesn't exist")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
