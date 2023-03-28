// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package offers

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// OffersServiceClient is the client API for OffersService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OffersServiceClient interface {
	GetOfferPrice(ctx context.Context, in *GetOfferPriceRequest, opts ...grpc.CallOption) (*GetOfferPriceResponse, error)
}

type offersServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOffersServiceClient(cc grpc.ClientConnInterface) OffersServiceClient {
	return &offersServiceClient{cc}
}

func (c *offersServiceClient) GetOfferPrice(ctx context.Context, in *GetOfferPriceRequest, opts ...grpc.CallOption) (*GetOfferPriceResponse, error) {
	out := new(GetOfferPriceResponse)
	err := c.cc.Invoke(ctx, "/offers.OffersService/GetOfferPrice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OffersServiceServer is the server API for OffersService service.
// All implementations must embed UnimplementedOffersServiceServer
// for forward compatibility
type OffersServiceServer interface {
	GetOfferPrice(context.Context, *GetOfferPriceRequest) (*GetOfferPriceResponse, error)
	mustEmbedUnimplementedOffersServiceServer()
}

// UnimplementedOffersServiceServer must be embedded to have forward compatible implementations.
type UnimplementedOffersServiceServer struct {
}

func (*UnimplementedOffersServiceServer) GetOfferPrice(context.Context, *GetOfferPriceRequest) (*GetOfferPriceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOfferPrice not implemented")
}
func (*UnimplementedOffersServiceServer) mustEmbedUnimplementedOffersServiceServer() {}

func RegisterOffersServiceServer(s *grpc.Server, srv OffersServiceServer) {
	s.RegisterService(&_OffersService_serviceDesc, srv)
}

func _OffersService_GetOfferPrice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOfferPriceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OffersServiceServer).GetOfferPrice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/offers.OffersService/GetOfferPrice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OffersServiceServer).GetOfferPrice(ctx, req.(*GetOfferPriceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _OffersService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "offers.OffersService",
	HandlerType: (*OffersServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOfferPrice",
			Handler:    _OffersService_GetOfferPrice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/offers.proto",
}
