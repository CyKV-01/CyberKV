// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.6.1
// source: version.proto

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

type VersionSet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Current uint64 `protobuf:"varint,1,opt,name=current,proto3" json:"current,omitempty"`
	Last    uint64 `protobuf:"varint,2,opt,name=last,proto3" json:"last,omitempty"`
}

func (x *VersionSet) Reset() {
	*x = VersionSet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_version_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VersionSet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VersionSet) ProtoMessage() {}

func (x *VersionSet) ProtoReflect() protoreflect.Message {
	mi := &file_version_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VersionSet.ProtoReflect.Descriptor instead.
func (*VersionSet) Descriptor() ([]byte, []int) {
	return file_version_proto_rawDescGZIP(), []int{0}
}

func (x *VersionSet) GetCurrent() uint64 {
	if x != nil {
		return x.Current
	}
	return 0
}

func (x *VersionSet) GetLast() uint64 {
	if x != nil {
		return x.Last
	}
	return 0
}

type Version struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VersionId uint64                  `protobuf:"varint,1,opt,name=version_id,json=versionId,proto3" json:"version_id,omitempty"`
	Sstables  map[int32]*SSTableLevel `protobuf:"bytes,2,rep,name=sstables,proto3" json:"sstables,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Version) Reset() {
	*x = Version{}
	if protoimpl.UnsafeEnabled {
		mi := &file_version_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Version) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Version) ProtoMessage() {}

func (x *Version) ProtoReflect() protoreflect.Message {
	mi := &file_version_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Version.ProtoReflect.Descriptor instead.
func (*Version) Descriptor() ([]byte, []int) {
	return file_version_proto_rawDescGZIP(), []int{1}
}

func (x *Version) GetVersionId() uint64 {
	if x != nil {
		return x.VersionId
	}
	return 0
}

func (x *Version) GetSstables() map[int32]*SSTableLevel {
	if x != nil {
		return x.Sstables
	}
	return nil
}

type SSTableLevel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Level    int32    `protobuf:"varint,1,opt,name=level,proto3" json:"level,omitempty"`
	Sstables []string `protobuf:"bytes,2,rep,name=sstables,proto3" json:"sstables,omitempty"`
}

func (x *SSTableLevel) Reset() {
	*x = SSTableLevel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_version_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SSTableLevel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SSTableLevel) ProtoMessage() {}

func (x *SSTableLevel) ProtoReflect() protoreflect.Message {
	mi := &file_version_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SSTableLevel.ProtoReflect.Descriptor instead.
func (*SSTableLevel) Descriptor() ([]byte, []int) {
	return file_version_proto_rawDescGZIP(), []int{2}
}

func (x *SSTableLevel) GetLevel() int32 {
	if x != nil {
		return x.Level
	}
	return 0
}

func (x *SSTableLevel) GetSstables() []string {
	if x != nil {
		return x.Sstables
	}
	return nil
}

var File_version_proto protoreflect.FileDescriptor

var file_version_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x3a, 0x0a, 0x0a, 0x56, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x53, 0x65, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x6c, 0x61, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04,
	0x6c, 0x61, 0x73, 0x74, 0x22, 0xb8, 0x01, 0x0a, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x1d, 0x0a, 0x0a, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12,
	0x3a, 0x0a, 0x08, 0x73, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x1e, 0x2e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x56, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x08, 0x73, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x1a, 0x52, 0x0a, 0x0d, 0x53,
	0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2b,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x53, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x4c,
	0x65, 0x76, 0x65, 0x6c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22,
	0x40, 0x0a, 0x0c, 0x53, 0x53, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x12,
	0x14, 0x0a, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05,
	0x6c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x73, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x73, 0x42, 0x22, 0x5a, 0x20, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x43, 0x79, 0x4b, 0x56, 0x2d, 0x30, 0x31, 0x2f, 0x43, 0x79, 0x62, 0x65, 0x72, 0x4b, 0x56, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_version_proto_rawDescOnce sync.Once
	file_version_proto_rawDescData = file_version_proto_rawDesc
)

func file_version_proto_rawDescGZIP() []byte {
	file_version_proto_rawDescOnce.Do(func() {
		file_version_proto_rawDescData = protoimpl.X.CompressGZIP(file_version_proto_rawDescData)
	})
	return file_version_proto_rawDescData
}

var file_version_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_version_proto_goTypes = []interface{}{
	(*VersionSet)(nil),   // 0: version.VersionSet
	(*Version)(nil),      // 1: version.Version
	(*SSTableLevel)(nil), // 2: version.SSTableLevel
	nil,                  // 3: version.Version.SstablesEntry
}
var file_version_proto_depIdxs = []int32{
	3, // 0: version.Version.sstables:type_name -> version.Version.SstablesEntry
	2, // 1: version.Version.SstablesEntry.value:type_name -> version.SSTableLevel
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_version_proto_init() }
func file_version_proto_init() {
	if File_version_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_version_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VersionSet); i {
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
		file_version_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Version); i {
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
		file_version_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SSTableLevel); i {
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
			RawDescriptor: file_version_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_version_proto_goTypes,
		DependencyIndexes: file_version_proto_depIdxs,
		MessageInfos:      file_version_proto_msgTypes,
	}.Build()
	File_version_proto = out.File
	file_version_proto_rawDesc = nil
	file_version_proto_goTypes = nil
	file_version_proto_depIdxs = nil
}
