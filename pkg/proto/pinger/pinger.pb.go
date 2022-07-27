// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0-devel
// 	protoc        (unknown)
// source: pinger/pinger.proto

package pinger

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

type PingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PingRequest) Reset() {
	*x = PingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pinger_pinger_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingRequest) ProtoMessage() {}

func (x *PingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pinger_pinger_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingRequest.ProtoReflect.Descriptor instead.
func (*PingRequest) Descriptor() ([]byte, []int) {
	return file_pinger_pinger_proto_rawDescGZIP(), []int{0}
}

type PingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PingResponse) Reset() {
	*x = PingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pinger_pinger_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingResponse) ProtoMessage() {}

func (x *PingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pinger_pinger_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingResponse.ProtoReflect.Descriptor instead.
func (*PingResponse) Descriptor() ([]byte, []int) {
	return file_pinger_pinger_proto_rawDescGZIP(), []int{1}
}

var File_pinger_pinger_proto protoreflect.FileDescriptor

var file_pinger_pinger_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x2f, 0x70, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x22, 0x0d, 0x0a,
	0x0b, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x0e, 0x0a, 0x0c,
	0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x5f, 0x0a, 0x0a,
	0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x42, 0x0b, 0x50, 0x69, 0x6e, 0x67,
	0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x0c, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x70, 0x69, 0x6e, 0x67, 0x65, 0x72, 0xa2, 0x02, 0x03, 0x50, 0x58, 0x58, 0xaa, 0x02, 0x06,
	0x50, 0x69, 0x6e, 0x67, 0x65, 0x72, 0xca, 0x02, 0x06, 0x50, 0x69, 0x6e, 0x67, 0x65, 0x72, 0xe2,
	0x02, 0x12, 0x50, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x06, 0x50, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pinger_pinger_proto_rawDescOnce sync.Once
	file_pinger_pinger_proto_rawDescData = file_pinger_pinger_proto_rawDesc
)

func file_pinger_pinger_proto_rawDescGZIP() []byte {
	file_pinger_pinger_proto_rawDescOnce.Do(func() {
		file_pinger_pinger_proto_rawDescData = protoimpl.X.CompressGZIP(file_pinger_pinger_proto_rawDescData)
	})
	return file_pinger_pinger_proto_rawDescData
}

var file_pinger_pinger_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pinger_pinger_proto_goTypes = []interface{}{
	(*PingRequest)(nil),  // 0: pinger.PingRequest
	(*PingResponse)(nil), // 1: pinger.PingResponse
}
var file_pinger_pinger_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pinger_pinger_proto_init() }
func file_pinger_pinger_proto_init() {
	if File_pinger_pinger_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pinger_pinger_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingRequest); i {
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
		file_pinger_pinger_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingResponse); i {
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
			RawDescriptor: file_pinger_pinger_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pinger_pinger_proto_goTypes,
		DependencyIndexes: file_pinger_pinger_proto_depIdxs,
		MessageInfos:      file_pinger_pinger_proto_msgTypes,
	}.Build()
	File_pinger_pinger_proto = out.File
	file_pinger_pinger_proto_rawDesc = nil
	file_pinger_pinger_proto_goTypes = nil
	file_pinger_pinger_proto_depIdxs = nil
}
