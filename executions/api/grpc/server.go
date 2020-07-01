package grpc

import (
	"context"

	"github.com/datapace/datapace/executions"

	executionsproto "github.com/datapace/datapace/proto/executions"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ executionsproto.ExecutionsServiceServer = (*grpcServer)(nil)

type grpcServer struct {
	executionsproto.UnimplementedExecutionsServiceServer
	createAlgo kitgrpc.Handler
	createData kitgrpc.Handler
}

// NewServer instantiates new Executions gRPC server.
func NewServer(svc executions.Service) executionsproto.ExecutionsServiceServer {
	createAlgo := kitgrpc.NewServer(
		createAlgoEndpoint(svc),
		decodeAlgoReq,
		encodeCreateRes,
	)
	createData := kitgrpc.NewServer(
		createDataEndpoint(svc),
		decodeDataReq,
		encodeCreateRes,
	)

	return &grpcServer{
		createAlgo: createAlgo,
		createData: createData,
	}
}

func (s *grpcServer) CreateAlgorithm(ctx context.Context, algo *executionsproto.Algorithm) (*empty.Empty, error) {
	_, res, err := s.createAlgo.ServeGRPC(ctx, algo)
	if err != nil {
		return nil, encodeError(err)
	}
	return res.(*empty.Empty), nil
}

func (s *grpcServer) CreateDataset(ctx context.Context, data *executionsproto.Dataset) (*empty.Empty, error) {
	_, res, err := s.createData.ServeGRPC(ctx, data)
	if err != nil {
		return nil, encodeError(err)
	}
	return res.(*empty.Empty), nil
}

func decodeAlgoReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*executionsproto.Algorithm)
	return algoReq{
		id:       req.GetId(),
		name:     req.GetName(),
		metadata: req.GetMetadata(),
	}, nil
}

func decodeDataReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*executionsproto.Dataset)
	return dataReq{
		id:       req.GetId(),
		metadata: req.GetMetadata(),
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
	case executions.ErrMalformedData:
		return status.Error(codes.InvalidArgument, "received invalid request")
	case executions.ErrConflict:
		return status.Error(codes.AlreadyExists, "entity already exists")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
