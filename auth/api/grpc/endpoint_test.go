package grpc_test

import (
	"context"
	"fmt"
	"monetasa"
	"monetasa/auth"
	grpcapi "monetasa/auth/api/grpc"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	email   = "john.doe@email.com"
	invalid = "invalid"
)

var user = auth.User{
	ID:       email,
	Email:    email,
	Password: "pass",
	Name:     "John Doe",
}

func TestIdentify(t *testing.T) {
	svc.Register(user)

	authAddr := fmt.Sprintf("localhost:%d", port)
	conn, _ := grpc.Dial(authAddr, grpc.WithInsecure())
	client := grpcapi.NewClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cases := map[string]struct {
		token string
		id    string
		err   error
	}{
		"identify user": {
			token: user.Email,
			id:    user.Email,
			err:   nil,
		},
		"identify user with empty token": {
			token: "",
			id:    "",
			err:   status.Error(codes.InvalidArgument, "received invalid key request"),
		},
		"identify user with invalid token": {
			token: invalid,
			id:    "",
			err:   status.Error(codes.Unauthenticated, "failed to identify user from key"),
		},
	}

	for desc, tc := range cases {
		id, err := client.Identify(ctx, &monetasa.Token{Value: tc.token})
		assert.Equal(t, tc.id, id.GetValue(), fmt.Sprintf("%s: expected %s got %s", desc, tc.id, id.GetValue()))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}
}
