package grpc_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/datapace/datapace/auth"
	grpcapi "github.com/datapace/datapace/auth/api/grpc"
	authproto "github.com/datapace/datapace/proto/auth"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	email   = "john.doe@email.com"
	invalid = "invalid"
)

var user = auth.User{
	ID:        email,
	Email:     email,
	Password:  "password",
	FirstName: "Joe",
	LastName:  "Doe",
}

var admin = auth.User{
	ID:        "admin@example.com",
	Email:     "admin@example.com",
	Password:  "password",
	FirstName: "Joe",
	LastName:  "Doe",
	Company:   "company",
	Address:   "address",
	Phone:     "+1234567890",
	Roles:     []string{"admin"},
}

func TestIdentify(t *testing.T) {
	svc.Register(k, user)

	authAddr := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.Dial(authAddr, grpc.WithInsecure())
	require.Nil(t, err, "unexpected error dialing GRPC: %s", err)
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
			err:   status.Error(codes.InvalidArgument, "received invalid request"),
		},
		"identify user with invalid token": {
			token: invalid,
			id:    "",
			err:   status.Error(codes.Unauthenticated, "failed to identify user from key"),
		},
	}

	for desc, tc := range cases {
		id, err := client.Identify(ctx, &authproto.Token{Value: tc.token})
		assert.Equal(t, tc.id, id.GetValue(), fmt.Sprintf("%s: expected %s got %s", desc, tc.id, id.GetValue()))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}
}

func TestEmail(t *testing.T) {
	authAddr := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.Dial(authAddr, grpc.WithInsecure())
	require.Nil(t, err, "unexpected error dialing GRPC: %s", err)
	client := grpcapi.NewClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cases := map[string]struct {
		token string
		email string
		err   error
	}{
		"get an email": {
			token: user.Email,
			email: user.Email,
			err:   nil,
		},
		"get an email with an empty id": {
			token: "",
			email: "",
			err:   status.Error(codes.InvalidArgument, "received invalid request"),
		},
		"get an email with an invalid id": {
			token: invalid,
			email: "",
			err:   status.Error(codes.Unauthenticated, "failed to identify user from key"),
		},
	}

	for desc, tc := range cases {
		email, err := client.Email(ctx, &authproto.Token{Value: tc.token})
		assert.Equal(t, tc.email, email.GetEmail(), fmt.Sprintf("%s: expected %s got %s", desc, tc.email, email.GetEmail()))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}
}
