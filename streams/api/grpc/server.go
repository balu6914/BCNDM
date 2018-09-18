package grpc

import (
	"context"
	"monetasa"
	"monetasa/streams"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ monetasa.StreamsServiceServer = (*grpcServer)(nil)

type grpcServer struct {
	handler kitgrpc.Handler
}

// NewServer instantiates new Auth gRPC server.
func NewServer(svc streams.Service) monetasa.StreamsServiceServer {
	return &grpcServer{
		handler: kitgrpc.NewServer(
			oneEndpoint(svc),
			decodeOneRequest,
			encodeOneResponse,
		),
	}
}

func (s grpcServer) One(ctx context.Context, id *monetasa.ID) (*monetasa.Stream, error) {
	_, res, err := s.handler.ServeGRPC(ctx, id)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*monetasa.Stream), nil
}

func decodeOneRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*monetasa.ID)
	return oneReq{id: req.GetValue()}, nil
}

func encodeOneResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(oneRes)
	stream := monetasa.Stream{
		Id:      res.id,
		Name:    res.name,
		Owner:   res.owner,
		Url:     res.url,
		Price:   res.price,
		Bq:      res.bq,
		Project: res.project,
		Dataset: res.dataset,
		Table:   res.table,
		Fields:  res.fields,
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
