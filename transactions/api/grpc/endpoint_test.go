package grpc_test

import (
	"context"
	"fmt"
	"monetasa"
	grpcapi "monetasa/transactions/api/grpc"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateUser(t *testing.T) {
	addr := fmt.Sprintf("localhost:%d", port)
	conn, _ := grpc.Dial(addr, grpc.WithInsecure())
	cli := grpcapi.NewClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cases := map[string]struct {
		id     string
		secret string
		key    []byte
		code   codes.Code
	}{
		"create new user": {
			id:     "new_user",
			secret: secret,
			key:    []byte(secret),
			code:   codes.OK,
		},
		"create existing user": {
			id:     id,
			secret: secret,
			key:    nil,
			code:   codes.Internal,
		},
		"create user with empty id": {
			id:     "",
			secret: secret,
			key:    nil,
			code:   codes.InvalidArgument,
		},
		"create user with empty secret": {
			id:     "other_user",
			secret: "",
			key:    nil,
			code:   codes.InvalidArgument,
		},
	}

	for desc, tc := range cases {
		key, err := cli.CreateUser(ctx, &monetasa.User{Id: tc.id, Secret: tc.secret})
		e, ok := status.FromError(err)
		assert.True(t, ok, "OK expected to be true")
		assert.Equal(t, tc.key, key.GetValue(), fmt.Sprintf("%s: expected %s got %s", desc, tc.key, key.GetValue()))
		assert.Equal(t, tc.code, e.Code(), fmt.Sprintf("%s: expected %s got %s", desc, tc.code, e.Code()))
	}
}
