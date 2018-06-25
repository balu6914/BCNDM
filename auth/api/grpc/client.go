package grpc

import (
	"context"
	"monetasa"

	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

var _ monetasa.AuthServiceClient = (*grpcClient)(nil)

type grpcClient struct {
	identify endpoint.Endpoint
}

// NewClient returns new gRPC client instance.
func NewClient(conn *grpc.ClientConn) monetasa.AuthServiceClient {
	endpoint := kitgrpc.NewClient(
		conn,
		"monetasa.AuthService",
		"Identify",
		encodeIdentifyRequest,
		decodeIdentifyResponse,
		monetasa.UserID{},
	).Endpoint()

	return &grpcClient{endpoint}
}

func (client grpcClient) Identify(ctx context.Context, token *monetasa.Token, _ ...grpc.CallOption) (*monetasa.UserID, error) {
	res, err := client.identify(ctx, identityReq{token.GetValue()})
	if err != nil {
		return nil, err
	}

	ir := res.(identityRes)
	return &monetasa.UserID{Value: ir.id}, ir.err
}

func encodeIdentifyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(identityReq)
	return &monetasa.Token{Value: req.token}, nil
}

func decodeIdentifyResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*monetasa.UserID)
	return identityRes{res.GetValue(), nil}, nil
}
