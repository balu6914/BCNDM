package grpc

import (
	"context"

	"github.com/datapace/datapace/dproxy/persistence"
	dproxyproto "github.com/datapace/datapace/proto/dproxy"
	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ dproxyproto.DproxyServiceClient = (*grpcClient)(nil)

type grpcClient struct {
	list endpoint.Endpoint
}

// NewClient returns new gRPC client instance.
func NewClient(conn *grpc.ClientConn) dproxyproto.DproxyServiceClient {
	return grpcClient{
		list: kitgrpc.NewClient(
			conn,
			"datapace.DproxyService",
			"List",
			encodeListRequest,
			decodeListResponse,
			dproxyproto.ListResponse{},
		).Endpoint(),
	}
}

func (client grpcClient) List(ctx context.Context, lr *dproxyproto.ListRequest, _ ...grpc.CallOption) (*dproxyproto.ListResponse, error) {
	res, err := client.list(ctx, decodeQuery(lr))
	if err != nil {
		return nil, err
	}

	listRes := res.(listResponse)
	listResp := dproxyproto.ListResponse{}
	return &listResp, listRes.err
}

func encodeListRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(listRequest)
	s := encodeSort(req.query.Sort)
	c, err := encodeCursor(req.query.Cursor)
	if err != nil {
		return nil, err
	}
	return &dproxyproto.ListRequest{
		Limit:  req.query.Limit,
		Cursor: &c,
		Sort:   &s,
	}, nil
}

func encodeSortOrder(src persistence.SortOrder) (dst dproxyproto.SortOrder) {
	switch src {
	case persistence.SortOrderAsc:
		dst = dproxyproto.SortOrder_ASC
	case persistence.SortOrderDesc:
		dst = dproxyproto.SortOrder_DESC
	}
	return
}

func encodeSortType(src persistence.SortBy) (dst dproxyproto.SortBy) {
	switch src {
	case persistence.SortByDate:
		dst = dproxyproto.SortBy_DATE
	}
	return
}

func encodeSort(src persistence.Sort) (dst dproxyproto.Sort) {

	dst = dproxyproto.Sort{
		Order: encodeSortOrder(src.Order),
		By:    encodeSortType(src.By),
	}
	return
}

func encodeCursor(src persistence.Event) (dst dproxyproto.AccessLog, err error) {
	dst = dproxyproto.AccessLog{
		SubId: src.SubID,
		Time:  timestamppb.New(src.Time),
	}
	return
}

func decodeListResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*dproxyproto.ListResponse)
	var evs []persistence.Event
	for _, p := range res.Page {
		evs = append(evs, decodeAccessLog(p))
	}
	return evs, nil
}

func decodeAccessLog(access *dproxyproto.AccessLog) (resp persistence.Event) {
	resp = persistence.Event{
		SubID: access.SubId,
		Time:  access.Time.AsTime(),
	}
	return
}
