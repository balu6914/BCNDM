package grpc

import (
	"context"

	"github.com/datapace/datapace/auth"
	"github.com/datapace/datapace/dproxy"
	"github.com/datapace/datapace/dproxy/persistence"
	dproxyproto "github.com/datapace/datapace/proto/dproxy"
	"google.golang.org/protobuf/types/known/timestamppb"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ dproxyproto.DproxyServiceServer = (*grpcServer)(nil)

type grpcServer struct {
	dproxyproto.UnimplementedDproxyServiceServer
	list kitgrpc.Handler
}

// NewServer instantiates new Auth gRPC server.
func NewServer(svc dproxy.Service) dproxyproto.DproxyServiceServer {
	list := kitgrpc.NewServer(
		listEndpoint(svc),
		decodeListRequest,
		encodeListResponse,
	)

	return &grpcServer{
		list: list,
	}
}

func (s *grpcServer) List(ctx context.Context, req *dproxyproto.ListRequest) (*dproxyproto.ListResponse, error) {
	_, res, err := s.list.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*dproxyproto.ListResponse), nil
}

func decodeListRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*dproxyproto.ListRequest)
	q := decodeQuery(req)
	return q, nil
}

func encodeListResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(listResponse)
	var results []*dproxyproto.AccessLog
	for _, al := range res.events {
		results = append(results, encodeAccessLog(al))
	}
	resp := dproxyproto.ListResponse{Page: results}
	return &resp, nil
}

func encodeError(err error) error {
	switch err {
	case nil:
		return nil
	case auth.ErrMalformedEntity:
		return status.Error(codes.InvalidArgument, "received invalid request")
	case auth.ErrUnauthorizedAccess:
		return status.Error(codes.Unauthenticated, "failed to identify user from key")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}

func decodeQuery(pageReq *dproxyproto.ListRequest) (q listRequest) {
	cursorReq := pageReq.Cursor
	cursor := persistence.Event{}
	if cursorReq != nil {
		cursor.SubID = cursorReq.SubId
		cursor.Time = cursorReq.Time.AsTime()
	}
	sortReq := pageReq.Sort
	sort := persistence.Sort{}
	if sortReq != nil {
		sort.By = persistence.SortBy(sortReq.By)
		sort.Order = persistence.SortOrder(sortReq.Order)
	}
	pq := persistence.Query{
		Limit:  pageReq.Limit,
		Cursor: cursor,
		Sort:   sort,
	}
	q = listRequest{query: pq}
	return
}

func encodeAccessLog(access persistence.Event) (resp *dproxyproto.AccessLog) {
	resp = &dproxyproto.AccessLog{

		SubId: access.SubID,
		Time:  timestamppb.New(access.Time),
	}
	return
}
