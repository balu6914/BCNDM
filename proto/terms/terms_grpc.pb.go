// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: proto/terms.proto

package terms

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TermsServiceClient is the client API for TermsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TermsServiceClient interface {
	CreateTerms(ctx context.Context, in *Terms, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type termsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTermsServiceClient(cc grpc.ClientConnInterface) TermsServiceClient {
	return &termsServiceClient{cc}
}

func (c *termsServiceClient) CreateTerms(ctx context.Context, in *Terms, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/datapace.TermsService/CreateTerms", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TermsServiceServer is the server API for TermsService service.
// All implementations must embed UnimplementedTermsServiceServer
// for forward compatibility
type TermsServiceServer interface {
	CreateTerms(context.Context, *Terms) (*emptypb.Empty, error)
	mustEmbedUnimplementedTermsServiceServer()
}

// UnimplementedTermsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTermsServiceServer struct {
}

func (UnimplementedTermsServiceServer) CreateTerms(context.Context, *Terms) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTerms not implemented")
}
func (UnimplementedTermsServiceServer) mustEmbedUnimplementedTermsServiceServer() {}

// UnsafeTermsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TermsServiceServer will
// result in compilation errors.
type UnsafeTermsServiceServer interface {
	mustEmbedUnimplementedTermsServiceServer()
}

func RegisterTermsServiceServer(s grpc.ServiceRegistrar, srv TermsServiceServer) {
	s.RegisterService(&TermsService_ServiceDesc, srv)
}

func _TermsService_CreateTerms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Terms)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TermsServiceServer).CreateTerms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/datapace.TermsService/CreateTerms",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TermsServiceServer).CreateTerms(ctx, req.(*Terms))
	}
	return interceptor(ctx, in, info, handler)
}

// TermsService_ServiceDesc is the grpc.ServiceDesc for TermsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TermsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "datapace.TermsService",
	HandlerType: (*TermsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateTerms",
			Handler:    _TermsService_CreateTerms_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/terms.proto",
}
