package mocks

import (
	"context"

	"monetasa"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ monetasa.AuthServiceClient = (*authClientMock)(nil)

type authClientMock struct {
	tokens map[string]string
}

// NewAuthClient creates mock of users service.
func NewAuthClient(tokens map[string]string) monetasa.AuthServiceClient {
	return &authClientMock{tokens}
}

func (svc authClientMock) Identify(_ context.Context, in *monetasa.Token, opts ...grpc.CallOption) (*monetasa.UserID, error) {
	if id, ok := svc.tokens[in.Value]; ok {
		return &monetasa.UserID{Value: id}, nil
	}

	return nil, status.Error(codes.Unauthenticated, "unauthenticated")
}
