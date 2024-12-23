// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        (unknown)
// source: receiver/receiver.proto

package receiver

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

type FlagResponse_Verdict int32

const (
	FlagResponse_VERDICT_UNSPECIFIED FlagResponse_Verdict = 0
	FlagResponse_VERDICT_ACCEPTED    FlagResponse_Verdict = 1
	FlagResponse_VERDICT_OWN         FlagResponse_Verdict = 2
	FlagResponse_VERDICT_OLD         FlagResponse_Verdict = 3
	FlagResponse_VERDICT_INVALID     FlagResponse_Verdict = 4
	FlagResponse_VERDICT_DUPLICATE   FlagResponse_Verdict = 5
)

// Enum value maps for FlagResponse_Verdict.
var (
	FlagResponse_Verdict_name = map[int32]string{
		0: "VERDICT_UNSPECIFIED",
		1: "VERDICT_ACCEPTED",
		2: "VERDICT_OWN",
		3: "VERDICT_OLD",
		4: "VERDICT_INVALID",
		5: "VERDICT_DUPLICATE",
	}
	FlagResponse_Verdict_value = map[string]int32{
		"VERDICT_UNSPECIFIED": 0,
		"VERDICT_ACCEPTED":    1,
		"VERDICT_OWN":         2,
		"VERDICT_OLD":         3,
		"VERDICT_INVALID":     4,
		"VERDICT_DUPLICATE":   5,
	}
)

func (x FlagResponse_Verdict) Enum() *FlagResponse_Verdict {
	p := new(FlagResponse_Verdict)
	*p = x
	return p
}

func (x FlagResponse_Verdict) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FlagResponse_Verdict) Descriptor() protoreflect.EnumDescriptor {
	return file_receiver_receiver_proto_enumTypes[0].Descriptor()
}

func (FlagResponse_Verdict) Type() protoreflect.EnumType {
	return &file_receiver_receiver_proto_enumTypes[0]
}

func (x FlagResponse_Verdict) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FlagResponse_Verdict.Descriptor instead.
func (FlagResponse_Verdict) EnumDescriptor() ([]byte, []int) {
	return file_receiver_receiver_proto_rawDescGZIP(), []int{1, 0}
}

type SubmitFlagsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Flags         []string               `protobuf:"bytes,1,rep,name=flags,proto3" json:"flags,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SubmitFlagsRequest) Reset() {
	*x = SubmitFlagsRequest{}
	mi := &file_receiver_receiver_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmitFlagsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitFlagsRequest) ProtoMessage() {}

func (x *SubmitFlagsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_receiver_receiver_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitFlagsRequest.ProtoReflect.Descriptor instead.
func (*SubmitFlagsRequest) Descriptor() ([]byte, []int) {
	return file_receiver_receiver_proto_rawDescGZIP(), []int{0}
}

func (x *SubmitFlagsRequest) GetFlags() []string {
	if x != nil {
		return x.Flags
	}
	return nil
}

type FlagResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Flag          string                 `protobuf:"bytes,1,opt,name=flag,proto3" json:"flag,omitempty"`
	Verdict       FlagResponse_Verdict   `protobuf:"varint,2,opt,name=verdict,proto3,enum=receiver.FlagResponse_Verdict" json:"verdict,omitempty"`
	Message       string                 `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	VictimId      int64                  `protobuf:"varint,4,opt,name=victim_id,json=victimId,proto3" json:"victim_id,omitempty"`
	ServiceId     int64                  `protobuf:"varint,5,opt,name=service_id,json=serviceId,proto3" json:"service_id,omitempty"`
	AttackerDelta float64                `protobuf:"fixed64,6,opt,name=attacker_delta,json=attackerDelta,proto3" json:"attacker_delta,omitempty"`
	VictimDelta   float64                `protobuf:"fixed64,7,opt,name=victim_delta,json=victimDelta,proto3" json:"victim_delta,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FlagResponse) Reset() {
	*x = FlagResponse{}
	mi := &file_receiver_receiver_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FlagResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlagResponse) ProtoMessage() {}

func (x *FlagResponse) ProtoReflect() protoreflect.Message {
	mi := &file_receiver_receiver_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlagResponse.ProtoReflect.Descriptor instead.
func (*FlagResponse) Descriptor() ([]byte, []int) {
	return file_receiver_receiver_proto_rawDescGZIP(), []int{1}
}

func (x *FlagResponse) GetFlag() string {
	if x != nil {
		return x.Flag
	}
	return ""
}

func (x *FlagResponse) GetVerdict() FlagResponse_Verdict {
	if x != nil {
		return x.Verdict
	}
	return FlagResponse_VERDICT_UNSPECIFIED
}

func (x *FlagResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *FlagResponse) GetVictimId() int64 {
	if x != nil {
		return x.VictimId
	}
	return 0
}

func (x *FlagResponse) GetServiceId() int64 {
	if x != nil {
		return x.ServiceId
	}
	return 0
}

func (x *FlagResponse) GetAttackerDelta() float64 {
	if x != nil {
		return x.AttackerDelta
	}
	return 0
}

func (x *FlagResponse) GetVictimDelta() float64 {
	if x != nil {
		return x.VictimDelta
	}
	return 0
}

type State struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Services      []*State_Service       `protobuf:"bytes,1,rep,name=services,proto3" json:"services,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *State) Reset() {
	*x = State{}
	mi := &file_receiver_receiver_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *State) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*State) ProtoMessage() {}

func (x *State) ProtoReflect() protoreflect.Message {
	mi := &file_receiver_receiver_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use State.ProtoReflect.Descriptor instead.
func (*State) Descriptor() ([]byte, []int) {
	return file_receiver_receiver_proto_rawDescGZIP(), []int{2}
}

func (x *State) GetServices() []*State_Service {
	if x != nil {
		return x.Services
	}
	return nil
}

type SubmitFlagsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Responses     []*FlagResponse        `protobuf:"bytes,1,rep,name=responses,proto3" json:"responses,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SubmitFlagsResponse) Reset() {
	*x = SubmitFlagsResponse{}
	mi := &file_receiver_receiver_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmitFlagsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitFlagsResponse) ProtoMessage() {}

func (x *SubmitFlagsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_receiver_receiver_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitFlagsResponse.ProtoReflect.Descriptor instead.
func (*SubmitFlagsResponse) Descriptor() ([]byte, []int) {
	return file_receiver_receiver_proto_rawDescGZIP(), []int{3}
}

func (x *SubmitFlagsResponse) GetResponses() []*FlagResponse {
	if x != nil {
		return x.Responses
	}
	return nil
}

type GetStateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetStateRequest) Reset() {
	*x = GetStateRequest{}
	mi := &file_receiver_receiver_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetStateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStateRequest) ProtoMessage() {}

func (x *GetStateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_receiver_receiver_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStateRequest.ProtoReflect.Descriptor instead.
func (*GetStateRequest) Descriptor() ([]byte, []int) {
	return file_receiver_receiver_proto_rawDescGZIP(), []int{4}
}

type GetStateResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	State         *State                 `protobuf:"bytes,1,opt,name=state,proto3" json:"state,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetStateResponse) Reset() {
	*x = GetStateResponse{}
	mi := &file_receiver_receiver_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetStateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStateResponse) ProtoMessage() {}

func (x *GetStateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_receiver_receiver_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStateResponse.ProtoReflect.Descriptor instead.
func (*GetStateResponse) Descriptor() ([]byte, []int) {
	return file_receiver_receiver_proto_rawDescGZIP(), []int{5}
}

func (x *GetStateResponse) GetState() *State {
	if x != nil {
		return x.State
	}
	return nil
}

type State_Team struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Points        float64                `protobuf:"fixed64,2,opt,name=points,proto3" json:"points,omitempty"`
	StolenFlags   int64                  `protobuf:"varint,3,opt,name=stolen_flags,json=stolenFlags,proto3" json:"stolen_flags,omitempty"`
	LostFlags     int64                  `protobuf:"varint,4,opt,name=lost_flags,json=lostFlags,proto3" json:"lost_flags,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *State_Team) Reset() {
	*x = State_Team{}
	mi := &file_receiver_receiver_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *State_Team) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*State_Team) ProtoMessage() {}

func (x *State_Team) ProtoReflect() protoreflect.Message {
	mi := &file_receiver_receiver_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use State_Team.ProtoReflect.Descriptor instead.
func (*State_Team) Descriptor() ([]byte, []int) {
	return file_receiver_receiver_proto_rawDescGZIP(), []int{2, 0}
}

func (x *State_Team) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *State_Team) GetPoints() float64 {
	if x != nil {
		return x.Points
	}
	return 0
}

func (x *State_Team) GetStolenFlags() int64 {
	if x != nil {
		return x.StolenFlags
	}
	return 0
}

func (x *State_Team) GetLostFlags() int64 {
	if x != nil {
		return x.LostFlags
	}
	return 0
}

type State_Service struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Teams         []*State_Team          `protobuf:"bytes,2,rep,name=teams,proto3" json:"teams,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *State_Service) Reset() {
	*x = State_Service{}
	mi := &file_receiver_receiver_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *State_Service) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*State_Service) ProtoMessage() {}

func (x *State_Service) ProtoReflect() protoreflect.Message {
	mi := &file_receiver_receiver_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use State_Service.ProtoReflect.Descriptor instead.
func (*State_Service) Descriptor() ([]byte, []int) {
	return file_receiver_receiver_proto_rawDescGZIP(), []int{2, 1}
}

func (x *State_Service) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *State_Service) GetTeams() []*State_Team {
	if x != nil {
		return x.Teams
	}
	return nil
}

var File_receiver_receiver_proto protoreflect.FileDescriptor

var file_receiver_receiver_proto_rawDesc = []byte{
	0x0a, 0x17, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x2f, 0x72, 0x65, 0x63, 0x65, 0x69,
	0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x72, 0x65, 0x63, 0x65, 0x69,
	0x76, 0x65, 0x72, 0x22, 0x2a, 0x0a, 0x12, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x46, 0x6c, 0x61,
	0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x6c, 0x61,
	0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x22,
	0x85, 0x03, 0x0a, 0x0c, 0x46, 0x6c, 0x61, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x66, 0x6c, 0x61, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x66, 0x6c, 0x61, 0x67, 0x12, 0x38, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x64, 0x69, 0x63, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72,
	0x2e, 0x46, 0x6c, 0x61, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x56, 0x65,
	0x72, 0x64, 0x69, 0x63, 0x74, 0x52, 0x07, 0x76, 0x65, 0x72, 0x64, 0x69, 0x63, 0x74, 0x12, 0x18,
	0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x76, 0x69, 0x63, 0x74,
	0x69, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x76, 0x69, 0x63,
	0x74, 0x69, 0x6d, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x61, 0x74, 0x74, 0x61, 0x63, 0x6b, 0x65, 0x72,
	0x5f, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0d, 0x61, 0x74,
	0x74, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x44, 0x65, 0x6c, 0x74, 0x61, 0x12, 0x21, 0x0a, 0x0c, 0x76,
	0x69, 0x63, 0x74, 0x69, 0x6d, 0x5f, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x0b, 0x76, 0x69, 0x63, 0x74, 0x69, 0x6d, 0x44, 0x65, 0x6c, 0x74, 0x61, 0x22, 0x86,
	0x01, 0x0a, 0x07, 0x56, 0x65, 0x72, 0x64, 0x69, 0x63, 0x74, 0x12, 0x17, 0x0a, 0x13, 0x56, 0x45,
	0x52, 0x44, 0x49, 0x43, 0x54, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45,
	0x44, 0x10, 0x00, 0x12, 0x14, 0x0a, 0x10, 0x56, 0x45, 0x52, 0x44, 0x49, 0x43, 0x54, 0x5f, 0x41,
	0x43, 0x43, 0x45, 0x50, 0x54, 0x45, 0x44, 0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x56, 0x45, 0x52,
	0x44, 0x49, 0x43, 0x54, 0x5f, 0x4f, 0x57, 0x4e, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x56, 0x45,
	0x52, 0x44, 0x49, 0x43, 0x54, 0x5f, 0x4f, 0x4c, 0x44, 0x10, 0x03, 0x12, 0x13, 0x0a, 0x0f, 0x56,
	0x45, 0x52, 0x44, 0x49, 0x43, 0x54, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x04,
	0x12, 0x15, 0x0a, 0x11, 0x56, 0x45, 0x52, 0x44, 0x49, 0x43, 0x54, 0x5f, 0x44, 0x55, 0x50, 0x4c,
	0x49, 0x43, 0x41, 0x54, 0x45, 0x10, 0x05, 0x22, 0xf5, 0x01, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x12, 0x33, 0x0a, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x2e, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x08, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x1a, 0x70, 0x0a, 0x04, 0x54, 0x65, 0x61, 0x6d, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16,
	0x0a, 0x06, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x74, 0x6f, 0x6c, 0x65, 0x6e,
	0x5f, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x73, 0x74,
	0x6f, 0x6c, 0x65, 0x6e, 0x46, 0x6c, 0x61, 0x67, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x6f, 0x73,
	0x74, 0x5f, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x6c,
	0x6f, 0x73, 0x74, 0x46, 0x6c, 0x61, 0x67, 0x73, 0x1a, 0x45, 0x0a, 0x07, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x2a, 0x0a, 0x05, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x2e, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x05, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x22,
	0x4b, 0x0a, 0x13, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x46, 0x6c, 0x61, 0x67, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x34, 0x0a, 0x09, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x72, 0x65, 0x63, 0x65,
	0x69, 0x76, 0x65, 0x72, 0x2e, 0x46, 0x6c, 0x61, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x52, 0x09, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x73, 0x22, 0x11, 0x0a, 0x0f,
	0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22,
	0x39, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x42, 0x8f, 0x01, 0x0a, 0x0c, 0x63,
	0x6f, 0x6d, 0x2e, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x42, 0x0d, 0x52, 0x65, 0x63,
	0x65, 0x69, 0x76, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x30, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x34, 0x74, 0x2d, 0x62, 0x75, 0x74,
	0x2d, 0x73, 0x34, 0x64, 0x2f, 0x66, 0x61, 0x73, 0x74, 0x61, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0xa2, 0x02,
	0x03, 0x52, 0x58, 0x58, 0xaa, 0x02, 0x08, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0xca,
	0x02, 0x08, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0xe2, 0x02, 0x14, 0x52, 0x65, 0x63,
	0x65, 0x69, 0x76, 0x65, 0x72, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x08, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_receiver_receiver_proto_rawDescOnce sync.Once
	file_receiver_receiver_proto_rawDescData = file_receiver_receiver_proto_rawDesc
)

func file_receiver_receiver_proto_rawDescGZIP() []byte {
	file_receiver_receiver_proto_rawDescOnce.Do(func() {
		file_receiver_receiver_proto_rawDescData = protoimpl.X.CompressGZIP(file_receiver_receiver_proto_rawDescData)
	})
	return file_receiver_receiver_proto_rawDescData
}

var file_receiver_receiver_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_receiver_receiver_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_receiver_receiver_proto_goTypes = []any{
	(FlagResponse_Verdict)(0),   // 0: receiver.FlagResponse.Verdict
	(*SubmitFlagsRequest)(nil),  // 1: receiver.SubmitFlagsRequest
	(*FlagResponse)(nil),        // 2: receiver.FlagResponse
	(*State)(nil),               // 3: receiver.State
	(*SubmitFlagsResponse)(nil), // 4: receiver.SubmitFlagsResponse
	(*GetStateRequest)(nil),     // 5: receiver.GetStateRequest
	(*GetStateResponse)(nil),    // 6: receiver.GetStateResponse
	(*State_Team)(nil),          // 7: receiver.State.Team
	(*State_Service)(nil),       // 8: receiver.State.Service
}
var file_receiver_receiver_proto_depIdxs = []int32{
	0, // 0: receiver.FlagResponse.verdict:type_name -> receiver.FlagResponse.Verdict
	8, // 1: receiver.State.services:type_name -> receiver.State.Service
	2, // 2: receiver.SubmitFlagsResponse.responses:type_name -> receiver.FlagResponse
	3, // 3: receiver.GetStateResponse.state:type_name -> receiver.State
	7, // 4: receiver.State.Service.teams:type_name -> receiver.State.Team
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_receiver_receiver_proto_init() }
func file_receiver_receiver_proto_init() {
	if File_receiver_receiver_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_receiver_receiver_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_receiver_receiver_proto_goTypes,
		DependencyIndexes: file_receiver_receiver_proto_depIdxs,
		EnumInfos:         file_receiver_receiver_proto_enumTypes,
		MessageInfos:      file_receiver_receiver_proto_msgTypes,
	}.Build()
	File_receiver_receiver_proto = out.File
	file_receiver_receiver_proto_rawDesc = nil
	file_receiver_receiver_proto_goTypes = nil
	file_receiver_receiver_proto_depIdxs = nil
}
