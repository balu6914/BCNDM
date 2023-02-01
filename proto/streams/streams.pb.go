// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.21.9
// source: proto/streams.proto

package streams

import (
	common "github.com/datapace/datapace/proto/common"
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

type Stream struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name       string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Owner      string `protobuf:"bytes,3,opt,name=owner,proto3" json:"owner,omitempty"`
	Url        string `protobuf:"bytes,4,opt,name=url,proto3" json:"url,omitempty"`
	Price      uint64 `protobuf:"varint,5,opt,name=price,proto3" json:"price,omitempty"`
	External   bool   `protobuf:"varint,6,opt,name=external,proto3" json:"external,omitempty"`
	Project    string `protobuf:"bytes,7,opt,name=project,proto3" json:"project,omitempty"`
	Dataset    string `protobuf:"bytes,8,opt,name=dataset,proto3" json:"dataset,omitempty"`
	Table      string `protobuf:"bytes,9,opt,name=table,proto3" json:"table,omitempty"`
	Fields     string `protobuf:"bytes,10,opt,name=fields,proto3" json:"fields,omitempty"`
	Visibility string `protobuf:"bytes,11,opt,name=visibility,proto3" json:"visibility,omitempty"`
	AccessType string `protobuf:"bytes,12,opt,name=accessType,proto3" json:"accessType,omitempty"`
	MaxCalls   uint64 `protobuf:"varint,13,opt,name=maxCalls,proto3" json:"maxCalls,omitempty"`
	MaxUnit    string `protobuf:"bytes,14,opt,name=maxUnit,proto3" json:"maxUnit,omitempty"`
	EndDate    string `protobuf:"bytes,15,opt,name=endDate,proto3" json:"endDate,omitempty"`
}

func (x *Stream) Reset() {
	*x = Stream{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_streams_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Stream) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stream) ProtoMessage() {}

func (x *Stream) ProtoReflect() protoreflect.Message {
	mi := &file_proto_streams_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stream.ProtoReflect.Descriptor instead.
func (*Stream) Descriptor() ([]byte, []int) {
	return file_proto_streams_proto_rawDescGZIP(), []int{0}
}

func (x *Stream) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Stream) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Stream) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *Stream) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Stream) GetPrice() uint64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *Stream) GetExternal() bool {
	if x != nil {
		return x.External
	}
	return false
}

func (x *Stream) GetProject() string {
	if x != nil {
		return x.Project
	}
	return ""
}

func (x *Stream) GetDataset() string {
	if x != nil {
		return x.Dataset
	}
	return ""
}

func (x *Stream) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *Stream) GetFields() string {
	if x != nil {
		return x.Fields
	}
	return ""
}

func (x *Stream) GetVisibility() string {
	if x != nil {
		return x.Visibility
	}
	return ""
}

func (x *Stream) GetAccessType() string {
	if x != nil {
		return x.AccessType
	}
	return ""
}

func (x *Stream) GetMaxCalls() uint64 {
	if x != nil {
		return x.MaxCalls
	}
	return 0
}

func (x *Stream) GetMaxUnit() string {
	if x != nil {
		return x.MaxUnit
	}
	return ""
}

func (x *Stream) GetEndDate() string {
	if x != nil {
		return x.EndDate
	}
	return ""
}

var File_proto_streams_proto protoreflect.FileDescriptor

var file_proto_streams_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x64, 0x61, 0x74, 0x61, 0x70, 0x61, 0x63, 0x65, 0x1a,
	0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xf8, 0x02, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72,
	0x69, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x08, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x12, 0x18, 0x0a, 0x07,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65,
	0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x1e,
	0x0a, 0x0a, 0x76, 0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x76, 0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x12, 0x1e,
	0x0a, 0x0a, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x18, 0x0c, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x6d, 0x61, 0x78, 0x43, 0x61, 0x6c, 0x6c, 0x73, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x08, 0x6d, 0x61, 0x78, 0x43, 0x61, 0x6c, 0x6c, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x61,
	0x78, 0x55, 0x6e, 0x69, 0x74, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x61, 0x78,
	0x55, 0x6e, 0x69, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x6e, 0x64, 0x44, 0x61, 0x74, 0x65, 0x18,
	0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x44, 0x61, 0x74, 0x65, 0x32, 0x39,
	0x0a, 0x0e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x27, 0x0a, 0x03, 0x4f, 0x6e, 0x65, 0x12, 0x0c, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x70, 0x61,
	0x63, 0x65, 0x2e, 0x49, 0x44, 0x1a, 0x10, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x70, 0x61, 0x63, 0x65,
	0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x22, 0x00, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x70, 0x61, 0x63, 0x65,
	0x2f, 0x64, 0x61, 0x74, 0x61, 0x70, 0x61, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_streams_proto_rawDescOnce sync.Once
	file_proto_streams_proto_rawDescData = file_proto_streams_proto_rawDesc
)

func file_proto_streams_proto_rawDescGZIP() []byte {
	file_proto_streams_proto_rawDescOnce.Do(func() {
		file_proto_streams_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_streams_proto_rawDescData)
	})
	return file_proto_streams_proto_rawDescData
}

var file_proto_streams_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_proto_streams_proto_goTypes = []interface{}{
	(*Stream)(nil),    // 0: datapace.Stream
	(*common.ID)(nil), // 1: datapace.ID
}
var file_proto_streams_proto_depIdxs = []int32{
	1, // 0: datapace.StreamsService.One:input_type -> datapace.ID
	0, // 1: datapace.StreamsService.One:output_type -> datapace.Stream
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_streams_proto_init() }
func file_proto_streams_proto_init() {
	if File_proto_streams_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_streams_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Stream); i {
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
			RawDescriptor: file_proto_streams_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_streams_proto_goTypes,
		DependencyIndexes: file_proto_streams_proto_depIdxs,
		MessageInfos:      file_proto_streams_proto_msgTypes,
	}.Build()
	File_proto_streams_proto = out.File
	file_proto_streams_proto_rawDesc = nil
	file_proto_streams_proto_goTypes = nil
	file_proto_streams_proto_depIdxs = nil
}
