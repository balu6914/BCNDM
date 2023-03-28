// Code generated by protoc-gen-go. DO NOT EDIT.
// source: common/configtx.proto

package common

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

// ConfigEnvelope is designed to contain _all_ configuration for a chain with no dependency
// on previous configuration transactions.
//
// It is generated with the following scheme:
//   1. Retrieve the existing configuration
//   2. Note the config properties (ConfigValue, ConfigPolicy, ConfigGroup) to be modified
//   3. Add any intermediate ConfigGroups to the ConfigUpdate.read_set (sparsely)
//   4. Add any additional desired dependencies to ConfigUpdate.read_set (sparsely)
//   5. Modify the config properties, incrementing each version by 1, set them in the ConfigUpdate.write_set
//      Note: any element not modified but specified should already be in the read_set, so may be specified sparsely
//   6. Create ConfigUpdate message and marshal it into ConfigUpdateEnvelope.update and encode the required signatures
//     a) Each signature is of type ConfigSignature
//     b) The ConfigSignature signature is over the concatenation of signature_header and the ConfigUpdate bytes (which includes a ChainHeader)
//   5. Submit new Config for ordering in Envelope signed by submitter
//     a) The Envelope Payload has data set to the marshaled ConfigEnvelope
//     b) The Envelope Payload has a header of type Header.Type.CONFIG_UPDATE
//
// The configuration manager will verify:
//   1. All items in the read_set exist at the read versions
//   2. All items in the write_set at a different version than, or not in, the read_set have been appropriately signed according to their mod_policy
//   3. The new configuration satisfies the ConfigSchema
type ConfigEnvelope struct {
	Config               *Config   `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
	LastUpdate           *Envelope `protobuf:"bytes,2,opt,name=last_update,json=lastUpdate,proto3" json:"last_update,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ConfigEnvelope) Reset()         { *m = ConfigEnvelope{} }
func (m *ConfigEnvelope) String() string { return proto.CompactTextString(m) }
func (*ConfigEnvelope) ProtoMessage()    {}
func (*ConfigEnvelope) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{0}
}

func (m *ConfigEnvelope) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigEnvelope.Unmarshal(m, b)
}
func (m *ConfigEnvelope) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigEnvelope.Marshal(b, m, deterministic)
}
func (m *ConfigEnvelope) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigEnvelope.Merge(m, src)
}
func (m *ConfigEnvelope) XXX_Size() int {
	return xxx_messageInfo_ConfigEnvelope.Size(m)
}
func (m *ConfigEnvelope) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigEnvelope.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigEnvelope proto.InternalMessageInfo

func (m *ConfigEnvelope) GetConfig() *Config {
	if m != nil {
		return m.Config
	}
	return nil
}

func (m *ConfigEnvelope) GetLastUpdate() *Envelope {
	if m != nil {
		return m.LastUpdate
	}
	return nil
}

type ConfigGroupSchema struct {
	Groups               map[string]*ConfigGroupSchema  `protobuf:"bytes,1,rep,name=groups,proto3" json:"groups,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Values               map[string]*ConfigValueSchema  `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Policies             map[string]*ConfigPolicySchema `protobuf:"bytes,3,rep,name=policies,proto3" json:"policies,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                       `json:"-"`
	XXX_unrecognized     []byte                         `json:"-"`
	XXX_sizecache        int32                          `json:"-"`
}

func (m *ConfigGroupSchema) Reset()         { *m = ConfigGroupSchema{} }
func (m *ConfigGroupSchema) String() string { return proto.CompactTextString(m) }
func (*ConfigGroupSchema) ProtoMessage()    {}
func (*ConfigGroupSchema) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{1}
}

func (m *ConfigGroupSchema) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigGroupSchema.Unmarshal(m, b)
}
func (m *ConfigGroupSchema) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigGroupSchema.Marshal(b, m, deterministic)
}
func (m *ConfigGroupSchema) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigGroupSchema.Merge(m, src)
}
func (m *ConfigGroupSchema) XXX_Size() int {
	return xxx_messageInfo_ConfigGroupSchema.Size(m)
}
func (m *ConfigGroupSchema) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigGroupSchema.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigGroupSchema proto.InternalMessageInfo

func (m *ConfigGroupSchema) GetGroups() map[string]*ConfigGroupSchema {
	if m != nil {
		return m.Groups
	}
	return nil
}

func (m *ConfigGroupSchema) GetValues() map[string]*ConfigValueSchema {
	if m != nil {
		return m.Values
	}
	return nil
}

func (m *ConfigGroupSchema) GetPolicies() map[string]*ConfigPolicySchema {
	if m != nil {
		return m.Policies
	}
	return nil
}

type ConfigValueSchema struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfigValueSchema) Reset()         { *m = ConfigValueSchema{} }
func (m *ConfigValueSchema) String() string { return proto.CompactTextString(m) }
func (*ConfigValueSchema) ProtoMessage()    {}
func (*ConfigValueSchema) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{2}
}

func (m *ConfigValueSchema) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigValueSchema.Unmarshal(m, b)
}
func (m *ConfigValueSchema) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigValueSchema.Marshal(b, m, deterministic)
}
func (m *ConfigValueSchema) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigValueSchema.Merge(m, src)
}
func (m *ConfigValueSchema) XXX_Size() int {
	return xxx_messageInfo_ConfigValueSchema.Size(m)
}
func (m *ConfigValueSchema) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigValueSchema.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigValueSchema proto.InternalMessageInfo

type ConfigPolicySchema struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfigPolicySchema) Reset()         { *m = ConfigPolicySchema{} }
func (m *ConfigPolicySchema) String() string { return proto.CompactTextString(m) }
func (*ConfigPolicySchema) ProtoMessage()    {}
func (*ConfigPolicySchema) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{3}
}

func (m *ConfigPolicySchema) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigPolicySchema.Unmarshal(m, b)
}
func (m *ConfigPolicySchema) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigPolicySchema.Marshal(b, m, deterministic)
}
func (m *ConfigPolicySchema) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigPolicySchema.Merge(m, src)
}
func (m *ConfigPolicySchema) XXX_Size() int {
	return xxx_messageInfo_ConfigPolicySchema.Size(m)
}
func (m *ConfigPolicySchema) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigPolicySchema.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigPolicySchema proto.InternalMessageInfo

// Config represents the config for a particular channel
type Config struct {
	Sequence             uint64       `protobuf:"varint,1,opt,name=sequence,proto3" json:"sequence,omitempty"`
	ChannelGroup         *ConfigGroup `protobuf:"bytes,2,opt,name=channel_group,json=channelGroup,proto3" json:"channel_group,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{4}
}

func (m *Config) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Config.Unmarshal(m, b)
}
func (m *Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Config.Marshal(b, m, deterministic)
}
func (m *Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config.Merge(m, src)
}
func (m *Config) XXX_Size() int {
	return xxx_messageInfo_Config.Size(m)
}
func (m *Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Config proto.InternalMessageInfo

func (m *Config) GetSequence() uint64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *Config) GetChannelGroup() *ConfigGroup {
	if m != nil {
		return m.ChannelGroup
	}
	return nil
}

type ConfigUpdateEnvelope struct {
	ConfigUpdate         []byte             `protobuf:"bytes,1,opt,name=config_update,json=configUpdate,proto3" json:"config_update,omitempty"`
	Signatures           []*ConfigSignature `protobuf:"bytes,2,rep,name=signatures,proto3" json:"signatures,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *ConfigUpdateEnvelope) Reset()         { *m = ConfigUpdateEnvelope{} }
func (m *ConfigUpdateEnvelope) String() string { return proto.CompactTextString(m) }
func (*ConfigUpdateEnvelope) ProtoMessage()    {}
func (*ConfigUpdateEnvelope) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{5}
}

func (m *ConfigUpdateEnvelope) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigUpdateEnvelope.Unmarshal(m, b)
}
func (m *ConfigUpdateEnvelope) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigUpdateEnvelope.Marshal(b, m, deterministic)
}
func (m *ConfigUpdateEnvelope) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigUpdateEnvelope.Merge(m, src)
}
func (m *ConfigUpdateEnvelope) XXX_Size() int {
	return xxx_messageInfo_ConfigUpdateEnvelope.Size(m)
}
func (m *ConfigUpdateEnvelope) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigUpdateEnvelope.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigUpdateEnvelope proto.InternalMessageInfo

func (m *ConfigUpdateEnvelope) GetConfigUpdate() []byte {
	if m != nil {
		return m.ConfigUpdate
	}
	return nil
}

func (m *ConfigUpdateEnvelope) GetSignatures() []*ConfigSignature {
	if m != nil {
		return m.Signatures
	}
	return nil
}

// ConfigUpdate is used to submit a subset of config and to have the orderer apply to Config
// it is always submitted inside a ConfigUpdateEnvelope which allows the addition of signatures
// resulting in a new total configuration.  The update is applied as follows:
// 1. The versions from all of the elements in the read_set is verified against the versions in the existing config.
//    If there is a mismatch in the read versions, then the config update fails and is rejected.
// 2. Any elements in the write_set with the same version as the read_set are ignored.
// 3. The corresponding mod_policy for every remaining element in the write_set is collected.
// 4. Each policy is checked against the signatures from the ConfigUpdateEnvelope, any failing to verify are rejected
// 5. The write_set is applied to the Config and the ConfigGroupSchema verifies that the updates were legal
type ConfigUpdate struct {
	ChannelId            string            `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	ReadSet              *ConfigGroup      `protobuf:"bytes,2,opt,name=read_set,json=readSet,proto3" json:"read_set,omitempty"`
	WriteSet             *ConfigGroup      `protobuf:"bytes,3,opt,name=write_set,json=writeSet,proto3" json:"write_set,omitempty"`
	IsolatedData         map[string][]byte `protobuf:"bytes,5,rep,name=isolated_data,json=isolatedData,proto3" json:"isolated_data,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ConfigUpdate) Reset()         { *m = ConfigUpdate{} }
func (m *ConfigUpdate) String() string { return proto.CompactTextString(m) }
func (*ConfigUpdate) ProtoMessage()    {}
func (*ConfigUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{6}
}

func (m *ConfigUpdate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigUpdate.Unmarshal(m, b)
}
func (m *ConfigUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigUpdate.Marshal(b, m, deterministic)
}
func (m *ConfigUpdate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigUpdate.Merge(m, src)
}
func (m *ConfigUpdate) XXX_Size() int {
	return xxx_messageInfo_ConfigUpdate.Size(m)
}
func (m *ConfigUpdate) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigUpdate.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigUpdate proto.InternalMessageInfo

func (m *ConfigUpdate) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func (m *ConfigUpdate) GetReadSet() *ConfigGroup {
	if m != nil {
		return m.ReadSet
	}
	return nil
}

func (m *ConfigUpdate) GetWriteSet() *ConfigGroup {
	if m != nil {
		return m.WriteSet
	}
	return nil
}

func (m *ConfigUpdate) GetIsolatedData() map[string][]byte {
	if m != nil {
		return m.IsolatedData
	}
	return nil
}

// ConfigGroup is the hierarchical data structure for holding config
type ConfigGroup struct {
	Version              uint64                   `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	Groups               map[string]*ConfigGroup  `protobuf:"bytes,2,rep,name=groups,proto3" json:"groups,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Values               map[string]*ConfigValue  `protobuf:"bytes,3,rep,name=values,proto3" json:"values,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Policies             map[string]*ConfigPolicy `protobuf:"bytes,4,rep,name=policies,proto3" json:"policies,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ModPolicy            string                   `protobuf:"bytes,5,opt,name=mod_policy,json=modPolicy,proto3" json:"mod_policy,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *ConfigGroup) Reset()         { *m = ConfigGroup{} }
func (m *ConfigGroup) String() string { return proto.CompactTextString(m) }
func (*ConfigGroup) ProtoMessage()    {}
func (*ConfigGroup) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{7}
}

func (m *ConfigGroup) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigGroup.Unmarshal(m, b)
}
func (m *ConfigGroup) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigGroup.Marshal(b, m, deterministic)
}
func (m *ConfigGroup) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigGroup.Merge(m, src)
}
func (m *ConfigGroup) XXX_Size() int {
	return xxx_messageInfo_ConfigGroup.Size(m)
}
func (m *ConfigGroup) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigGroup.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigGroup proto.InternalMessageInfo

func (m *ConfigGroup) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *ConfigGroup) GetGroups() map[string]*ConfigGroup {
	if m != nil {
		return m.Groups
	}
	return nil
}

func (m *ConfigGroup) GetValues() map[string]*ConfigValue {
	if m != nil {
		return m.Values
	}
	return nil
}

func (m *ConfigGroup) GetPolicies() map[string]*ConfigPolicy {
	if m != nil {
		return m.Policies
	}
	return nil
}

func (m *ConfigGroup) GetModPolicy() string {
	if m != nil {
		return m.ModPolicy
	}
	return ""
}

// ConfigValue represents an individual piece of config data
type ConfigValue struct {
	Version              uint64   `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	Value                []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	ModPolicy            string   `protobuf:"bytes,3,opt,name=mod_policy,json=modPolicy,proto3" json:"mod_policy,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfigValue) Reset()         { *m = ConfigValue{} }
func (m *ConfigValue) String() string { return proto.CompactTextString(m) }
func (*ConfigValue) ProtoMessage()    {}
func (*ConfigValue) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{8}
}

func (m *ConfigValue) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigValue.Unmarshal(m, b)
}
func (m *ConfigValue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigValue.Marshal(b, m, deterministic)
}
func (m *ConfigValue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigValue.Merge(m, src)
}
func (m *ConfigValue) XXX_Size() int {
	return xxx_messageInfo_ConfigValue.Size(m)
}
func (m *ConfigValue) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigValue.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigValue proto.InternalMessageInfo

func (m *ConfigValue) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *ConfigValue) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *ConfigValue) GetModPolicy() string {
	if m != nil {
		return m.ModPolicy
	}
	return ""
}

type ConfigPolicy struct {
	Version              uint64   `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	Policy               *Policy  `protobuf:"bytes,2,opt,name=policy,proto3" json:"policy,omitempty"`
	ModPolicy            string   `protobuf:"bytes,3,opt,name=mod_policy,json=modPolicy,proto3" json:"mod_policy,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfigPolicy) Reset()         { *m = ConfigPolicy{} }
func (m *ConfigPolicy) String() string { return proto.CompactTextString(m) }
func (*ConfigPolicy) ProtoMessage()    {}
func (*ConfigPolicy) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{9}
}

func (m *ConfigPolicy) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigPolicy.Unmarshal(m, b)
}
func (m *ConfigPolicy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigPolicy.Marshal(b, m, deterministic)
}
func (m *ConfigPolicy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigPolicy.Merge(m, src)
}
func (m *ConfigPolicy) XXX_Size() int {
	return xxx_messageInfo_ConfigPolicy.Size(m)
}
func (m *ConfigPolicy) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigPolicy.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigPolicy proto.InternalMessageInfo

func (m *ConfigPolicy) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *ConfigPolicy) GetPolicy() *Policy {
	if m != nil {
		return m.Policy
	}
	return nil
}

func (m *ConfigPolicy) GetModPolicy() string {
	if m != nil {
		return m.ModPolicy
	}
	return ""
}

type ConfigSignature struct {
	SignatureHeader      []byte   `protobuf:"bytes,1,opt,name=signature_header,json=signatureHeader,proto3" json:"signature_header,omitempty"`
	Signature            []byte   `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfigSignature) Reset()         { *m = ConfigSignature{} }
func (m *ConfigSignature) String() string { return proto.CompactTextString(m) }
func (*ConfigSignature) ProtoMessage()    {}
func (*ConfigSignature) Descriptor() ([]byte, []int) {
	return fileDescriptor_5190bbf196fa7499, []int{10}
}

func (m *ConfigSignature) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigSignature.Unmarshal(m, b)
}
func (m *ConfigSignature) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigSignature.Marshal(b, m, deterministic)
}
func (m *ConfigSignature) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigSignature.Merge(m, src)
}
func (m *ConfigSignature) XXX_Size() int {
	return xxx_messageInfo_ConfigSignature.Size(m)
}
func (m *ConfigSignature) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigSignature.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigSignature proto.InternalMessageInfo

func (m *ConfigSignature) GetSignatureHeader() []byte {
	if m != nil {
		return m.SignatureHeader
	}
	return nil
}

func (m *ConfigSignature) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func init() {
	proto.RegisterType((*ConfigEnvelope)(nil), "common.ConfigEnvelope")
	proto.RegisterType((*ConfigGroupSchema)(nil), "common.ConfigGroupSchema")
	proto.RegisterMapType((map[string]*ConfigGroupSchema)(nil), "common.ConfigGroupSchema.GroupsEntry")
	proto.RegisterMapType((map[string]*ConfigPolicySchema)(nil), "common.ConfigGroupSchema.PoliciesEntry")
	proto.RegisterMapType((map[string]*ConfigValueSchema)(nil), "common.ConfigGroupSchema.ValuesEntry")
	proto.RegisterType((*ConfigValueSchema)(nil), "common.ConfigValueSchema")
	proto.RegisterType((*ConfigPolicySchema)(nil), "common.ConfigPolicySchema")
	proto.RegisterType((*Config)(nil), "common.Config")
	proto.RegisterType((*ConfigUpdateEnvelope)(nil), "common.ConfigUpdateEnvelope")
	proto.RegisterType((*ConfigUpdate)(nil), "common.ConfigUpdate")
	proto.RegisterMapType((map[string][]byte)(nil), "common.ConfigUpdate.IsolatedDataEntry")
	proto.RegisterType((*ConfigGroup)(nil), "common.ConfigGroup")
	proto.RegisterMapType((map[string]*ConfigGroup)(nil), "common.ConfigGroup.GroupsEntry")
	proto.RegisterMapType((map[string]*ConfigPolicy)(nil), "common.ConfigGroup.PoliciesEntry")
	proto.RegisterMapType((map[string]*ConfigValue)(nil), "common.ConfigGroup.ValuesEntry")
	proto.RegisterType((*ConfigValue)(nil), "common.ConfigValue")
	proto.RegisterType((*ConfigPolicy)(nil), "common.ConfigPolicy")
	proto.RegisterType((*ConfigSignature)(nil), "common.ConfigSignature")
}

func init() { proto.RegisterFile("common/configtx.proto", fileDescriptor_5190bbf196fa7499) }

var fileDescriptor_5190bbf196fa7499 = []byte{
	// 744 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x56, 0xe1, 0x6e, 0xd3, 0x3c,
	0x14, 0x55, 0x9b, 0xb6, 0x6b, 0x6f, 0xdb, 0xad, 0xf3, 0xfa, 0xe9, 0x0b, 0x15, 0x88, 0x11, 0x60,
	0x6c, 0x48, 0x4b, 0xc7, 0xf8, 0xb1, 0x09, 0x69, 0x42, 0x62, 0x4c, 0xb0, 0x21, 0x4d, 0x90, 0xc1,
	0x90, 0x26, 0xa4, 0xca, 0x4b, 0xbc, 0x34, 0x2c, 0x8d, 0x43, 0xe2, 0x0e, 0xfa, 0x48, 0x3c, 0x13,
	0x6f, 0xc0, 0x53, 0xa0, 0xd8, 0x4e, 0x70, 0xd6, 0xb4, 0x15, 0xbf, 0x56, 0x5f, 0x9f, 0x73, 0xee,
	0xb5, 0xef, 0xf5, 0xc9, 0xe0, 0x3f, 0x9b, 0x8e, 0x46, 0x34, 0xe8, 0xdb, 0x34, 0xb8, 0xf2, 0x5c,
	0xf6, 0xc3, 0x0c, 0x23, 0xca, 0x28, 0xaa, 0x89, 0x70, 0x6f, 0x2d, 0xdb, 0x4e, 0xfe, 0x88, 0xcd,
	0x5e, 0xca, 0x09, 0xa9, 0xef, 0xd9, 0x1e, 0x89, 0x45, 0xd8, 0xb8, 0x86, 0xe5, 0x43, 0xae, 0x72,
	0x14, 0xdc, 0x10, 0x9f, 0x86, 0x04, 0x6d, 0x40, 0x4d, 0xe8, 0xea, 0xa5, 0xf5, 0xd2, 0x66, 0x73,
	0x77, 0xd9, 0x94, 0x3a, 0x02, 0x67, 0xc9, 0x5d, 0xf4, 0x0c, 0x9a, 0x3e, 0x8e, 0xd9, 0x60, 0x1c,
	0x3a, 0x98, 0x11, 0xbd, 0xcc, 0xc1, 0x9d, 0x14, 0x9c, 0xca, 0x59, 0x90, 0x80, 0x3e, 0x71, 0x8c,
	0xf1, 0x4b, 0x83, 0x55, 0xa1, 0xf2, 0x26, 0xa2, 0xe3, 0xf0, 0xcc, 0x1e, 0x92, 0x11, 0x46, 0x07,
	0x50, 0x73, 0x93, 0x65, 0xac, 0x97, 0xd6, 0xb5, 0xcd, 0xe6, 0xee, 0xe3, 0x7c, 0x42, 0x05, 0x6a,
	0xf2, 0xdf, 0xf1, 0x51, 0xc0, 0xa2, 0x89, 0x25, 0x49, 0x09, 0xfd, 0x06, 0xfb, 0x63, 0x12, 0xeb,
	0xe5, 0x45, 0xf4, 0x73, 0x8e, 0x93, 0x74, 0x41, 0x42, 0x87, 0x50, 0x4f, 0xaf, 0x44, 0xd7, 0xb8,
	0xc0, 0x93, 0xd9, 0x02, 0xef, 0x25, 0x52, 0x48, 0x64, 0xc4, 0xde, 0x47, 0x68, 0x2a, 0xa5, 0xa1,
	0x0e, 0x68, 0xd7, 0x64, 0xc2, 0xef, 0xaf, 0x61, 0x25, 0x3f, 0x51, 0x1f, 0xaa, 0x3c, 0x9f, 0xbc,
	0xa6, 0x3b, 0x33, 0x53, 0x58, 0x02, 0xf7, 0xa2, 0xbc, 0x5f, 0x4a, 0x54, 0x95, 0x8a, 0xff, 0x59,
	0x95, 0x73, 0xa7, 0x55, 0x3f, 0x43, 0x3b, 0x77, 0x8c, 0x02, 0xdd, 0x9d, 0xbc, 0x6e, 0x2f, 0xaf,
	0xcb, 0xd9, 0x93, 0x29, 0x61, 0x63, 0x2d, 0x6d, 0xae, 0x92, 0xd8, 0xe8, 0x02, 0x9a, 0x66, 0x19,
	0x5f, 0xa1, 0x26, 0xa2, 0xa8, 0x07, 0xf5, 0x98, 0x7c, 0x1b, 0x93, 0xc0, 0x26, 0xbc, 0x82, 0x8a,
	0x95, 0xad, 0xd1, 0x3e, 0xb4, 0xed, 0x21, 0x0e, 0x02, 0xe2, 0x0f, 0x78, 0xaf, 0x65, 0x39, 0x6b,
	0x05, 0x97, 0x67, 0xb5, 0x24, 0x92, 0xaf, 0x4e, 0x2a, 0x75, 0xad, 0x53, 0xb1, 0x2a, 0x6c, 0x12,
	0x12, 0x83, 0x41, 0x57, 0x00, 0xc5, 0x10, 0x66, 0x73, 0xfe, 0x10, 0xda, 0x62, 0x92, 0xd3, 0x09,
	0x4e, 0xd2, 0xb7, 0xac, 0x96, 0xad, 0x80, 0xd1, 0x1e, 0x40, 0xec, 0xb9, 0x01, 0x66, 0xe3, 0x28,
	0x1b, 0xb0, 0xff, 0xf3, 0xf9, 0xcf, 0xd2, 0x7d, 0x4b, 0x81, 0x1a, 0x3f, 0xcb, 0xd0, 0x52, 0xd3,
	0xa2, 0x7b, 0x00, 0xe9, 0x61, 0x3c, 0x47, 0x5e, 0x76, 0x43, 0x46, 0x8e, 0x1d, 0x64, 0x42, 0x3d,
	0x22, 0xd8, 0x19, 0xc4, 0x84, 0xcd, 0x3b, 0xe6, 0x52, 0x02, 0x3a, 0x23, 0x0c, 0xed, 0x40, 0xe3,
	0x7b, 0xe4, 0x31, 0xc2, 0x09, 0xda, 0x6c, 0x42, 0x9d, 0xa3, 0x12, 0xc6, 0x3b, 0x68, 0x7b, 0x31,
	0xf5, 0x31, 0x23, 0xce, 0xc0, 0xc1, 0x0c, 0xeb, 0x55, 0x7e, 0x9a, 0x8d, 0x3c, 0x4b, 0x54, 0x6b,
	0x1e, 0x4b, 0xe4, 0x6b, 0xcc, 0xb0, 0x18, 0xf6, 0x96, 0xa7, 0x84, 0x7a, 0x2f, 0x61, 0x75, 0x0a,
	0x52, 0x30, 0x48, 0x5d, 0x75, 0x90, 0x5a, 0xca, 0xb0, 0x9c, 0x54, 0xea, 0x95, 0x4e, 0x55, 0x76,
	0xe8, 0xb7, 0x06, 0x4d, 0xa5, 0x66, 0xa4, 0xc3, 0xd2, 0x0d, 0x89, 0x62, 0x8f, 0x06, 0x72, 0x24,
	0xd2, 0x25, 0xda, 0xcb, 0xac, 0x42, 0xb4, 0xe2, 0x7e, 0xc1, 0x91, 0x0b, 0x4d, 0x62, 0x2f, 0x33,
	0x09, 0x6d, 0x36, 0xb1, 0xc8, 0x1e, 0x0e, 0x14, 0x7b, 0xa8, 0x70, 0xea, 0x83, 0x22, 0xea, 0x0c,
	0x63, 0x48, 0xba, 0x3e, 0xa2, 0xce, 0x80, 0xaf, 0x27, 0x7a, 0x55, 0x74, 0x7d, 0x44, 0x1d, 0xf1,
	0x1a, 0x7a, 0xa7, 0x8b, 0x7c, 0x63, 0x2b, 0xff, 0x12, 0x0b, 0x5b, 0xac, 0xbc, 0xed, 0xd3, 0x45,
	0x8e, 0x31, 0x5f, 0x8f, 0x73, 0x55, 0xbd, 0x0f, 0x8b, 0xbd, 0xe2, 0x69, 0x5e, 0xb1, 0x5b, 0xe4,
	0x15, 0xaa, 0x4b, 0x7c, 0x49, 0x7b, 0xcd, 0x93, 0xcd, 0xe9, 0x75, 0xe1, 0xec, 0xdc, 0xba, 0x50,
	0xed, 0xd6, 0x85, 0x1a, 0x34, 0x7d, 0x75, 0x62, 0x3d, 0x47, 0x7e, 0x03, 0x6a, 0x52, 0xa4, 0x9c,
	0xff, 0xcc, 0xc9, 0x92, 0xe5, 0xee, 0xa2, 0x84, 0x17, 0xb0, 0x72, 0xcb, 0x06, 0xd0, 0x16, 0x74,
	0x32, 0x23, 0x18, 0x0c, 0x09, 0x76, 0x48, 0x24, 0xbd, 0x65, 0x25, 0x8b, 0xbf, 0xe5, 0x61, 0x74,
	0x17, 0x1a, 0x59, 0x48, 0x9e, 0xf3, 0x6f, 0xe0, 0xd5, 0x39, 0x3c, 0xa2, 0x91, 0x6b, 0x0e, 0x27,
	0x21, 0x89, 0x7c, 0xe2, 0xb8, 0x24, 0x32, 0xaf, 0xf0, 0x65, 0xe4, 0xd9, 0xe2, 0xdb, 0x1d, 0xcb,
	0x8a, 0x2f, 0x4c, 0xd7, 0x63, 0xc3, 0xf1, 0x65, 0xb2, 0xec, 0x2b, 0xe0, 0xbe, 0x00, 0x6f, 0x0b,
	0xf0, 0xb6, 0x4b, 0xe5, 0x3f, 0x04, 0x97, 0x35, 0x1e, 0x79, 0xfe, 0x27, 0x00, 0x00, 0xff, 0xff,
	0x96, 0xdf, 0xe6, 0x5b, 0x47, 0x08, 0x00, 0x00,
}
