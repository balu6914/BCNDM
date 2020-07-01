package grpc

import (
	"context"

	commonproto "github.com/datapace/datapace/proto/common"
	transactionsproto "github.com/datapace/datapace/proto/transactions"
	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

var _ transactionsproto.TransactionsServiceClient = (*grpcClient)(nil)

type grpcClient struct {
	createUser endpoint.Endpoint
	transfer   endpoint.Endpoint
}

// NewClient returns new gRPC client instance.
func NewClient(conn *grpc.ClientConn) transactionsproto.TransactionsServiceClient {
	createUser := kitgrpc.NewClient(
		conn,
		"datapace.TransactionsService",
		"CreateUser",
		encodeCreateUserRequest,
		decodeCreateUserResponse,
		empty.Empty{},
	).Endpoint()

	transfer := kitgrpc.NewClient(
		conn,
		"datapace.TransactionsService",
		"Transfer",
		encodeTransferRequest,
		decodeTransferResponse,
		empty.Empty{},
	).Endpoint()

	return &grpcClient{
		createUser: createUser,
		transfer:   transfer,
	}
}

func (client grpcClient) CreateUser(ctx context.Context, user *commonproto.ID, _ ...grpc.CallOption) (*empty.Empty, error) {
	res, err := client.createUser(ctx, createUserReq{id: user.GetValue()})
	if err != nil {
		return nil, err
	}

	cur := res.(createUserRes)
	return &empty.Empty{}, cur.err
}

func (client grpcClient) Transfer(ctx context.Context, td *transactionsproto.TransferData, _ ...grpc.CallOption) (*empty.Empty, error) {
	req := transferReq{
		streamID: td.GetStreamID(),
		from:     td.GetFrom(),
		to:       td.GetTo(),
		value:    td.GetValue(),
	}

	res, err := client.transfer(ctx, req)
	if err != nil {
		return nil, err
	}
	tr := res.(transferRes)
	return &empty.Empty{}, tr.err
}

func encodeCreateUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(createUserReq)
	return &commonproto.ID{Value: req.id}, nil
}

func decodeCreateUserResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	return createUserRes{}, nil
}

func encodeTransferRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(transferReq)
	return &transactionsproto.TransferData{
		StreamID: req.streamID,
		From:     req.from,
		To:       req.to,
		Value:    req.value,
	}, nil
}

func decodeTransferResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	return transferRes{}, nil
}
