package grpc

import (
	"context"
	"datapace"

	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

var _ datapace.AccessServiceClient = (*grpcClient)(nil)

type grpcClient struct {
	partners  endpoint.Endpoint
	potential endpoint.Endpoint
}

// NewClient returns new gRPC client instance.
func NewClient(conn *grpc.ClientConn) datapace.AccessServiceClient {
	partners := kitgrpc.NewClient(
		conn,
		"datapace.AccessService",
		"Partners",
		encodePartnersRequest,
		decodePartnersResponse,
		datapace.PartnersList{},
	).Endpoint()

	potential := kitgrpc.NewClient(
		conn,
		"datapace.AccessService",
		"PotentialPartners",
		encodePartnersRequest,
		decodePartnersResponse,
		datapace.PartnersList{},
	).Endpoint()

	return &grpcClient{
		partners:  partners,
		potential: potential,
	}
}

func (client grpcClient) Partners(ctx context.Context, id *datapace.ID, _ ...grpc.CallOption) (*datapace.PartnersList, error) {
	res, err := client.partners(ctx, partnersReq{id.GetValue()})
	if err != nil {
		return nil, err
	}

	partnersRes := res.(partnersRes)
	return &datapace.PartnersList{Value: partnersRes.partners}, partnersRes.err
}

func (client grpcClient) PotentialPartners(ctx context.Context, id *datapace.ID, _ ...grpc.CallOption) (*datapace.PartnersList, error) {
	res, err := client.potential(ctx, partnersReq{id.GetValue()})
	if err != nil {
		return nil, err
	}

	pres := res.(partnersRes)
	return &datapace.PartnersList{Value: pres.partners}, pres.err
}

func encodePartnersRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(partnersReq)
	return &datapace.ID{Value: req.id}, nil
}

func decodePartnersResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*datapace.PartnersList)
	return partnersRes{res.GetValue(), nil}, nil
}
