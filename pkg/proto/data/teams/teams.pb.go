// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: data/teams/teams.proto

package teams

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

type Team struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      int32             `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name    string            `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Address string            `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"`
	Token   string            `protobuf:"bytes,4,opt,name=token,proto3" json:"token,omitempty"`
	Labels  map[string]string `protobuf:"bytes,5,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Team) Reset() {
	*x = Team{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_teams_teams_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Team) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Team) ProtoMessage() {}

func (x *Team) ProtoReflect() protoreflect.Message {
	mi := &file_data_teams_teams_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Team.ProtoReflect.Descriptor instead.
func (*Team) Descriptor() ([]byte, []int) {
	return file_data_teams_teams_proto_rawDescGZIP(), []int{0}
}

func (x *Team) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Team) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Team) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *Team) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *Team) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

type ListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LastUpdate int64 `protobuf:"varint,1,opt,name=last_update,json=lastUpdate,proto3" json:"last_update,omitempty"`
}

func (x *ListRequest) Reset() {
	*x = ListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_teams_teams_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRequest) ProtoMessage() {}

func (x *ListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_data_teams_teams_proto_msgTypes[1]
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
	return file_data_teams_teams_proto_rawDescGZIP(), []int{1}
}

func (x *ListRequest) GetLastUpdate() int64 {
	if x != nil {
		return x.LastUpdate
	}
	return 0
}

type ListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Teams      []*Team `protobuf:"bytes,1,rep,name=teams,proto3" json:"teams,omitempty"`
	LastUpdate int64   `protobuf:"varint,2,opt,name=last_update,json=lastUpdate,proto3" json:"last_update,omitempty"`
}

func (x *ListResponse) Reset() {
	*x = ListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_teams_teams_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListResponse) ProtoMessage() {}

func (x *ListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_data_teams_teams_proto_msgTypes[2]
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
	return file_data_teams_teams_proto_rawDescGZIP(), []int{2}
}

func (x *ListResponse) GetTeams() []*Team {
	if x != nil {
		return x.Teams
	}
	return nil
}

func (x *ListResponse) GetLastUpdate() int64 {
	if x != nil {
		return x.LastUpdate
	}
	return 0
}

type CreateBatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Teams []*Team `protobuf:"bytes,1,rep,name=teams,proto3" json:"teams,omitempty"`
}

func (x *CreateBatchRequest) Reset() {
	*x = CreateBatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_teams_teams_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateBatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateBatchRequest) ProtoMessage() {}

func (x *CreateBatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_data_teams_teams_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateBatchRequest.ProtoReflect.Descriptor instead.
func (*CreateBatchRequest) Descriptor() ([]byte, []int) {
	return file_data_teams_teams_proto_rawDescGZIP(), []int{3}
}

func (x *CreateBatchRequest) GetTeams() []*Team {
	if x != nil {
		return x.Teams
	}
	return nil
}

type CreateBatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Teams []*Team `protobuf:"bytes,1,rep,name=teams,proto3" json:"teams,omitempty"`
}

func (x *CreateBatchResponse) Reset() {
	*x = CreateBatchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_teams_teams_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateBatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateBatchResponse) ProtoMessage() {}

func (x *CreateBatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_data_teams_teams_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateBatchResponse.ProtoReflect.Descriptor instead.
func (*CreateBatchResponse) Descriptor() ([]byte, []int) {
	return file_data_teams_teams_proto_rawDescGZIP(), []int{4}
}

func (x *CreateBatchResponse) GetTeams() []*Team {
	if x != nil {
		return x.Teams
	}
	return nil
}

var File_data_teams_teams_proto protoreflect.FileDescriptor

var file_data_teams_teams_proto_rawDesc = []byte{
	0x0a, 0x16, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x2f, 0x74, 0x65, 0x61,
	0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x74,
	0x65, 0x61, 0x6d, 0x73, 0x22, 0xcb, 0x01, 0x0a, 0x04, 0x54, 0x65, 0x61, 0x6d, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x12, 0x34, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x2e, 0x54,
	0x65, 0x61, 0x6d, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62, 0x65, 0x6c,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x2e, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x6c, 0x61, 0x73, 0x74, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x22, 0x57, 0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x26, 0x0a, 0x05, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x2e, 0x54,
	0x65, 0x61, 0x6d, 0x52, 0x05, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x6c, 0x61,
	0x73, 0x74, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0a, 0x6c, 0x61, 0x73, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x22, 0x3c, 0x0a, 0x12, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x26, 0x0a, 0x05, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x2e, 0x54, 0x65,
	0x61, 0x6d, 0x52, 0x05, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x22, 0x3d, 0x0a, 0x13, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x26, 0x0a, 0x05, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x2e, 0x54, 0x65, 0x61,
	0x6d, 0x52, 0x05, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x42, 0x99, 0x01, 0x0a, 0x0e, 0x63, 0x6f, 0x6d,
	0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x42, 0x0a, 0x54, 0x65, 0x61,
	0x6d, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x34, 0x74, 0x2d, 0x62, 0x75, 0x74, 0x2d, 0x73, 0x34,
	0x64, 0x2f, 0x66, 0x61, 0x73, 0x74, 0x61, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x74, 0x65, 0x61, 0x6d, 0x73, 0xa2, 0x02, 0x03,
	0x44, 0x54, 0x58, 0xaa, 0x02, 0x0a, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x54, 0x65, 0x61, 0x6d, 0x73,
	0xca, 0x02, 0x0a, 0x44, 0x61, 0x74, 0x61, 0x5c, 0x54, 0x65, 0x61, 0x6d, 0x73, 0xe2, 0x02, 0x16,
	0x44, 0x61, 0x74, 0x61, 0x5c, 0x54, 0x65, 0x61, 0x6d, 0x73, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0b, 0x44, 0x61, 0x74, 0x61, 0x3a, 0x3a, 0x54,
	0x65, 0x61, 0x6d, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_data_teams_teams_proto_rawDescOnce sync.Once
	file_data_teams_teams_proto_rawDescData = file_data_teams_teams_proto_rawDesc
)

func file_data_teams_teams_proto_rawDescGZIP() []byte {
	file_data_teams_teams_proto_rawDescOnce.Do(func() {
		file_data_teams_teams_proto_rawDescData = protoimpl.X.CompressGZIP(file_data_teams_teams_proto_rawDescData)
	})
	return file_data_teams_teams_proto_rawDescData
}

var file_data_teams_teams_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_data_teams_teams_proto_goTypes = []interface{}{
	(*Team)(nil),                // 0: data.teams.Team
	(*ListRequest)(nil),         // 1: data.teams.ListRequest
	(*ListResponse)(nil),        // 2: data.teams.ListResponse
	(*CreateBatchRequest)(nil),  // 3: data.teams.CreateBatchRequest
	(*CreateBatchResponse)(nil), // 4: data.teams.CreateBatchResponse
	nil,                         // 5: data.teams.Team.LabelsEntry
}
var file_data_teams_teams_proto_depIdxs = []int32{
	5, // 0: data.teams.Team.labels:type_name -> data.teams.Team.LabelsEntry
	0, // 1: data.teams.ListResponse.teams:type_name -> data.teams.Team
	0, // 2: data.teams.CreateBatchRequest.teams:type_name -> data.teams.Team
	0, // 3: data.teams.CreateBatchResponse.teams:type_name -> data.teams.Team
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_data_teams_teams_proto_init() }
func file_data_teams_teams_proto_init() {
	if File_data_teams_teams_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_data_teams_teams_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Team); i {
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
		file_data_teams_teams_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_data_teams_teams_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_data_teams_teams_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateBatchRequest); i {
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
		file_data_teams_teams_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateBatchResponse); i {
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
			RawDescriptor: file_data_teams_teams_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_data_teams_teams_proto_goTypes,
		DependencyIndexes: file_data_teams_teams_proto_depIdxs,
		MessageInfos:      file_data_teams_teams_proto_msgTypes,
	}.Build()
	File_data_teams_teams_proto = out.File
	file_data_teams_teams_proto_rawDesc = nil
	file_data_teams_teams_proto_goTypes = nil
	file_data_teams_teams_proto_depIdxs = nil
}