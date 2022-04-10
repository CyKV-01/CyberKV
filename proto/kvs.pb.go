// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.6.1
// source: kvs.proto

package proto

import (
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

type AssignSlotRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SlotID int32 `protobuf:"varint,1,opt,name=slotID,proto3" json:"slotID,omitempty"`
}

func (x *AssignSlotRequest) Reset() {
	*x = AssignSlotRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kvs_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssignSlotRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssignSlotRequest) ProtoMessage() {}

func (x *AssignSlotRequest) ProtoReflect() protoreflect.Message {
	mi := &file_kvs_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssignSlotRequest.ProtoReflect.Descriptor instead.
func (*AssignSlotRequest) Descriptor() ([]byte, []int) {
	return file_kvs_proto_rawDescGZIP(), []int{0}
}

func (x *AssignSlotRequest) GetSlotID() int32 {
	if x != nil {
		return x.SlotID
	}
	return 0
}

type AssignSlotResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AssignSlotResponse) Reset() {
	*x = AssignSlotResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kvs_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssignSlotResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssignSlotResponse) ProtoMessage() {}

func (x *AssignSlotResponse) ProtoReflect() protoreflect.Message {
	mi := &file_kvs_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssignSlotResponse.ProtoReflect.Descriptor instead.
func (*AssignSlotResponse) Descriptor() ([]byte, []int) {
	return file_kvs_proto_rawDescGZIP(), []int{1}
}

type ReadOption struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConsistencyLevel int32 `protobuf:"varint,1,opt,name=consistency_level,json=consistencyLevel,proto3" json:"consistency_level,omitempty"`
}

func (x *ReadOption) Reset() {
	*x = ReadOption{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kvs_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadOption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadOption) ProtoMessage() {}

func (x *ReadOption) ProtoReflect() protoreflect.Message {
	mi := &file_kvs_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadOption.ProtoReflect.Descriptor instead.
func (*ReadOption) Descriptor() ([]byte, []int) {
	return file_kvs_proto_rawDescGZIP(), []int{2}
}

func (x *ReadOption) GetConsistencyLevel() int32 {
	if x != nil {
		return x.ConsistencyLevel
	}
	return 0
}

type ReadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key    string      `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Ts     uint64      `protobuf:"varint,2,opt,name=ts,proto3" json:"ts,omitempty"`
	Option *ReadOption `protobuf:"bytes,3,opt,name=option,proto3" json:"option,omitempty"`
	Info   []*NodeInfo `protobuf:"bytes,5,rep,name=info,proto3" json:"info,omitempty"`
}

func (x *ReadRequest) Reset() {
	*x = ReadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kvs_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadRequest) ProtoMessage() {}

func (x *ReadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_kvs_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadRequest.ProtoReflect.Descriptor instead.
func (*ReadRequest) Descriptor() ([]byte, []int) {
	return file_kvs_proto_rawDescGZIP(), []int{3}
}

func (x *ReadRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *ReadRequest) GetTs() uint64 {
	if x != nil {
		return x.Ts
	}
	return 0
}

func (x *ReadRequest) GetOption() *ReadOption {
	if x != nil {
		return x.Option
	}
	return nil
}

func (x *ReadRequest) GetInfo() []*NodeInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

type ReadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value  string  `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	Ts     uint64  `protobuf:"varint,2,opt,name=ts,proto3" json:"ts,omitempty"`
	Status *Status `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *ReadResponse) Reset() {
	*x = ReadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kvs_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadResponse) ProtoMessage() {}

func (x *ReadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_kvs_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadResponse.ProtoReflect.Descriptor instead.
func (*ReadResponse) Descriptor() ([]byte, []int) {
	return file_kvs_proto_rawDescGZIP(), []int{4}
}

func (x *ReadResponse) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *ReadResponse) GetTs() uint64 {
	if x != nil {
		return x.Ts
	}
	return 0
}

func (x *ReadResponse) GetStatus() *Status {
	if x != nil {
		return x.Status
	}
	return nil
}

type WriteOption struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sync bool `protobuf:"varint,1,opt,name=sync,proto3" json:"sync,omitempty"`
}

func (x *WriteOption) Reset() {
	*x = WriteOption{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kvs_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WriteOption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WriteOption) ProtoMessage() {}

func (x *WriteOption) ProtoReflect() protoreflect.Message {
	mi := &file_kvs_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WriteOption.ProtoReflect.Descriptor instead.
func (*WriteOption) Descriptor() ([]byte, []int) {
	return file_kvs_proto_rawDescGZIP(), []int{5}
}

func (x *WriteOption) GetSync() bool {
	if x != nil {
		return x.Sync
	}
	return false
}

type WriteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key      string       `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value    string       `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	IsRemove bool         `protobuf:"varint,3,opt,name=is_remove,json=isRemove,proto3" json:"is_remove,omitempty"`
	Ts       uint64       `protobuf:"varint,4,opt,name=ts,proto3" json:"ts,omitempty"`
	Option   *WriteOption `protobuf:"bytes,5,opt,name=option,proto3" json:"option,omitempty"`
	Info     []*NodeInfo  `protobuf:"bytes,7,rep,name=info,proto3" json:"info,omitempty"`
}

func (x *WriteRequest) Reset() {
	*x = WriteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kvs_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WriteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WriteRequest) ProtoMessage() {}

func (x *WriteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_kvs_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WriteRequest.ProtoReflect.Descriptor instead.
func (*WriteRequest) Descriptor() ([]byte, []int) {
	return file_kvs_proto_rawDescGZIP(), []int{6}
}

func (x *WriteRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *WriteRequest) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *WriteRequest) GetIsRemove() bool {
	if x != nil {
		return x.IsRemove
	}
	return false
}

func (x *WriteRequest) GetTs() uint64 {
	if x != nil {
		return x.Ts
	}
	return 0
}

func (x *WriteRequest) GetOption() *WriteOption {
	if x != nil {
		return x.Option
	}
	return nil
}

func (x *WriteRequest) GetInfo() []*NodeInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

type WriteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status *Status `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *WriteResponse) Reset() {
	*x = WriteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kvs_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WriteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WriteResponse) ProtoMessage() {}

func (x *WriteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_kvs_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WriteResponse.ProtoReflect.Descriptor instead.
func (*WriteResponse) Descriptor() ([]byte, []int) {
	return file_kvs_proto_rawDescGZIP(), []int{7}
}

func (x *WriteResponse) GetStatus() *Status {
	if x != nil {
		return x.Status
	}
	return nil
}

var File_kvs_proto protoreflect.FileDescriptor

var file_kvs_proto_rawDesc = []byte{
	0x0a, 0x09, 0x6b, 0x76, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x6b, 0x76, 0x73,
	0x1a, 0x0c, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0a,
	0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2b, 0x0a, 0x11, 0x41, 0x73,
	0x73, 0x69, 0x67, 0x6e, 0x53, 0x6c, 0x6f, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x6c, 0x6f, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x06, 0x73, 0x6c, 0x6f, 0x74, 0x49, 0x44, 0x22, 0x14, 0x0a, 0x12, 0x41, 0x73, 0x73, 0x69, 0x67,
	0x6e, 0x53, 0x6c, 0x6f, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x39, 0x0a,
	0x0a, 0x52, 0x65, 0x61, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2b, 0x0a, 0x11, 0x63,
	0x6f, 0x6e, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x6c, 0x65, 0x76, 0x65, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x10, 0x63, 0x6f, 0x6e, 0x73, 0x69, 0x73, 0x74, 0x65,
	0x6e, 0x63, 0x79, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x22, 0x7c, 0x0a, 0x0b, 0x52, 0x65, 0x61, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x74, 0x73, 0x12, 0x27, 0x0a, 0x06, 0x6f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6b, 0x76, 0x73, 0x2e,
	0x52, 0x65, 0x61, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x6f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0e, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x22, 0x5c, 0x0a, 0x0c, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x0e, 0x0a, 0x02,
	0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x74, 0x73, 0x12, 0x26, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x22, 0x21, 0x0a, 0x0b, 0x57, 0x72, 0x69, 0x74, 0x65, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x79, 0x6e, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x04, 0x73, 0x79, 0x6e, 0x63, 0x22, 0xb1, 0x01, 0x0a, 0x0c, 0x57, 0x72, 0x69, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x72, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x74, 0x73, 0x12, 0x28, 0x0a,
	0x06, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e,
	0x6b, 0x76, 0x73, 0x2e, 0x57, 0x72, 0x69, 0x74, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x06, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18,
	0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x4e, 0x6f, 0x64,
	0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x22, 0x37, 0x0a, 0x0d, 0x57,
	0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x32, 0xdc, 0x01, 0x0a, 0x08, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x3f, 0x0a, 0x0a, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x53, 0x6c, 0x6f, 0x74, 0x12,
	0x16, 0x2e, 0x6b, 0x76, 0x73, 0x2e, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x53, 0x6c, 0x6f, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x6b, 0x76, 0x73, 0x2e, 0x41, 0x73,
	0x73, 0x69, 0x67, 0x6e, 0x53, 0x6c, 0x6f, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x2c, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x10, 0x2e, 0x6b, 0x76, 0x73, 0x2e,
	0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x6b, 0x76,
	0x73, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x2e, 0x0a, 0x03, 0x53, 0x65, 0x74, 0x12, 0x11, 0x2e, 0x6b, 0x76, 0x73, 0x2e, 0x57, 0x72,
	0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x6b, 0x76, 0x73,
	0x2e, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x31, 0x0a, 0x06, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x12, 0x11, 0x2e, 0x6b, 0x76, 0x73,
	0x2e, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e,
	0x6b, 0x76, 0x73, 0x2e, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x22, 0x5a, 0x20, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x43, 0x79, 0x4b, 0x56, 0x2d, 0x30, 0x31, 0x2f, 0x43, 0x79, 0x62, 0x65, 0x72, 0x4b,
	0x56, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_kvs_proto_rawDescOnce sync.Once
	file_kvs_proto_rawDescData = file_kvs_proto_rawDesc
)

func file_kvs_proto_rawDescGZIP() []byte {
	file_kvs_proto_rawDescOnce.Do(func() {
		file_kvs_proto_rawDescData = protoimpl.X.CompressGZIP(file_kvs_proto_rawDescData)
	})
	return file_kvs_proto_rawDescData
}

var file_kvs_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_kvs_proto_goTypes = []interface{}{
	(*AssignSlotRequest)(nil),  // 0: kvs.AssignSlotRequest
	(*AssignSlotResponse)(nil), // 1: kvs.AssignSlotResponse
	(*ReadOption)(nil),         // 2: kvs.ReadOption
	(*ReadRequest)(nil),        // 3: kvs.ReadRequest
	(*ReadResponse)(nil),       // 4: kvs.ReadResponse
	(*WriteOption)(nil),        // 5: kvs.WriteOption
	(*WriteRequest)(nil),       // 6: kvs.WriteRequest
	(*WriteResponse)(nil),      // 7: kvs.WriteResponse
	(*NodeInfo)(nil),           // 8: node.NodeInfo
	(*Status)(nil),             // 9: status.Status
}
var file_kvs_proto_depIdxs = []int32{
	2,  // 0: kvs.ReadRequest.option:type_name -> kvs.ReadOption
	8,  // 1: kvs.ReadRequest.info:type_name -> node.NodeInfo
	9,  // 2: kvs.ReadResponse.status:type_name -> status.Status
	5,  // 3: kvs.WriteRequest.option:type_name -> kvs.WriteOption
	8,  // 4: kvs.WriteRequest.info:type_name -> node.NodeInfo
	9,  // 5: kvs.WriteResponse.status:type_name -> status.Status
	0,  // 6: kvs.KeyValue.AssignSlot:input_type -> kvs.AssignSlotRequest
	3,  // 7: kvs.KeyValue.Get:input_type -> kvs.ReadRequest
	6,  // 8: kvs.KeyValue.Set:input_type -> kvs.WriteRequest
	6,  // 9: kvs.KeyValue.Remove:input_type -> kvs.WriteRequest
	1,  // 10: kvs.KeyValue.AssignSlot:output_type -> kvs.AssignSlotResponse
	4,  // 11: kvs.KeyValue.Get:output_type -> kvs.ReadResponse
	7,  // 12: kvs.KeyValue.Set:output_type -> kvs.WriteResponse
	7,  // 13: kvs.KeyValue.Remove:output_type -> kvs.WriteResponse
	10, // [10:14] is the sub-list for method output_type
	6,  // [6:10] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_kvs_proto_init() }
func file_kvs_proto_init() {
	if File_kvs_proto != nil {
		return
	}
	file_status_proto_init()
	file_node_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_kvs_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssignSlotRequest); i {
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
		file_kvs_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssignSlotResponse); i {
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
		file_kvs_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadOption); i {
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
		file_kvs_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadRequest); i {
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
		file_kvs_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadResponse); i {
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
		file_kvs_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WriteOption); i {
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
		file_kvs_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WriteRequest); i {
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
		file_kvs_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WriteResponse); i {
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
			RawDescriptor: file_kvs_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_kvs_proto_goTypes,
		DependencyIndexes: file_kvs_proto_depIdxs,
		MessageInfos:      file_kvs_proto_msgTypes,
	}.Build()
	File_kvs_proto = out.File
	file_kvs_proto_rawDesc = nil
	file_kvs_proto_goTypes = nil
	file_kvs_proto_depIdxs = nil
}
