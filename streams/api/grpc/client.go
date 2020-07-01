package grpc

import (
	"context"

	commonproto "github.com/datapace/datapace/proto/common"
	streamsproto "github.com/datapace/datapace/proto/streams"
	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

var _ streamsproto.StreamsServiceClient = (*grpcClient)(nil)

type grpcClient struct {
	one endpoint.Endpoint
}

// NewClient returns new gRPC client instance.
func NewClient(conn *grpc.ClientConn) streamsproto.StreamsServiceClient {
	return grpcClient{
		one: kitgrpc.NewClient(
			conn,
			"datapace.StreamsService",
			"One",
			encodeOneRequest,
			decodeOneResponse,
			streamsproto.Stream{},
		).Endpoint(),
	}
}

func (client grpcClient) One(ctx context.Context, id *commonproto.ID, _ ...grpc.CallOption) (*streamsproto.Stream, error) {
	res, err := client.one(ctx, oneReq{id: id.GetValue()})
	if err != nil {
		return nil, err
	}

	sr := res.(oneRes)
	stream := streamsproto.Stream{
		Id:       sr.id,
		Name:     sr.name,
		Owner:    sr.owner,
		Url:      sr.url,
		Price:    sr.price,
		External: sr.external,
		Project:  sr.project,
		Dataset:  sr.dataset,
		Table:    sr.table,
		Fields:   sr.fields,
	}

	return &stream, sr.err
}

func encodeOneRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(oneReq)
	return &commonproto.ID{Value: req.id}, nil
}

func decodeOneResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*streamsproto.Stream)
	stream := oneRes{
		id:       res.GetId(),
		name:     res.GetName(),
		owner:    res.GetOwner(),
		url:      res.GetUrl(),
		price:    res.GetPrice(),
		external: res.GetExternal(),
		project:  res.GetProject(),
		dataset:  res.GetDataset(),
		table:    res.GetTable(),
		fields:   res.GetFields(),
	}

	return stream, nil
}
