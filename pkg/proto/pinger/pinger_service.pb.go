// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: pinger/pinger_service.proto

package pinger

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_pinger_pinger_service_proto protoreflect.FileDescriptor

var file_pinger_pinger_service_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x2f, 0x70, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70,
	0x69, 0x6e, 0x67, 0x65, 0x72, 0x1a, 0x13, 0x70, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x2f, 0x70, 0x69,
	0x6e, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x42, 0x0a, 0x0d, 0x50, 0x69,
	0x6e, 0x67, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x31, 0x0a, 0x04, 0x50,
	0x69, 0x6e, 0x67, 0x12, 0x13, 0x2e, 0x70, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x2e, 0x50, 0x69, 0x6e,
	0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x70, 0x69, 0x6e, 0x67, 0x65,
	0x72, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x88,
	0x01, 0x0a, 0x0a, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x42, 0x12, 0x50,
	0x69, 0x6e, 0x67, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x50, 0x01, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x63, 0x34, 0x74, 0x2d, 0x62, 0x75, 0x74, 0x2d, 0x73, 0x34, 0x64, 0x2f, 0x66, 0x61, 0x73, 0x74,
	0x61, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x69, 0x6e,
	0x67, 0x65, 0x72, 0xa2, 0x02, 0x03, 0x50, 0x58, 0x58, 0xaa, 0x02, 0x06, 0x50, 0x69, 0x6e, 0x67,
	0x65, 0x72, 0xca, 0x02, 0x06, 0x50, 0x69, 0x6e, 0x67, 0x65, 0x72, 0xe2, 0x02, 0x12, 0x50, 0x69,
	0x6e, 0x67, 0x65, 0x72, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0xea, 0x02, 0x06, 0x50, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var file_pinger_pinger_service_proto_goTypes = []any{
	(*PingRequest)(nil),  // 0: pinger.PingRequest
	(*PingResponse)(nil), // 1: pinger.PingResponse
}
var file_pinger_pinger_service_proto_depIdxs = []int32{
	0, // 0: pinger.PingerService.Ping:input_type -> pinger.PingRequest
	1, // 1: pinger.PingerService.Ping:output_type -> pinger.PingResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pinger_pinger_service_proto_init() }
func file_pinger_pinger_service_proto_init() {
	if File_pinger_pinger_service_proto != nil {
		return
	}
	file_pinger_pinger_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pinger_pinger_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pinger_pinger_service_proto_goTypes,
		DependencyIndexes: file_pinger_pinger_service_proto_depIdxs,
	}.Build()
	File_pinger_pinger_service_proto = out.File
	file_pinger_pinger_service_proto_rawDesc = nil
	file_pinger_pinger_service_proto_goTypes = nil
	file_pinger_pinger_service_proto_depIdxs = nil
}
