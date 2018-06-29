package grpc

import (
	"context"
	"monetasa"

	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
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
		monetasa.Key{},
	).Endpoint()

	return &grpcClient{endpoint}
}

func (client grpcClient) CreateUser(ctx context.Context, user *monetasa.User, _ ...grpc.CallOption) (*monetasa.Key, error) {
	res, err := client.createUser(ctx, createUserReq{id: user.GetId(), secret: user.GetSecret()})
	if err != nil {
		return nil, err
	}

	cur := res.(createUserRes)
	return &monetasa.Key{Value: cur.key}, cur.err
}

func encodeCreateUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(createUserReq)
	return &monetasa.User{Id: req.id, Secret: req.secret}, nil
}

func decodeCreateUserResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*monetasa.Key)
	return createUserRes{res.GetValue(), nil}, nil
}
