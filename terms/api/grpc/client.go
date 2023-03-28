package grpc

import (
	"context"
	termsproto "github.com/datapace/datapace/proto/terms"
	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

var _ termsproto.TermsServiceClient = (*grpcClient)(nil)

type grpcClient struct {
	createTerms endpoint.Endpoint
}

// NewClient returns new gRPC client instance.
func NewClient(conn *grpc.ClientConn) termsproto.TermsServiceClient {
	createTerms := kitgrpc.NewClient(
		conn,
		"datapace.TermsService",
		"CreateTerms",
		encodeTermsReq,
		decodeCreateRes,
		empty.Empty{},
	).Endpoint()

	return &grpcClient{
		createTerms: createTerms,
	}
}

func (client grpcClient) CreateTerms(ctx context.Context, terms *termsproto.Terms, _ ...grpc.CallOption) (*empty.Empty, error) {
	req := termsReq{
		termsUrl: terms.GetUrl(),
		streamId: terms.GetStreamId(),
	}
	res, err := client.createTerms(ctx, req)
	if err != nil {
		return nil, err
	}

	cr := res.(createRes)
	return &empty.Empty{}, cr.err
}

func encodeTermsReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(termsReq)
	return &termsproto.Terms{
		Url:      req.termsUrl,
		StreamId: req.streamId,
	}, nil
}

func decodeCreateRes(_ context.Context, grpcRes interface{}) (interface{}, error) {
	return createRes{}, nil
}
