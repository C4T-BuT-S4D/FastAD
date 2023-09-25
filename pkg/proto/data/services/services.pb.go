// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: data/services/services.proto

package services

import (
	checker "github.com/c4t-but-s4d/fastad/pkg/proto/checker"
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

type Service struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           int32            `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name         string           `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Checker      *Service_Checker `protobuf:"bytes,3,opt,name=checker,proto3" json:"checker,omitempty"`
	DefaultScore float64          `protobuf:"fixed64,4,opt,name=default_score,json=defaultScore,proto3" json:"default_score,omitempty"`
	Gets         int32            `protobuf:"varint,5,opt,name=gets,proto3" json:"gets,omitempty"`
	Puts         int32            `protobuf:"varint,6,opt,name=puts,proto3" json:"puts,omitempty"`
}

func (x *Service) Reset() {
	*x = Service{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_services_services_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Service) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Service) ProtoMessage() {}

func (x *Service) ProtoReflect() protoreflect.Message {
	mi := &file_data_services_services_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Service.ProtoReflect.Descriptor instead.
func (*Service) Descriptor() ([]byte, []int) {
	return file_data_services_services_proto_rawDescGZIP(), []int{0}
}

func (x *Service) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Service) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Service) GetChecker() *Service_Checker {
	if x != nil {
		return x.Checker
	}
	return nil
}

func (x *Service) GetDefaultScore() float64 {
	if x != nil {
		return x.DefaultScore
	}
	return 0
}

func (x *Service) GetGets() int32 {
	if x != nil {
		return x.Gets
	}
	return 0
}

func (x *Service) GetPuts() int32 {
	if x != nil {
		return x.Puts
	}
	return 0
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
		mi := &file_data_services_services_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRequest) ProtoMessage() {}

func (x *ListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_data_services_services_proto_msgTypes[1]
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
	return file_data_services_services_proto_rawDescGZIP(), []int{1}
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

	Services   []*Service `protobuf:"bytes,1,rep,name=services,proto3" json:"services,omitempty"`
	LastUpdate int64      `protobuf:"varint,2,opt,name=last_update,json=lastUpdate,proto3" json:"last_update,omitempty"`
}

func (x *ListResponse) Reset() {
	*x = ListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_services_services_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListResponse) ProtoMessage() {}

func (x *ListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_data_services_services_proto_msgTypes[2]
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
	return file_data_services_services_proto_rawDescGZIP(), []int{2}
}

func (x *ListResponse) GetServices() []*Service {
	if x != nil {
		return x.Services
	}
	return nil
}

func (x *ListResponse) GetLastUpdate() int64 {
	if x != nil {
		return x.LastUpdate
	}
	return 0
}

type Service_Checker struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type           checker.Type `protobuf:"varint,1,opt,name=type,proto3,enum=checker.Type" json:"type,omitempty"`
	Path           string       `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	TimeoutSeconds int32        `protobuf:"varint,3,opt,name=timeout_seconds,json=timeoutSeconds,proto3" json:"timeout_seconds,omitempty"`
}

func (x *Service_Checker) Reset() {
	*x = Service_Checker{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_services_services_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Service_Checker) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Service_Checker) ProtoMessage() {}

func (x *Service_Checker) ProtoReflect() protoreflect.Message {
	mi := &file_data_services_services_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Service_Checker.ProtoReflect.Descriptor instead.
func (*Service_Checker) Descriptor() ([]byte, []int) {
	return file_data_services_services_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Service_Checker) GetType() checker.Type {
	if x != nil {
		return x.Type
	}
	return checker.Type(0)
}

func (x *Service_Checker) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *Service_Checker) GetTimeoutSeconds() int32 {
	if x != nil {
		return x.TimeoutSeconds
	}
	return 0
}

var File_data_services_services_proto protoreflect.FileDescriptor

var file_data_services_services_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d,
	0x64, 0x61, 0x74, 0x61, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x1a, 0x15, 0x63,
	0x68, 0x65, 0x63, 0x6b, 0x65, 0x72, 0x2f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9f, 0x02, 0x0a, 0x07, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x38, 0x0a, 0x07, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x65, 0x72, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x65, 0x72, 0x52, 0x07, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x65, 0x72, 0x12, 0x23,
	0x0a, 0x0d, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0c, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x53, 0x63,
	0x6f, 0x72, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x67, 0x65, 0x74, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x04, 0x67, 0x65, 0x74, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x75, 0x74, 0x73, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x75, 0x74, 0x73, 0x1a, 0x69, 0x0a, 0x07, 0x43,
	0x68, 0x65, 0x63, 0x6b, 0x65, 0x72, 0x12, 0x21, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74,
	0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x27, 0x0a,
	0x0f, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x5f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x53,
	0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x22, 0x2e, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x6c, 0x61, 0x73, 0x74,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x22, 0x63, 0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x52, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x6c, 0x61,
	0x73, 0x74, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0a, 0x6c, 0x61, 0x73, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0xae, 0x01, 0x0a, 0x11,
	0x63, 0x6f, 0x6d, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x42, 0x0d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x35, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63,
	0x34, 0x74, 0x2d, 0x62, 0x75, 0x74, 0x2d, 0x73, 0x34, 0x64, 0x2f, 0x66, 0x61, 0x73, 0x74, 0x61,
	0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x61, 0x74, 0x61,
	0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0xa2, 0x02, 0x03, 0x44, 0x53, 0x58, 0xaa,
	0x02, 0x0d, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0xca,
	0x02, 0x0d, 0x44, 0x61, 0x74, 0x61, 0x5c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0xe2,
	0x02, 0x19, 0x44, 0x61, 0x74, 0x61, 0x5c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x5c,
	0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0e, 0x44, 0x61,
	0x74, 0x61, 0x3a, 0x3a, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_data_services_services_proto_rawDescOnce sync.Once
	file_data_services_services_proto_rawDescData = file_data_services_services_proto_rawDesc
)

func file_data_services_services_proto_rawDescGZIP() []byte {
	file_data_services_services_proto_rawDescOnce.Do(func() {
		file_data_services_services_proto_rawDescData = protoimpl.X.CompressGZIP(file_data_services_services_proto_rawDescData)
	})
	return file_data_services_services_proto_rawDescData
}

var file_data_services_services_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_data_services_services_proto_goTypes = []interface{}{
	(*Service)(nil),         // 0: data.services.Service
	(*ListRequest)(nil),     // 1: data.services.ListRequest
	(*ListResponse)(nil),    // 2: data.services.ListResponse
	(*Service_Checker)(nil), // 3: data.services.Service.Checker
	(checker.Type)(0),       // 4: checker.Type
}
var file_data_services_services_proto_depIdxs = []int32{
	3, // 0: data.services.Service.checker:type_name -> data.services.Service.Checker
	0, // 1: data.services.ListResponse.services:type_name -> data.services.Service
	4, // 2: data.services.Service.Checker.type:type_name -> checker.Type
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_data_services_services_proto_init() }
func file_data_services_services_proto_init() {
	if File_data_services_services_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_data_services_services_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Service); i {
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
		file_data_services_services_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_data_services_services_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_data_services_services_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Service_Checker); i {
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
			RawDescriptor: file_data_services_services_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_data_services_services_proto_goTypes,
		DependencyIndexes: file_data_services_services_proto_depIdxs,
		MessageInfos:      file_data_services_services_proto_msgTypes,
	}.Build()
	File_data_services_services_proto = out.File
	file_data_services_services_proto_rawDesc = nil
	file_data_services_services_proto_goTypes = nil
	file_data_services_services_proto_depIdxs = nil
}