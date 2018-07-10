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
		id   string
		key  []byte
		code codes.Code
	}{
		"create new user": {
			id:   "5281b83afbb7f35cb62d0835",
			key:  []byte(secret),
			code: codes.OK,
		},
		"create existing user": {
			id:   id,
			key:  nil,
			code: codes.Internal,
		},
		"create user with empty id": {
			id:   "",
			key:  nil,
			code: codes.InvalidArgument,
		},
	}

	for desc, tc := range cases {
		_, err := cli.CreateUser(ctx, &monetasa.ID{Value: tc.id})
		e, ok := status.FromError(err)
		assert.True(t, ok, "OK expected to be true")
		assert.Equal(t, tc.code, e.Code(), fmt.Sprintf("%s: expected %s got %s", desc, tc.code, e.Code()))
	}
}
