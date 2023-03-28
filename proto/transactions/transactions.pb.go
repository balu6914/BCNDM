// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.21.9
// source: proto/transactions.proto

package transactions

import (
	common "github.com/datapace/datapace/proto/common"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TransferData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StreamID string `protobuf:"bytes,1,opt,name=streamID,proto3" json:"streamID,omitempty"`
	From     string `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	To       string `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`
	Value    uint64 `protobuf:"varint,4,opt,name=value,proto3" json:"value,omitempty"`
	DateTime string `protobuf:"bytes,5,opt,name=dateTime,proto3" json:"dateTime,omitempty"`
}

func (x *TransferData) Reset() {
	*x = TransferData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_transactions_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransferData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransferData) ProtoMessage() {}

func (x *TransferData) ProtoReflect() protoreflect.Message {
	mi := &file_proto_transactions_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransferData.ProtoReflect.Descriptor instead.
func (*TransferData) Descriptor() ([]byte, []int) {
	return file_proto_transactions_proto_rawDescGZIP(), []int{0}
}

func (x *TransferData) GetStreamID() string {
	if x != nil {
		return x.StreamID
	}
	return ""
}

func (x *TransferData) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *TransferData) GetTo() string {
	if x != nil {
		return x.To
	}
	return ""
}

func (x *TransferData) GetValue() uint64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *TransferData) GetDateTime() string {
	if x != nil {
		return x.DateTime
	}
	return ""
}

var File_proto_transactions_proto protoreflect.FileDescriptor

var file_proto_transactions_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x64, 0x61, 0x74, 0x61,
	0x70, 0x61, 0x63, 0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x80, 0x01, 0x0a, 0x0c, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66,
	0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x6f, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x74, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x32, 0x89, 0x01, 0x0a, 0x13, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x34, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0c,
	0x2e, 0x64, 0x61, 0x74, 0x61, 0x70, 0x61, 0x63, 0x65, 0x2e, 0x49, 0x44, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x3c, 0x0a, 0x08, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66,
	0x65, 0x72, 0x12, 0x16, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x70, 0x61, 0x63, 0x65, 0x2e, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x22, 0x00, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x70, 0x61, 0x63, 0x65, 0x2f, 0x64, 0x61, 0x74, 0x61,
	0x70, 0x61, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_transactions_proto_rawDescOnce sync.Once
	file_proto_transactions_proto_rawDescData = file_proto_transactions_proto_rawDesc
)

func file_proto_transactions_proto_rawDescGZIP() []byte {
	file_proto_transactions_proto_rawDescOnce.Do(func() {
		file_proto_transactions_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_transactions_proto_rawDescData)
	})
	return file_proto_transactions_proto_rawDescData
}

var file_proto_transactions_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_proto_transactions_proto_goTypes = []interface{}{
	(*TransferData)(nil),  // 0: datapace.TransferData
	(*common.ID)(nil),     // 1: datapace.ID
	(*emptypb.Empty)(nil), // 2: google.protobuf.Empty
}
var file_proto_transactions_proto_depIdxs = []int32{
	1, // 0: datapace.TransactionsService.CreateUser:input_type -> datapace.ID
	0, // 1: datapace.TransactionsService.Transfer:input_type -> datapace.TransferData
	2, // 2: datapace.TransactionsService.CreateUser:output_type -> google.protobuf.Empty
	2, // 3: datapace.TransactionsService.Transfer:output_type -> google.protobuf.Empty
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_transactions_proto_init() }
func file_proto_transactions_proto_init() {
	if File_proto_transactions_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_transactions_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TransferData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_transactions_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_transactions_proto_goTypes,
		DependencyIndexes: file_proto_transactions_proto_depIdxs,
		MessageInfos:      file_proto_transactions_proto_msgTypes,
	}.Build()
	File_proto_transactions_proto = out.File
	file_proto_transactions_proto_rawDesc = nil
	file_proto_transactions_proto_goTypes = nil
	file_proto_transactions_proto_depIdxs = nil
}
