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
	email    endpoint.Endpoint
}

// NewClient returns new gRPC client instance.
func NewClient(conn *grpc.ClientConn) monetasa.AuthServiceClient {
	identify := kitgrpc.NewClient(
		conn,
		"monetasa.AuthService",
		"Identify",
		encodeIdentifyRequest,
		decodeIdentifyResponse,
		monetasa.UserID{},
	).Endpoint()

	email := kitgrpc.NewClient(
		conn,
		"monetasa.AuthService",
		"Email",
		encodeEmailRequest,
		decodeEmailResponse,
		monetasa.UserEmail{},
	).Endpoint()
	return &grpcClient{identify, email}
}

func (client grpcClient) Identify(ctx context.Context, token *monetasa.Token, _ ...grpc.CallOption) (*monetasa.UserID, error) {
	res, err := client.identify(ctx, identityReq{token.GetValue()})
	if err != nil {
		return nil, err
	}

	ir := res.(identityRes)
	return &monetasa.UserID{Value: ir.id}, ir.err
}

func (client grpcClient) Email(ctx context.Context, id *monetasa.UserID, _ ...grpc.CallOption) (*monetasa.UserEmail, error) {
	res, err := client.email(ctx, emailReq{id.GetValue()})
	if err != nil {
		return nil, err
	}

	er := res.(emailRes)
	return &monetasa.UserEmail{Value: er.email}, er.err
}

func encodeIdentifyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(identityReq)
	return &monetasa.Token{Value: req.token}, nil
}

func decodeIdentifyResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*monetasa.UserID)
	return identityRes{res.GetValue(), nil}, nil
}

func encodeEmailRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(emailReq)
	return &monetasa.UserID{Value: req.id}, nil
}

func decodeEmailResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*monetasa.UserEmail)
	return emailRes{res.GetValue(), nil}, nil
}
