package mocks

import (
	"context"
	"monetasa"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ monetasa.AuthServiceClient = (*mockAuthClient)(nil)

type mockAuthClient struct {
	users map[string]string
	mutex *sync.Mutex
}

// NewAuthClient returns auth client mock instance.
func NewAuthClient(users map[string]string) monetasa.AuthServiceClient {
	return mockAuthClient{users: users, mutex: &sync.Mutex{}}
}

func (mac mockAuthClient) Identify(_ context.Context, token *monetasa.Token, _ ...grpc.CallOption) (*monetasa.UserID, error) {
	mac.mutex.Lock()
	defer mac.mutex.Unlock()

	key := token.GetValue()
	id, ok := mac.users[key]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "failed to identify user from key")
	}

	return &monetasa.UserID{Value: id}, nil
}
