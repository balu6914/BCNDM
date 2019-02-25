package mocks

import (
	"context"

	"datapace"

	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ datapace.AuthServiceClient = (*authClientMock)(nil)

type authClientMock struct {
	tokens map[string]string
	emails map[string]string
}

// NewAuthClient creates mock of users service.
func NewAuthClient(tokens map[string]string, emails map[string]string) datapace.AuthServiceClient {
	return &authClientMock{tokens, emails}
}

func (svc authClientMock) Identify(_ context.Context, in *datapace.Token, opts ...grpc.CallOption) (*datapace.ID, error) {
	if id, ok := svc.tokens[in.Value]; ok {
		return &datapace.ID{Value: id}, nil
	}

	return nil, status.Error(codes.Unauthenticated, "unauthenticated")
}

func (svc authClientMock) Email(_ context.Context, in *datapace.Token, opts ...grpc.CallOption) (*datapace.UserEmail, error) {
	if id, ok := svc.emails[in.Value]; ok {
		return &datapace.UserEmail{Email: id, ContactEmail: ""}, nil
	}

	return nil, status.Error(codes.Unauthenticated, "unauthenticated")
}

func (svc authClientMock) Exists(_ context.Context, id *datapace.ID, opts ...grpc.CallOption) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
