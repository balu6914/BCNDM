// Code generated by protoc-gen-go. DO NOT EDIT.
// source: access.proto

package datapace

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type PartnersList struct {
	Value                []string `protobuf:"bytes,1,rep,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PartnersList) Reset()         { *m = PartnersList{} }
func (m *PartnersList) String() string { return proto.CompactTextString(m) }
func (*PartnersList) ProtoMessage()    {}
func (*PartnersList) Descriptor() ([]byte, []int) {
	return fileDescriptor_a098e900d2c3a6f2, []int{0}
}

func (m *PartnersList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PartnersList.Unmarshal(m, b)
}
func (m *PartnersList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PartnersList.Marshal(b, m, deterministic)
}
func (m *PartnersList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PartnersList.Merge(m, src)
}
func (m *PartnersList) XXX_Size() int {
	return xxx_messageInfo_PartnersList.Size(m)
}
func (m *PartnersList) XXX_DiscardUnknown() {
	xxx_messageInfo_PartnersList.DiscardUnknown(m)
}

var xxx_messageInfo_PartnersList proto.InternalMessageInfo

func (m *PartnersList) GetValue() []string {
	if m != nil {
		return m.Value
	}
	return nil
}

func init() {
	proto.RegisterType((*PartnersList)(nil), "datapace.PartnersList")
}

func init() { proto.RegisterFile("access.proto", fileDescriptor_a098e900d2c3a6f2) }

var fileDescriptor_a098e900d2c3a6f2 = []byte{
	// 150 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0x4c, 0x4e, 0x4e,
	0x2d, 0x2e, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x48, 0x49, 0x2c, 0x49, 0x2c, 0x48,
	0x4c, 0x4e, 0x95, 0xe2, 0x49, 0xce, 0xcf, 0xcd, 0xcd, 0xcf, 0x83, 0x88, 0x2b, 0xa9, 0x70, 0xf1,
	0x04, 0x24, 0x16, 0x95, 0xe4, 0xa5, 0x16, 0x15, 0xfb, 0x64, 0x16, 0x97, 0x08, 0x89, 0x70, 0xb1,
	0x96, 0x25, 0xe6, 0x94, 0xa6, 0x4a, 0x30, 0x2a, 0x30, 0x6b, 0x70, 0x06, 0x41, 0x38, 0x46, 0x0d,
	0x8c, 0x5c, 0xbc, 0x8e, 0x60, 0xe3, 0x82, 0x53, 0x8b, 0xca, 0x32, 0x93, 0x53, 0x85, 0x8c, 0xb8,
	0x38, 0x60, 0xfa, 0x84, 0x78, 0xf4, 0x60, 0x86, 0xeb, 0x79, 0xba, 0x48, 0x89, 0x21, 0x78, 0xc8,
	0x26, 0x2b, 0x31, 0x08, 0x59, 0x73, 0x09, 0x06, 0xe4, 0x97, 0xa4, 0xe6, 0x95, 0x64, 0x26, 0xe6,
	0x90, 0xaa, 0x39, 0x89, 0x0d, 0xec, 0x5e, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x7c, 0x31,
	0x52, 0x25, 0xd7, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AccessServiceClient is the client API for AccessService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AccessServiceClient interface {
	Partners(ctx context.Context, in *ID, opts ...grpc.CallOption) (*PartnersList, error)
	PotentialPartners(ctx context.Context, in *ID, opts ...grpc.CallOption) (*PartnersList, error)
}

type accessServiceClient struct {
	cc *grpc.ClientConn
}

func NewAccessServiceClient(cc *grpc.ClientConn) AccessServiceClient {
	return &accessServiceClient{cc}
}

func (c *accessServiceClient) Partners(ctx context.Context, in *ID, opts ...grpc.CallOption) (*PartnersList, error) {
	out := new(PartnersList)
	err := c.cc.Invoke(ctx, "/datapace.AccessService/Partners", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessServiceClient) PotentialPartners(ctx context.Context, in *ID, opts ...grpc.CallOption) (*PartnersList, error) {
	out := new(PartnersList)
	err := c.cc.Invoke(ctx, "/datapace.AccessService/PotentialPartners", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccessServiceServer is the server API for AccessService service.
type AccessServiceServer interface {
	Partners(context.Context, *ID) (*PartnersList, error)
	PotentialPartners(context.Context, *ID) (*PartnersList, error)
}

func RegisterAccessServiceServer(s *grpc.Server, srv AccessServiceServer) {
	s.RegisterService(&_AccessService_serviceDesc, srv)
}

func _AccessService_Partners_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
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
		return srv.(AccessServiceServer).Partners(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessService_PotentialPartners_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
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
		return srv.(AccessServiceServer).PotentialPartners(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}

var _AccessService_serviceDesc = grpc.ServiceDesc{
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
	Metadata: "access.proto",
}
