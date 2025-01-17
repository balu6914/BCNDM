// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.21.12
// source: proto/dproxy.proto

package dproxy

import (
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SortOrder int32

const (
	SortOrder_ASC  SortOrder = 0
	SortOrder_DESC SortOrder = 1
)

// Enum value maps for SortOrder.
var (
	SortOrder_name = map[int32]string{
		0: "ASC",
		1: "DESC",
	}
	SortOrder_value = map[string]int32{
		"ASC":  0,
		"DESC": 1,
	}
)

func (x SortOrder) Enum() *SortOrder {
	p := new(SortOrder)
	*p = x
	return p
}

func (x SortOrder) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SortOrder) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_dproxy_proto_enumTypes[0].Descriptor()
}

func (SortOrder) Type() protoreflect.EnumType {
	return &file_proto_dproxy_proto_enumTypes[0]
}

func (x SortOrder) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SortOrder.Descriptor instead.
func (SortOrder) EnumDescriptor() ([]byte, []int) {
	return file_proto_dproxy_proto_rawDescGZIP(), []int{0}
}

type SortBy int32

const (
	SortBy_DATE SortBy = 0
)

// Enum value maps for SortBy.
var (
	SortBy_name = map[int32]string{
		0: "DATE",
	}
	SortBy_value = map[string]int32{
		"DATE": 0,
	}
)

func (x SortBy) Enum() *SortBy {
	p := new(SortBy)
	*p = x
	return p
}

func (x SortBy) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SortBy) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_dproxy_proto_enumTypes[1].Descriptor()
}

func (SortBy) Type() protoreflect.EnumType {
	return &file_proto_dproxy_proto_enumTypes[1]
}

func (x SortBy) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SortBy.Descriptor instead.
func (SortBy) EnumDescriptor() ([]byte, []int) {
	return file_proto_dproxy_proto_rawDescGZIP(), []int{1}
}

type AccessLog struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SubId string               `protobuf:"bytes,1,opt,name=subId,proto3" json:"subId,omitempty"`
	Time  *timestamp.Timestamp `protobuf:"bytes,2,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *AccessLog) Reset() {
	*x = AccessLog{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_dproxy_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessLog) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessLog) ProtoMessage() {}

func (x *AccessLog) ProtoReflect() protoreflect.Message {
	mi := &file_proto_dproxy_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessLog.ProtoReflect.Descriptor instead.
func (*AccessLog) Descriptor() ([]byte, []int) {
	return file_proto_dproxy_proto_rawDescGZIP(), []int{0}
}

func (x *AccessLog) GetSubId() string {
	if x != nil {
		return x.SubId
	}
	return ""
}

func (x *AccessLog) GetTime() *timestamp.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

type ListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Limit  uint32     `protobuf:"varint,1,opt,name=limit,proto3" json:"limit,omitempty"`
	Cursor *AccessLog `protobuf:"bytes,2,opt,name=cursor,proto3" json:"cursor,omitempty"`
	Sort   *Sort      `protobuf:"bytes,3,opt,name=sort,proto3" json:"sort,omitempty"`
}

func (x *ListRequest) Reset() {
	*x = ListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_dproxy_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRequest) ProtoMessage() {}

func (x *ListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_dproxy_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRequest.ProtoReflect.Descriptor instead.
func (*ListRequest) Descriptor() ([]byte, []int) {
	return file_proto_dproxy_proto_rawDescGZIP(), []int{1}
}

func (x *ListRequest) GetLimit() uint32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *ListRequest) GetCursor() *AccessLog {
	if x != nil {
		return x.Cursor
	}
	return nil
}

func (x *ListRequest) GetSort() *Sort {
	if x != nil {
		return x.Sort
	}
	return nil
}

type Sort struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Order SortOrder `protobuf:"varint,1,opt,name=order,proto3,enum=dproxy.SortOrder" json:"order,omitempty"`
	By    SortBy    `protobuf:"varint,2,opt,name=by,proto3,enum=dproxy.SortBy" json:"by,omitempty"`
}

func (x *Sort) Reset() {
	*x = Sort{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_dproxy_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sort) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sort) ProtoMessage() {}

func (x *Sort) ProtoReflect() protoreflect.Message {
	mi := &file_proto_dproxy_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sort.ProtoReflect.Descriptor instead.
func (*Sort) Descriptor() ([]byte, []int) {
	return file_proto_dproxy_proto_rawDescGZIP(), []int{2}
}

func (x *Sort) GetOrder() SortOrder {
	if x != nil {
		return x.Order
	}
	return SortOrder_ASC
}

func (x *Sort) GetBy() SortBy {
	if x != nil {
		return x.By
	}
	return SortBy_DATE
}

type ListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page []*AccessLog `protobuf:"bytes,1,rep,name=page,proto3" json:"page,omitempty"`
}

func (x *ListResponse) Reset() {
	*x = ListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_dproxy_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListResponse) ProtoMessage() {}

func (x *ListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_dproxy_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListResponse.ProtoReflect.Descriptor instead.
func (*ListResponse) Descriptor() ([]byte, []int) {
	return file_proto_dproxy_proto_rawDescGZIP(), []int{3}
}

func (x *ListResponse) GetPage() []*AccessLog {
	if x != nil {
		return x.Page
	}
	return nil
}

var File_proto_dproxy_proto protoreflect.FileDescriptor

var file_proto_dproxy_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x64, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x51, 0x0a,
	0x09, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4c, 0x6f, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x75,
	0x62, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x75, 0x62, 0x49, 0x64,
	0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65,
	0x22, 0x70, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05,
	0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x29, 0x0a, 0x06, 0x63, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x64, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x41,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x4c, 0x6f, 0x67, 0x52, 0x06, 0x63, 0x75, 0x72, 0x73, 0x6f, 0x72,
	0x12, 0x20, 0x0a, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c,
	0x2e, 0x64, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x53, 0x6f, 0x72, 0x74, 0x52, 0x04, 0x73, 0x6f,
	0x72, 0x74, 0x22, 0x4f, 0x0a, 0x04, 0x53, 0x6f, 0x72, 0x74, 0x12, 0x27, 0x0a, 0x05, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x64, 0x70, 0x72, 0x6f,
	0x78, 0x79, 0x2e, 0x53, 0x6f, 0x72, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x05, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x12, 0x1e, 0x0a, 0x02, 0x62, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x0e, 0x2e, 0x64, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x53, 0x6f, 0x72, 0x74, 0x42, 0x79, 0x52,
	0x02, 0x62, 0x79, 0x22, 0x35, 0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x11, 0x2e, 0x64, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x4c, 0x6f, 0x67, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x2a, 0x1e, 0x0a, 0x09, 0x53, 0x6f,
	0x72, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x07, 0x0a, 0x03, 0x41, 0x53, 0x43, 0x10, 0x00,
	0x12, 0x08, 0x0a, 0x04, 0x44, 0x45, 0x53, 0x43, 0x10, 0x01, 0x2a, 0x12, 0x0a, 0x06, 0x53, 0x6f,
	0x72, 0x74, 0x42, 0x79, 0x12, 0x08, 0x0a, 0x04, 0x44, 0x41, 0x54, 0x45, 0x10, 0x00, 0x32, 0x44,
	0x0a, 0x0d, 0x44, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x33, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x13, 0x2e, 0x64, 0x70, 0x72, 0x6f, 0x78, 0x79,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x64,
	0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x70, 0x61, 0x63, 0x65, 0x2f, 0x64, 0x61, 0x74, 0x61,
	0x70, 0x61, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x70, 0x72, 0x6f, 0x78,
	0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_dproxy_proto_rawDescOnce sync.Once
	file_proto_dproxy_proto_rawDescData = file_proto_dproxy_proto_rawDesc
)

func file_proto_dproxy_proto_rawDescGZIP() []byte {
	file_proto_dproxy_proto_rawDescOnce.Do(func() {
		file_proto_dproxy_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_dproxy_proto_rawDescData)
	})
	return file_proto_dproxy_proto_rawDescData
}

var file_proto_dproxy_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_proto_dproxy_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_dproxy_proto_goTypes = []interface{}{
	(SortOrder)(0),              // 0: dproxy.SortOrder
	(SortBy)(0),                 // 1: dproxy.SortBy
	(*AccessLog)(nil),           // 2: dproxy.AccessLog
	(*ListRequest)(nil),         // 3: dproxy.ListRequest
	(*Sort)(nil),                // 4: dproxy.Sort
	(*ListResponse)(nil),        // 5: dproxy.ListResponse
	(*timestamp.Timestamp)(nil), // 6: google.protobuf.Timestamp
}
var file_proto_dproxy_proto_depIdxs = []int32{
	6, // 0: dproxy.AccessLog.time:type_name -> google.protobuf.Timestamp
	2, // 1: dproxy.ListRequest.cursor:type_name -> dproxy.AccessLog
	4, // 2: dproxy.ListRequest.sort:type_name -> dproxy.Sort
	0, // 3: dproxy.Sort.order:type_name -> dproxy.SortOrder
	1, // 4: dproxy.Sort.by:type_name -> dproxy.SortBy
	2, // 5: dproxy.ListResponse.page:type_name -> dproxy.AccessLog
	3, // 6: dproxy.DproxyService.List:input_type -> dproxy.ListRequest
	5, // 7: dproxy.DproxyService.List:output_type -> dproxy.ListResponse
	7, // [7:8] is the sub-list for method output_type
	6, // [6:7] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_proto_dproxy_proto_init() }
func file_proto_dproxy_proto_init() {
	if File_proto_dproxy_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_dproxy_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessLog); i {
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
		file_proto_dproxy_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRequest); i {
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
		file_proto_dproxy_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Sort); i {
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
		file_proto_dproxy_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListResponse); i {
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
			RawDescriptor: file_proto_dproxy_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_dproxy_proto_goTypes,
		DependencyIndexes: file_proto_dproxy_proto_depIdxs,
		EnumInfos:         file_proto_dproxy_proto_enumTypes,
		MessageInfos:      file_proto_dproxy_proto_msgTypes,
	}.Build()
	File_proto_dproxy_proto = out.File
	file_proto_dproxy_proto_rawDesc = nil
	file_proto_dproxy_proto_goTypes = nil
	file_proto_dproxy_proto_depIdxs = nil
}
