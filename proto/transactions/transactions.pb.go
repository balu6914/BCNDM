// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/transactions.proto

package transactions

import (
	fmt "fmt"
	_ "github.com/datapace/datapace/proto/common"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/empty"
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

type TransferData struct {
	StreamID             string   `protobuf:"bytes,1,opt,name=streamID,proto3" json:"streamID,omitempty"`
	From                 string   `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	To                   string   `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`
	Value                uint64   `protobuf:"varint,4,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TransferData) Reset()         { *m = TransferData{} }
func (m *TransferData) String() string { return proto.CompactTextString(m) }
func (*TransferData) ProtoMessage()    {}
func (*TransferData) Descriptor() ([]byte, []int) {
	return fileDescriptor_1edb37854b34f19b, []int{0}
}

func (m *TransferData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransferData.Unmarshal(m, b)
}
func (m *TransferData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransferData.Marshal(b, m, deterministic)
}
func (m *TransferData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransferData.Merge(m, src)
}
func (m *TransferData) XXX_Size() int {
	return xxx_messageInfo_TransferData.Size(m)
}
func (m *TransferData) XXX_DiscardUnknown() {
	xxx_messageInfo_TransferData.DiscardUnknown(m)
}

var xxx_messageInfo_TransferData proto.InternalMessageInfo

func (m *TransferData) GetStreamID() string {
	if m != nil {
		return m.StreamID
	}
	return ""
}

func (m *TransferData) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

func (m *TransferData) GetTo() string {
	if m != nil {
		return m.To
	}
	return ""
}

func (m *TransferData) GetValue() uint64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func init() {
	proto.RegisterType((*TransferData)(nil), "datapace.TransferData")
}

func init() {
	proto.RegisterFile("proto/transactions.proto", fileDescriptor_1edb37854b34f19b)
}

var fileDescriptor_1edb37854b34f19b = []byte{
	// 254 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0x4d, 0x4b, 0xc3, 0x40,
	0x10, 0x86, 0x4d, 0x8c, 0x12, 0x97, 0xe2, 0x61, 0x94, 0x12, 0xe2, 0xa5, 0xf4, 0xd4, 0xd3, 0x06,
	0x3f, 0x8e, 0x9e, 0x34, 0x1e, 0x7a, 0xad, 0xf5, 0xe2, 0x6d, 0x92, 0x4e, 0x62, 0xa0, 0x9b, 0x09,
	0x9b, 0x49, 0xc1, 0x9f, 0xe0, 0xbf, 0x96, 0xee, 0x12, 0x5b, 0x90, 0xde, 0x66, 0x1e, 0xf6, 0x1d,
	0xde, 0x67, 0x55, 0xd2, 0x59, 0x16, 0xce, 0xc4, 0x62, 0xdb, 0x63, 0x29, 0x0d, 0xb7, 0xbd, 0x76,
	0x08, 0xe2, 0x0d, 0x0a, 0x76, 0x58, 0x52, 0x7a, 0x57, 0x33, 0xd7, 0x5b, 0xca, 0x1c, 0x2f, 0x86,
	0x2a, 0x23, 0xd3, 0xc9, 0xb7, 0x7f, 0x96, 0x82, 0x3f, 0x50, 0xb2, 0x31, 0xdc, 0x7a, 0x36, 0xdf,
	0xa8, 0xc9, 0x7a, 0x7f, 0xb0, 0x22, 0x9b, 0xa3, 0x20, 0xa4, 0x2a, 0xee, 0xc5, 0x12, 0x9a, 0x65,
	0x9e, 0x04, 0xb3, 0x60, 0x71, 0xb5, 0xfa, 0xdb, 0x01, 0x54, 0x54, 0x59, 0x36, 0x49, 0xe8, 0xb8,
	0x9b, 0xe1, 0x5a, 0x85, 0xc2, 0xc9, 0xb9, 0x23, 0xa1, 0x30, 0xdc, 0xaa, 0x8b, 0x1d, 0x6e, 0x07,
	0x4a, 0xa2, 0x59, 0xb0, 0x88, 0x56, 0x7e, 0x79, 0xf8, 0x09, 0xd4, 0xcd, 0xfa, 0xa8, 0xf7, 0x3b,
	0xd9, 0x5d, 0x53, 0x12, 0x3c, 0x29, 0xf5, 0x6a, 0x09, 0x85, 0x3e, 0x7a, 0xb2, 0x30, 0xd1, 0xa3,
	0x87, 0x5e, 0xe6, 0xe9, 0x54, 0x7b, 0x17, 0x3d, 0xba, 0xe8, 0xb7, 0xbd, 0xcb, 0xfc, 0x0c, 0x9e,
	0x55, 0x3c, 0x76, 0x86, 0xe9, 0x21, 0x73, 0xec, 0x71, 0x3a, 0xfd, 0x72, 0xff, 0x99, 0xd5, 0x8d,
	0x7c, 0x0d, 0x85, 0x2e, 0xd9, 0x64, 0x63, 0xfa, 0x30, 0xfc, 0xff, 0xe5, 0xe2, 0xd2, 0xb1, 0xc7,
	0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf0, 0x21, 0xac, 0xfb, 0x82, 0x01, 0x00, 0x00,
}
