package grpc

import (
	"context"

	accessproto "github.com/datapace/datapace/proto/access"
	commonproto "github.com/datapace/datapace/proto/common"
	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

var _ accessproto.AccessServiceClient = (*grpcClient)(nil)

type grpcClient struct {
	partners  endpoint.Endpoint
	potential endpoint.Endpoint
}

// NewClient returns new gRPC client instance.
func NewClient(conn *grpc.ClientConn) accessproto.AccessServiceClient {
	partners := kitgrpc.NewClient(
		conn,
		"datapace.AccessService",
		"Partners",
		encodePartnersRequest,
		decodePartnersResponse,
		accessproto.PartnersList{},
	).Endpoint()

	potential := kitgrpc.NewClient(
		conn,
		"datapace.AccessService",
		"PotentialPartners",
		encodePartnersRequest,
		decodePartnersResponse,
		accessproto.PartnersList{},
	).Endpoint()

	return &grpcClient{
		partners:  partners,
		potential: potential,
	}
}

func (client grpcClient) Partners(ctx context.Context, id *commonproto.ID, _ ...grpc.CallOption) (*accessproto.PartnersList, error) {
	res, err := client.partners(ctx, partnersReq{id.GetValue()})
	if err != nil {
		return nil, err
	}

	partnersRes := res.(partnersRes)
	return &accessproto.PartnersList{Value: partnersRes.partners}, partnersRes.err
}

func (client grpcClient) PotentialPartners(ctx context.Context, id *commonproto.ID, _ ...grpc.CallOption) (*accessproto.PartnersList, error) {
	res, err := client.potential(ctx, partnersReq{id.GetValue()})
	if err != nil {
		return nil, err
	}

	pres := res.(partnersRes)
	return &accessproto.PartnersList{Value: pres.partners}, pres.err
}

func encodePartnersRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(partnersReq)
	return &commonproto.ID{Value: req.id}, nil
}

func decodePartnersResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*accessproto.PartnersList)
	return partnersRes{res.GetValue(), nil}, nil
}
