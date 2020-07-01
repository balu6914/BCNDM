package mocks

import (
	"context"
	"sync"

	authproto "github.com/datapace/datapace/proto/auth"
	commonproto "github.com/datapace/datapace/proto/common"
	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ authproto.AuthServiceClient = (*mockAuthClient)(nil)

type mockAuthClient struct {
	users map[string]string
	mutex *sync.Mutex
}

// NewAuthClient returns auth client mock instance.
func NewAuthClient(users map[string]string) authproto.AuthServiceClient {
	return &mockAuthClient{users: users, mutex: &sync.Mutex{}}
}

func (mac *mockAuthClient) Identify(_ context.Context, token *authproto.Token, _ ...grpc.CallOption) (*commonproto.ID, error) {
	mac.mutex.Lock()
	defer mac.mutex.Unlock()

	key := token.GetValue()
	id, ok := mac.users[key]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "failed to identify user from key")
	}

	return &commonproto.ID{Value: id}, nil
}

func (mac *mockAuthClient) Email(_ context.Context, token *authproto.Token, _ ...grpc.CallOption) (*authproto.UserEmail, error) {
	return &authproto.UserEmail{Email: "", ContactEmail: ""}, nil
}

func (mac *mockAuthClient) Exists(_ context.Context, id *commonproto.ID, _ ...grpc.CallOption) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
