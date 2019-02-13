package grpc

import (
	"context"
	"datapace"

	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

var _ datapace.AuthServiceClient = (*grpcClient)(nil)

type grpcClient struct {
	identify endpoint.Endpoint
	email    endpoint.Endpoint
	partners endpoint.Endpoint
}

// NewClient returns new gRPC client instance.
func NewClient(conn *grpc.ClientConn) datapace.AuthServiceClient {
	identify := kitgrpc.NewClient(
		conn,
		"datapace.AuthService",
		"Identify",
		encodeIdentifyRequest,
		decodeIdentifyResponse,
		datapace.UserID{},
	).Endpoint()

	email := kitgrpc.NewClient(
		conn,
		"datapace.AuthService",
		"Email",
		encodeIdentifyRequest,
		decodeEmailResponse,
		datapace.UserEmail{},
	).Endpoint()

	partners := kitgrpc.NewClient(
		conn,
		"datapace.AuthService",
		"Partners",
		encodePartnersRequest,
		decodePartnersResponse,
		datapace.PartnersList{},
	).Endpoint()

	return &grpcClient{
		identify: identify,
		email:    email,
		partners: partners,
	}
}

func (client grpcClient) Identify(ctx context.Context, token *datapace.Token, _ ...grpc.CallOption) (*datapace.UserID, error) {
	res, err := client.identify(ctx, identityReq{token.GetValue()})
	if err != nil {
		return nil, err
	}

	idRes := res.(identityRes)
	return &datapace.UserID{Value: idRes.id}, idRes.err
}

func (client grpcClient) Email(ctx context.Context, id *datapace.Token, _ ...grpc.CallOption) (*datapace.UserEmail, error) {
	res, err := client.email(ctx, identityReq{id.GetValue()})
	if err != nil {
		return nil, err
	}

	emailRes := res.(emailRes)
	return &datapace.UserEmail{Email: emailRes.email, ContactEmail: emailRes.contactEmail}, emailRes.err
}

func (client grpcClient) Partners(ctx context.Context, id *datapace.UserID, _ ...grpc.CallOption) (*datapace.PartnersList, error) {
	res, err := client.partners(ctx, partnersReq{id.GetValue()})
	if err != nil {
		return nil, err
	}

	partnersRes := res.(partnersRes)
	return &datapace.PartnersList{Value: partnersRes.partners}, partnersRes.err
}

func encodeIdentifyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(identityReq)
	return &datapace.Token{Value: req.token}, nil
}

func decodeIdentifyResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*datapace.UserID)
	return identityRes{res.GetValue(), nil}, nil
}

func decodeEmailResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*datapace.UserEmail)
	return emailRes{res.GetEmail(), res.GetContactEmail(), nil}, nil
}

func encodePartnersRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(partnersReq)
	return &datapace.UserID{Value: req.id}, nil
}

func decodePartnersResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*datapace.PartnersList)
	return partnersRes{res.GetValue(), nil}, nil
}
