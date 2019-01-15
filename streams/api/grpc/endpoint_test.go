package grpc_test

import (
	"context"
	"fmt"
	"datapace"
	"datapace/streams"
	grpcapi "datapace/streams/api/grpc"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/mgo.v2/bson"
)

var stream = streams.Stream{
	Owner: owner,
	ID:    bson.NewObjectId(),
	URL:   "http://test.com",
	Price: 100,
}

func TestOne(t *testing.T) {
	svc.AddStream(stream)

	conn := createConn(t)

	cli := grpcapi.NewClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cases := map[string]struct {
		id     string
		stream *datapace.Stream
		err    error
	}{
		"get existing stream": {
			id: stream.ID.Hex(),
			stream: &datapace.Stream{
				Owner: stream.Owner,
				Id:    stream.ID.Hex(),
				Url:   stream.URL,
				Price: stream.Price,
			},
			err: nil,
		},
		"get non-existent stream": {
			id:     "non-existent",
			stream: nil,
			err:    status.Error(codes.NotFound, "stream doesn't exist"),
		},
	}

	for desc, tc := range cases {
		stream, err := cli.One(ctx, &datapace.ID{Value: tc.id})
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		assert.Equal(t, tc.stream, stream, fmt.Sprintf("%s: expected %v got %v", desc, tc.stream, stream))
	}
}

func createConn(t *testing.T) *grpc.ClientConn {
	addr := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	return conn
}
