package grpc

import (
	"context"

	authproto "github.com/datapace/datapace/proto/auth"
	commonproto "github.com/datapace/datapace/proto/common"

	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

var _ authproto.AuthServiceClient = (*grpcClient)(nil)

type grpcClient struct {
	identify  endpoint.Endpoint
	email     endpoint.Endpoint
	exists    endpoint.Endpoint
	authorize endpoint.Endpoint
}

// NewClient returns new gRPC client instance.
func NewClient(conn *grpc.ClientConn) authproto.AuthServiceClient {
	identify := kitgrpc.NewClient(
		conn,
		"datapace.AuthService",
		"Identify",
		encodeIdentifyRequest,
		decodeIdentifyResponse,
		commonproto.ID{},
	).Endpoint()

	email := kitgrpc.NewClient(
		conn,
		"datapace.AuthService",
		"Email",
		encodeIdentifyRequest,
		decodeEmailResponse,
		authproto.UserEmail{},
	).Endpoint()

	exists := kitgrpc.NewClient(
		conn,
		"datapace.AuthService",
		"Exists",
		encodeExistsRequest,
		decodeExistsResponse,
		empty.Empty{},
	).Endpoint()

	authorize := kitgrpc.NewClient(
		conn,
		"datapace.AuthService",
		"Authorize",
		encodeAuthorizeRequest,
		decodeIdentifyResponse,
		commonproto.ID{},
	).Endpoint()

	return &grpcClient{
		identify:  identify,
		email:     email,
		exists:    exists,
		authorize: authorize,
	}
}

func (client grpcClient) Identify(ctx context.Context, token *authproto.Token, _ ...grpc.CallOption) (*commonproto.ID, error) {
	res, err := client.identify(ctx, identityReq{token.GetValue()})
	if err != nil {
		return nil, err
	}

	idRes := res.(identityRes)
	return &commonproto.ID{Value: idRes.id}, idRes.err
}

func (client grpcClient) Email(ctx context.Context, token *authproto.Token, _ ...grpc.CallOption) (*authproto.UserEmail, error) {
	res, err := client.email(ctx, identityReq{token.GetValue()})
	if err != nil {
		return nil, err
	}

	emailRes := res.(emailRes)
	return &authproto.UserEmail{Email: emailRes.email, ContactEmail: emailRes.contactEmail}, emailRes.err
}

func (client grpcClient) Exists(ctx context.Context, id *commonproto.ID, _ ...grpc.CallOption) (*empty.Empty, error) {
	res, err := client.exists(ctx, existsReq{id.GetValue()})
	if err != nil {
		return nil, err
	}

	existsRes := res.(existsRes)
	return &empty.Empty{}, existsRes.err
}

func (client grpcClient) Authorize(ctx context.Context, ar *authproto.AuthRequest, _ ...grpc.CallOption) (*commonproto.ID, error) {
	req := authReq{
		token:        ar.Token,
		action:       ar.Action,
		resourceType: ar.Type,
		attributes:   ar.Attributes,
	}
	res, err := client.authorize(ctx, req)
	if err != nil {
		return nil, err
	}

	idRes := res.(identityRes)
	return &commonproto.ID{Value: idRes.id}, idRes.err
}

func encodeIdentifyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(identityReq)
	return &authproto.Token{Value: req.token}, nil
}

func decodeIdentifyResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*commonproto.ID)
	return identityRes{res.GetValue(), nil}, nil
}

func decodeEmailResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*authproto.UserEmail)
	return emailRes{res.GetEmail(), res.GetContactEmail(), nil}, nil
}

func encodeExistsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(existsReq)
	return &commonproto.ID{Value: req.id}, nil
}

func decodeExistsResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	return existsRes{nil}, nil
}

func encodeAuthorizeRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(authReq)
	ret := &authproto.AuthRequest{
		Token:      req.token,
		Action:     req.action,
		Type:       req.resourceType,
		Attributes: req.attributes,
	}
	return ret, nil
}
