package grpc

import (
	"context"
	"datapace"

	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

var _ datapace.ExecutionsServiceClient = (*grpcClient)(nil)

type grpcClient struct {
	createAlgo endpoint.Endpoint
	createData endpoint.Endpoint
}

// NewClient returns new gRPC client instance.
func NewClient(conn *grpc.ClientConn) datapace.ExecutionsServiceClient {
	createAlgo := kitgrpc.NewClient(
		conn,
		"datapace.ExecutionsService",
		"CreateAlgorithm",
		encodeAlgoReq,
		decodeCreateRes,
		empty.Empty{},
	).Endpoint()

	createData := kitgrpc.NewClient(
		conn,
		"datapace.ExecutionsService",
		"CreateDataset",
		encodeDataReq,
		decodeCreateRes,
		empty.Empty{},
	).Endpoint()

	return &grpcClient{
		createAlgo: createAlgo,
		createData: createData,
	}
}

func (client grpcClient) CreateAlgorithm(ctx context.Context, algo *datapace.Algorithm, _ ...grpc.CallOption) (*empty.Empty, error) {
	req := algoReq{
		id:         algo.GetId(),
		name:       algo.GetName(),
		path:       algo.GetPath(),
		modelToken: algo.GetModelToken(),
		modelName:  algo.GetModelName(),
	}

	res, err := client.createAlgo(ctx, req)
	if err != nil {
		return nil, err
	}

	cr := res.(createRes)
	return &empty.Empty{}, cr.err
}

func (client grpcClient) CreateDataset(ctx context.Context, data *datapace.Dataset, _ ...grpc.CallOption) (*empty.Empty, error) {
	req := dataReq{
		id:   data.GetId(),
		path: data.GetPath(),
	}

	res, err := client.createData(ctx, req)
	if err != nil {
		return nil, err
	}

	cr := res.(createRes)
	return &empty.Empty{}, cr.err
}

func encodeAlgoReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(algoReq)
	return &datapace.Algorithm{
		Id:         req.id,
		Name:       req.name,
		Path:       req.path,
		ModelToken: req.modelToken,
		ModelName:  req.modelName,
	}, nil
}

func encodeDataReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(dataReq)
	return &datapace.Dataset{
		Id:   req.id,
		Path: req.path,
	}, nil
}

func decodeCreateRes(_ context.Context, grpcRes interface{}) (interface{}, error) {
	return createRes{}, nil
}