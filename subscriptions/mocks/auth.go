package mocks

import (
	"context"

	authproto "github.com/datapace/datapace/proto/auth"
	commonproto "github.com/datapace/datapace/proto/common"

	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ authproto.AuthServiceClient = (*authClientMock)(nil)

type authClientMock struct {
	tokens map[string]string
	emails map[string]string
}

// NewAuthClient creates mock of users service.
func NewAuthClient(tokens map[string]string, emails map[string]string) authproto.AuthServiceClient {
	return &authClientMock{tokens, emails}
}

func (svc authClientMock) Identify(_ context.Context, in *authproto.Token, opts ...grpc.CallOption) (*commonproto.ID, error) {
	if id, ok := svc.tokens[in.Value]; ok {
		return &commonproto.ID{Value: id}, nil
	}

	return nil, status.Error(codes.Unauthenticated, "unauthenticated")
}

func (svc authClientMock) Email(_ context.Context, in *authproto.Token, opts ...grpc.CallOption) (*authproto.UserEmail, error) {
	if id, ok := svc.emails[in.Value]; ok {
		return &authproto.UserEmail{Email: id, ContactEmail: ""}, nil
	}

	return nil, status.Error(codes.Unauthenticated, "unauthenticated")
}

func (svc authClientMock) Exists(_ context.Context, id *commonproto.ID, opts ...grpc.CallOption) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (svc authClientMock) Authorize(ctx context.Context, ar *authproto.AuthRequest, opts ...grpc.CallOption) (*commonproto.ID, error) {
	if id, ok := svc.tokens[ar.Token]; ok {
		return &commonproto.ID{Value: id}, nil
	}

	return nil, status.Error(codes.Unauthenticated, "unauthenticated")
}
