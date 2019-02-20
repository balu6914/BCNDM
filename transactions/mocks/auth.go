package mocks

import (
	"context"
	"datapace"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ datapace.AuthServiceClient = (*mockAuthClient)(nil)

type mockAuthClient struct {
	users map[string]string
	mutex *sync.Mutex
}

// NewAuthClient returns auth client mock instance.
func NewAuthClient(users map[string]string) datapace.AuthServiceClient {
	return &mockAuthClient{users: users, mutex: &sync.Mutex{}}
}

func (mac *mockAuthClient) Identify(_ context.Context, token *datapace.Token, _ ...grpc.CallOption) (*datapace.UserID, error) {
	mac.mutex.Lock()
	defer mac.mutex.Unlock()

	key := token.GetValue()
	id, ok := mac.users[key]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "failed to identify user from key")
	}

	return &datapace.UserID{Value: id}, nil
}

func (mac *mockAuthClient) Email(_ context.Context, token *datapace.Token, _ ...grpc.CallOption) (*datapace.UserEmail, error) {
	return &datapace.UserEmail{Email: "", ContactEmail: ""}, nil
}

func (mac *mockAuthClient) Partners(_ context.Context, id *datapace.UserID, _ ...grpc.CallOption) (*datapace.PartnersList, error) {
	return &datapace.PartnersList{Value: []string{}}, nil
}
