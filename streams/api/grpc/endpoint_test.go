package grpc_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	commonproto "github.com/datapace/datapace/proto/common"
	streamsproto "github.com/datapace/datapace/proto/streams"
	"github.com/datapace/datapace/streams"
	grpcapi "github.com/datapace/datapace/streams/api/grpc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/mgo.v2/bson"
)

var stream = streams.Stream{
	Owner:      owner,
	ID:         bson.NewObjectId().Hex(),
	URL:        "http://test.com",
	Price:      100,
	Visibility: streams.Private,
}

func TestOne(t *testing.T) {
	id, err := svc.AddStream(stream)
	require.Nil(t, err, fmt.Sprintf("received unexpected error: %s", err))

	conn := createConn(t)

	cli := grpcapi.NewClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cases := map[string]struct {
		id     string
		stream *streamsproto.Stream
		err    error
	}{
		"get existing stream": {
			id: id,
			stream: &streamsproto.Stream{
				Owner:      stream.Owner,
				Id:         id,
				Url:        stream.URL,
				Price:      stream.Price,
				Visibility: string(stream.Visibility),
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
		stream, err := cli.One(ctx, &commonproto.ID{Value: tc.id})
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
