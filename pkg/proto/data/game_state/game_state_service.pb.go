// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        (unknown)
// source: data/game_state/game_state_service.proto

package game_state

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

var File_data_game_state_game_state_service_proto protoreflect.FileDescriptor

var file_data_game_state_game_state_service_proto_rawDesc = []byte{
	0x0a, 0x28, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x64, 0x61, 0x74, 0x61,
	0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x1a, 0x1c, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x64, 0x61, 0x74, 0x61, 0x2f,
	0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x5f,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xca, 0x02, 0x0a, 0x10,
	0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x59, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x1b, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x67,
	0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x67, 0x61, 0x6d, 0x65,
	0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x17, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11, 0x12, 0x0f, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x62, 0x0a, 0x06, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1e, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x67, 0x61, 0x6d,
	0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x67, 0x61, 0x6d,
	0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x17, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11, 0x22, 0x0f,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12,
	0x77, 0x0a, 0x0b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x23,
	0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x75, 0x6e,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x17, 0x22, 0x15, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x61,
	0x74, 0x65, 0x2f, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x42, 0xbe, 0x01, 0x0a, 0x13, 0x63, 0x6f, 0x6d,
	0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x42, 0x15, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x34, 0x74, 0x2d, 0x62, 0x75, 0x74, 0x2d, 0x73, 0x34,
	0x64, 0x2f, 0x66, 0x61, 0x73, 0x74, 0x61, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x61,
	0x74, 0x65, 0xa2, 0x02, 0x03, 0x44, 0x47, 0x58, 0xaa, 0x02, 0x0e, 0x44, 0x61, 0x74, 0x61, 0x2e,
	0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0xca, 0x02, 0x0e, 0x44, 0x61, 0x74, 0x61,
	0x5c, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0xe2, 0x02, 0x1a, 0x44, 0x61, 0x74,
	0x61, 0x5c, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x5c, 0x47, 0x50, 0x42, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0f, 0x44, 0x61, 0x74, 0x61, 0x3a, 0x3a,
	0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var file_data_game_state_game_state_service_proto_goTypes = []any{
	(*GetRequest)(nil),          // 0: data.game_state.GetRequest
	(*UpdateRequest)(nil),       // 1: data.game_state.UpdateRequest
	(*UpdateRoundRequest)(nil),  // 2: data.game_state.UpdateRoundRequest
	(*GetResponse)(nil),         // 3: data.game_state.GetResponse
	(*UpdateResponse)(nil),      // 4: data.game_state.UpdateResponse
	(*UpdateRoundResponse)(nil), // 5: data.game_state.UpdateRoundResponse
}
var file_data_game_state_game_state_service_proto_depIdxs = []int32{
	0, // 0: data.game_state.GameStateService.Get:input_type -> data.game_state.GetRequest
	1, // 1: data.game_state.GameStateService.Update:input_type -> data.game_state.UpdateRequest
	2, // 2: data.game_state.GameStateService.UpdateRound:input_type -> data.game_state.UpdateRoundRequest
	3, // 3: data.game_state.GameStateService.Get:output_type -> data.game_state.GetResponse
	4, // 4: data.game_state.GameStateService.Update:output_type -> data.game_state.UpdateResponse
	5, // 5: data.game_state.GameStateService.UpdateRound:output_type -> data.game_state.UpdateRoundResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_data_game_state_game_state_service_proto_init() }
func file_data_game_state_game_state_service_proto_init() {
	if File_data_game_state_game_state_service_proto != nil {
		return
	}
	file_data_game_state_game_state_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_data_game_state_game_state_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_data_game_state_game_state_service_proto_goTypes,
		DependencyIndexes: file_data_game_state_game_state_service_proto_depIdxs,
	}.Build()
	File_data_game_state_game_state_service_proto = out.File
	file_data_game_state_game_state_service_proto_rawDesc = nil
	file_data_game_state_game_state_service_proto_goTypes = nil
	file_data_game_state_game_state_service_proto_depIdxs = nil
}
