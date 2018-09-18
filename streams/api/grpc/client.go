package grpc

import (
	"context"
	"monetasa"

	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

var _ monetasa.StreamsServiceClient = (*grpcClient)(nil)

type grpcClient struct {
	one endpoint.Endpoint
}

// NewClient returns new gRPC client instance.
func NewClient(conn *grpc.ClientConn) monetasa.StreamsServiceClient {
	return grpcClient{
		one: kitgrpc.NewClient(
			conn,
			"monetasa.StreamsService",
			"One",
			encodeOneRequest,
			decodeOneResponse,
			monetasa.Stream{},
		).Endpoint(),
	}
}

func (client grpcClient) One(ctx context.Context, id *monetasa.ID, _ ...grpc.CallOption) (*monetasa.Stream, error) {
	res, err := client.one(ctx, oneReq{id: id.GetValue()})
	if err != nil {
		return nil, err
	}

	sr := res.(oneRes)
	stream := monetasa.Stream{
		Id:      sr.id,
		Name:    sr.name,
		Owner:   sr.owner,
		Url:     sr.url,
		Price:   sr.price,
		Bq:      sr.bq,
		Project: sr.project,
		Dataset: sr.dataset,
		Table:   sr.table,
		Fields:  sr.fields,
	}

	return &stream, sr.err
}

func encodeOneRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(oneReq)
	return &monetasa.ID{Value: req.id}, nil
}

func decodeOneResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*monetasa.Stream)
	stream := oneRes{
		id:      res.GetId(),
		name:    res.GetName(),
		owner:   res.GetOwner(),
		url:     res.GetUrl(),
		price:   res.GetPrice(),
		bq:      res.GetBq(),
		project: res.GetProject(),
		dataset: res.GetDataset(),
		table:   res.GetTable(),
		fields:  res.GetFields(),
	}

	return stream, nil
}
