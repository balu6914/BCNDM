package grpc

import (
	"context"
	"monetasa"

	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

var _ monetasa.TransactionsServiceClient = (*grpcClient)(nil)

type grpcClient struct {
	createUser endpoint.Endpoint
}

// NewClient returns new gRPC client instance.
func NewClient(conn *grpc.ClientConn) monetasa.TransactionsServiceClient {
	endpoint := kitgrpc.NewClient(
		conn,
		"monetasa.TransactionsService",
		"CreateUser",
		encodeCreateUserRequest,
		decodeCreateUserResponse,
		empty.Empty{},
	).Endpoint()

	return &grpcClient{endpoint}
}

func (client grpcClient) CreateUser(ctx context.Context, user *monetasa.ID, _ ...grpc.CallOption) (*empty.Empty, error) {
	res, err := client.createUser(ctx, createUserReq{id: user.GetValue()})
	if err != nil {
		return nil, err
	}

	cur := res.(createUserRes)
	return &empty.Empty{}, cur.err
}

func encodeCreateUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(createUserReq)
	return &monetasa.ID{Value: req.id}, nil
}

func decodeCreateUserResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	return createUserRes{}, nil
}
