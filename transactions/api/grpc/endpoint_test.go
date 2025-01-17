package grpc_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	commonproto "github.com/datapace/datapace/proto/common"
	transactionsproto "github.com/datapace/datapace/proto/transactions"
	grpcapi "github.com/datapace/datapace/transactions/api/grpc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const streamID = "stream_id"

func TestCreateUser(t *testing.T) {
	conn := createConn(t)

	cli := grpcapi.NewClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cases := []struct {
		desc   string
		id     string
		key    []byte
		status codes.Code
	}{
		{
			desc:   "create new user",
			id:     "5281b83afbb7f35cb62d0836",
			key:    []byte(secret),
			status: codes.OK,
		},
		{
			desc:   "create existing user",
			id:     id1,
			key:    nil,
			status: codes.Internal,
		},
		{
			desc:   "create user with empty id",
			id:     "",
			key:    nil,
			status: codes.InvalidArgument,
		},
	}

	for _, tc := range cases {
		_, err := cli.CreateUser(ctx, &commonproto.ID{Value: tc.id})
		e, ok := status.FromError(err)
		assert.True(t, ok, "OK expected to be true")
		assert.Equal(t, tc.status, e.Code(), fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.status, e.Code()))
	}
}

func TestTransfer(t *testing.T) {
	conn := createConn(t)

	cli := grpcapi.NewClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cases := []struct {
		desc   string
		from   string
		to     string
		value  uint64
		status codes.Code
	}{
		{
			desc:   "transfer money from one account to another",
			from:   id1,
			to:     id2,
			value:  balance,
			status: codes.OK,
		},
		{
			desc:   "transfer money with invalid request",
			from:   "",
			to:     id2,
			value:  balance,
			status: codes.InvalidArgument,
		},
		{
			desc:   "transfer money from non-existent account",
			from:   "non-existent",
			to:     id2,
			value:  balance,
			status: codes.Internal,
		},
	}

	for _, tc := range cases {
		req := transactionsproto.TransferData{
			StreamID: streamID,
			From:     tc.from,
			To:       tc.to,
			Value:    tc.value,
		}

		_, err := cli.Transfer(ctx, &req)
		e, ok := status.FromError(err)
		assert.True(t, ok, "OK expected to be true")
		assert.Equal(t, tc.status, e.Code(), fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.status, e.Code()))
	}
}

func createConn(t *testing.T) *grpc.ClientConn {
	addr := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	return conn
}
