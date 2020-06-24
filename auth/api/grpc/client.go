package grpc

import (
	"context"

	"github.com/datapace/datapace"

	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

var _ datapace.AuthServiceClient = (*grpcClient)(nil)

type grpcClient struct {
	identify endpoint.Endpoint
	email    endpoint.Endpoint
	exists   endpoint.Endpoint
}

// NewClient returns new gRPC client instance.
func NewClient(conn *grpc.ClientConn) datapace.AuthServiceClient {
	identify := kitgrpc.NewClient(
		conn,
		"datapace.AuthService",
		"Identify",
		encodeIdentifyRequest,
		decodeIdentifyResponse,
		datapace.ID{},
	).Endpoint()

	email := kitgrpc.NewClient(
		conn,
		"datapace.AuthService",
		"Email",
		encodeIdentifyRequest,
		decodeEmailResponse,
		datapace.UserEmail{},
	).Endpoint()

	exists := kitgrpc.NewClient(
		conn,
		"datapace.AuthService",
		"Exists",
		encodeExistsRequest,
		decodeExistsResponse,
		empty.Empty{},
	).Endpoint()

	return &grpcClient{
		identify: identify,
		email:    email,
		exists:   exists,
	}
}

func (client grpcClient) Identify(ctx context.Context, token *datapace.Token, _ ...grpc.CallOption) (*datapace.ID, error) {
	res, err := client.identify(ctx, identityReq{token.GetValue()})
	if err != nil {
		return nil, err
	}

	idRes := res.(identityRes)
	return &datapace.ID{Value: idRes.id}, idRes.err
}

func (client grpcClient) Email(ctx context.Context, token *datapace.Token, _ ...grpc.CallOption) (*datapace.UserEmail, error) {
	res, err := client.email(ctx, identityReq{token.GetValue()})
	if err != nil {
		return nil, err
	}

	emailRes := res.(emailRes)
	return &datapace.UserEmail{Email: emailRes.email, ContactEmail: emailRes.contactEmail}, emailRes.err
}

func (client grpcClient) Exists(ctx context.Context, id *datapace.ID, _ ...grpc.CallOption) (*empty.Empty, error) {
	res, err := client.exists(ctx, existsReq{id.GetValue()})
	if err != nil {
		return nil, err
	}

	existsRes := res.(existsRes)
	return &empty.Empty{}, existsRes.err
}

func encodeIdentifyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(identityReq)
	return &datapace.Token{Value: req.token}, nil
}

func decodeIdentifyResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*datapace.ID)
	return identityRes{res.GetValue(), nil}, nil
}

func decodeEmailResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*datapace.UserEmail)
	return emailRes{res.GetEmail(), res.GetContactEmail(), nil}, nil
}

func encodeExistsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(existsReq)
	return &datapace.ID{Value: req.id}, nil
}

func decodeExistsResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	return existsRes{nil}, nil
}
