// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: proto/access.proto

package access

import (
	context "context"
	common "github.com/datapace/datapace/proto/common"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AccessServiceClient is the client API for AccessService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccessServiceClient interface {
	Partners(ctx context.Context, in *common.ID, opts ...grpc.CallOption) (*PartnersList, error)
	PotentialPartners(ctx context.Context, in *common.ID, opts ...grpc.CallOption) (*PartnersList, error)
}

type accessServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAccessServiceClient(cc grpc.ClientConnInterface) AccessServiceClient {
	return &accessServiceClient{cc}
}

func (c *accessServiceClient) Partners(ctx context.Context, in *common.ID, opts ...grpc.CallOption) (*PartnersList, error) {
	out := new(PartnersList)
	err := c.cc.Invoke(ctx, "/datapace.AccessService/Partners", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessServiceClient) PotentialPartners(ctx context.Context, in *common.ID, opts ...grpc.CallOption) (*PartnersList, error) {
	out := new(PartnersList)
	err := c.cc.Invoke(ctx, "/datapace.AccessService/PotentialPartners", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccessServiceServer is the server API for AccessService service.
// All implementations must embed UnimplementedAccessServiceServer
// for forward compatibility
type AccessServiceServer interface {
	Partners(context.Context, *common.ID) (*PartnersList, error)
	PotentialPartners(context.Context, *common.ID) (*PartnersList, error)
	mustEmbedUnimplementedAccessServiceServer()
}

// UnimplementedAccessServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAccessServiceServer struct {
}

func (UnimplementedAccessServiceServer) Partners(context.Context, *common.ID) (*PartnersList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Partners not implemented")
}
func (UnimplementedAccessServiceServer) PotentialPartners(context.Context, *common.ID) (*PartnersList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PotentialPartners not implemented")
}
func (UnimplementedAccessServiceServer) mustEmbedUnimplementedAccessServiceServer() {}

// UnsafeAccessServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccessServiceServer will
// result in compilation errors.
type UnsafeAccessServiceServer interface {
	mustEmbedUnimplementedAccessServiceServer()
}

func RegisterAccessServiceServer(s grpc.ServiceRegistrar, srv AccessServiceServer) {
	s.RegisterService(&AccessService_ServiceDesc, srv)
}

func _AccessService_Partners_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessServiceServer).Partners(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/datapace.AccessService/Partners",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessServiceServer).Partners(ctx, req.(*common.ID))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessService_PotentialPartners_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessServiceServer).PotentialPartners(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/datapace.AccessService/PotentialPartners",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessServiceServer).PotentialPartners(ctx, req.(*common.ID))
	}
	return interceptor(ctx, in, info, handler)
}

// AccessService_ServiceDesc is the grpc.ServiceDesc for AccessService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccessService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "datapace.AccessService",
	HandlerType: (*AccessServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Partners",
			Handler:    _AccessService_Partners_Handler,
		},
		{
			MethodName: "PotentialPartners",
			Handler:    _AccessService_PotentialPartners_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/access.proto",
}
