// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.12.1
// source: enc.proto

package google

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

type R1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uint64 uint64 `protobuf:"varint,4,opt,name=uint64,proto3" json:"uint64,omitempty"`
	Uint32 uint32 `protobuf:"varint,200,opt,name=uint32,proto3" json:"uint32,omitempty"`
	Int64  int64  `protobuf:"varint,201,opt,name=int64,proto3" json:"int64,omitempty"`
	Int32  int32  `protobuf:"varint,202,opt,name=int32,proto3" json:"int32,omitempty"`
	Sint64 int64  `protobuf:"zigzag64,203,opt,name=sint64,proto3" json:"sint64,omitempty"`
	Sint32 int32  `protobuf:"zigzag32,204,opt,name=sint32,proto3" json:"sint32,omitempty"`
}

func (x *R1) Reset() {
	*x = R1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_enc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *R1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*R1) ProtoMessage() {}

func (x *R1) ProtoReflect() protoreflect.Message {
	mi := &file_enc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use R1.ProtoReflect.Descriptor instead.
func (*R1) Descriptor() ([]byte, []int) {
	return file_enc_proto_rawDescGZIP(), []int{0}
}

func (x *R1) GetUint64() uint64 {
	if x != nil {
		return x.Uint64
	}
	return 0
}

func (x *R1) GetUint32() uint32 {
	if x != nil {
		return x.Uint32
	}
	return 0
}

func (x *R1) GetInt64() int64 {
	if x != nil {
		return x.Int64
	}
	return 0
}

func (x *R1) GetInt32() int32 {
	if x != nil {
		return x.Int32
	}
	return 0
}

func (x *R1) GetSint64() int64 {
	if x != nil {
		return x.Sint64
	}
	return 0
}

func (x *R1) GetSint32() int32 {
	if x != nil {
		return x.Sint32
	}
	return 0
}

var File_enc_proto protoreflect.FileDescriptor

var file_enc_proto_rawDesc = []byte{
	0x0a, 0x09, 0x65, 0x6e, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x74, 0x65, 0x73,
	0x74, 0x22, 0x95, 0x01, 0x0a, 0x02, 0x52, 0x31, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x69, 0x6e, 0x74,
	0x36, 0x34, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x69, 0x6e, 0x74, 0x36, 0x34,
	0x12, 0x17, 0x0a, 0x06, 0x75, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x18, 0xc8, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x06, 0x75, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x12, 0x15, 0x0a, 0x05, 0x69, 0x6e, 0x74,
	0x36, 0x34, 0x18, 0xc9, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x69, 0x6e, 0x74, 0x36, 0x34,
	0x12, 0x15, 0x0a, 0x05, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x18, 0xca, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x05, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x12, 0x17, 0x0a, 0x06, 0x73, 0x69, 0x6e, 0x74, 0x36,
	0x34, 0x18, 0xcb, 0x01, 0x20, 0x01, 0x28, 0x12, 0x52, 0x06, 0x73, 0x69, 0x6e, 0x74, 0x36, 0x34,
	0x12, 0x17, 0x0a, 0x06, 0x73, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x18, 0xcc, 0x01, 0x20, 0x01, 0x28,
	0x11, 0x52, 0x06, 0x73, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x42, 0x15, 0x5a, 0x13, 0x74, 0x65, 0x73,
	0x74, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2f, 0x74, 0x65, 0x73, 0x74,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_enc_proto_rawDescOnce sync.Once
	file_enc_proto_rawDescData = file_enc_proto_rawDesc
)

func file_enc_proto_rawDescGZIP() []byte {
	file_enc_proto_rawDescOnce.Do(func() {
		file_enc_proto_rawDescData = protoimpl.X.CompressGZIP(file_enc_proto_rawDescData)
	})
	return file_enc_proto_rawDescData
}

var file_enc_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_enc_proto_goTypes = []interface{}{
	(*R1)(nil), // 0: test.R1
}
var file_enc_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_enc_proto_init() }
func file_enc_proto_init() {
	if File_enc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_enc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*R1); i {
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
			RawDescriptor: file_enc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_enc_proto_goTypes,
		DependencyIndexes: file_enc_proto_depIdxs,
		MessageInfos:      file_enc_proto_msgTypes,
	}.Build()
	File_enc_proto = out.File
	file_enc_proto_rawDesc = nil
	file_enc_proto_goTypes = nil
	file_enc_proto_depIdxs = nil
}