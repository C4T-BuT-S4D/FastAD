// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: data/teams/teams_service.proto

package teams

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

var File_data_teams_teams_service_proto protoreflect.FileDescriptor

var file_data_teams_teams_service_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x2f, 0x74, 0x65, 0x61,
	0x6d, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0a, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x64, 0x61, 0x74, 0x61,
	0x2f, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x2f, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x32, 0xad, 0x01, 0x0a, 0x0c, 0x54, 0x65, 0x61, 0x6d, 0x73, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x4d, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x17, 0x2e, 0x64, 0x61,
	0x74, 0x61, 0x2e, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x74, 0x65, 0x61, 0x6d,
	0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x12,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0c, 0x12, 0x0a, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x74, 0x65, 0x61,
	0x6d, 0x73, 0x12, 0x4e, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x61, 0x74, 0x63,
	0x68, 0x12, 0x1e, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1f, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0xa0, 0x01, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e,
	0x74, 0x65, 0x61, 0x6d, 0x73, 0x42, 0x11, 0x54, 0x65, 0x61, 0x6d, 0x73, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x34, 0x74, 0x2d, 0x62, 0x75, 0x74, 0x2d, 0x73,
	0x34, 0x64, 0x2f, 0x66, 0x61, 0x73, 0x74, 0x61, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x74, 0x65, 0x61, 0x6d, 0x73, 0xa2, 0x02,
	0x03, 0x44, 0x54, 0x58, 0xaa, 0x02, 0x0a, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x54, 0x65, 0x61, 0x6d,
	0x73, 0xca, 0x02, 0x0a, 0x44, 0x61, 0x74, 0x61, 0x5c, 0x54, 0x65, 0x61, 0x6d, 0x73, 0xe2, 0x02,
	0x16, 0x44, 0x61, 0x74, 0x61, 0x5c, 0x54, 0x65, 0x61, 0x6d, 0x73, 0x5c, 0x47, 0x50, 0x42, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0b, 0x44, 0x61, 0x74, 0x61, 0x3a, 0x3a,
	0x54, 0x65, 0x61, 0x6d, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_data_teams_teams_service_proto_goTypes = []any{
	(*ListRequest)(nil),         // 0: data.teams.ListRequest
	(*CreateBatchRequest)(nil),  // 1: data.teams.CreateBatchRequest
	(*ListResponse)(nil),        // 2: data.teams.ListResponse
	(*CreateBatchResponse)(nil), // 3: data.teams.CreateBatchResponse
}
var file_data_teams_teams_service_proto_depIdxs = []int32{
	0, // 0: data.teams.TeamsService.List:input_type -> data.teams.ListRequest
	1, // 1: data.teams.TeamsService.CreateBatch:input_type -> data.teams.CreateBatchRequest
	2, // 2: data.teams.TeamsService.List:output_type -> data.teams.ListResponse
	3, // 3: data.teams.TeamsService.CreateBatch:output_type -> data.teams.CreateBatchResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_data_teams_teams_service_proto_init() }
func file_data_teams_teams_service_proto_init() {
	if File_data_teams_teams_service_proto != nil {
		return
	}
	file_data_teams_teams_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_data_teams_teams_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_data_teams_teams_service_proto_goTypes,
		DependencyIndexes: file_data_teams_teams_service_proto_depIdxs,
	}.Build()
	File_data_teams_teams_service_proto = out.File
	file_data_teams_teams_service_proto_rawDesc = nil
	file_data_teams_teams_service_proto_goTypes = nil
	file_data_teams_teams_service_proto_depIdxs = nil
}
