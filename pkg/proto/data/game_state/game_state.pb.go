// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: data/game_state/game_state.proto

package game_state

import (
	version "github.com/c4t-but-s4d/fastad/pkg/proto/data/version"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GameMode int32

const (
	GameMode_GAME_MODE_UNSPECIFIED GameMode = 0
	GameMode_GAME_MODE_CLASSIC     GameMode = 1
)

// Enum value maps for GameMode.
var (
	GameMode_name = map[int32]string{
		0: "GAME_MODE_UNSPECIFIED",
		1: "GAME_MODE_CLASSIC",
	}
	GameMode_value = map[string]int32{
		"GAME_MODE_UNSPECIFIED": 0,
		"GAME_MODE_CLASSIC":     1,
	}
)

func (x GameMode) Enum() *GameMode {
	p := new(GameMode)
	*p = x
	return p
}

func (x GameMode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GameMode) Descriptor() protoreflect.EnumDescriptor {
	return file_data_game_state_game_state_proto_enumTypes[0].Descriptor()
}

func (GameMode) Type() protoreflect.EnumType {
	return &file_data_game_state_game_state_proto_enumTypes[0]
}

func (x GameMode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GameMode.Descriptor instead.
func (GameMode) EnumDescriptor() ([]byte, []int) {
	return file_data_game_state_game_state_proto_rawDescGZIP(), []int{0}
}

type GameState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartTime          *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime            *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	TotalRounds        uint32                 `protobuf:"varint,3,opt,name=total_rounds,json=totalRounds,proto3" json:"total_rounds,omitempty"`
	Paused             bool                   `protobuf:"varint,4,opt,name=paused,proto3" json:"paused,omitempty"`
	FlagLifetimeRounds uint32                 `protobuf:"varint,5,opt,name=flag_lifetime_rounds,json=flagLifetimeRounds,proto3" json:"flag_lifetime_rounds,omitempty"`
	RoundDuration      *durationpb.Duration   `protobuf:"bytes,6,opt,name=round_duration,json=roundDuration,proto3" json:"round_duration,omitempty"`
	Mode               GameMode               `protobuf:"varint,7,opt,name=mode,proto3,enum=data.game_state.GameMode" json:"mode,omitempty"`
	RunningRound       uint32                 `protobuf:"varint,8,opt,name=running_round,json=runningRound,proto3" json:"running_round,omitempty"`
	RunningRoundStart  *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=running_round_start,json=runningRoundStart,proto3" json:"running_round_start,omitempty"`
}

func (x *GameState) Reset() {
	*x = GameState{}
	mi := &file_data_game_state_game_state_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GameState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameState) ProtoMessage() {}

func (x *GameState) ProtoReflect() protoreflect.Message {
	mi := &file_data_game_state_game_state_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GameState.ProtoReflect.Descriptor instead.
func (*GameState) Descriptor() ([]byte, []int) {
	return file_data_game_state_game_state_proto_rawDescGZIP(), []int{0}
}

func (x *GameState) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *GameState) GetEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

func (x *GameState) GetTotalRounds() uint32 {
	if x != nil {
		return x.TotalRounds
	}
	return 0
}

func (x *GameState) GetPaused() bool {
	if x != nil {
		return x.Paused
	}
	return false
}

func (x *GameState) GetFlagLifetimeRounds() uint32 {
	if x != nil {
		return x.FlagLifetimeRounds
	}
	return 0
}

func (x *GameState) GetRoundDuration() *durationpb.Duration {
	if x != nil {
		return x.RoundDuration
	}
	return nil
}

func (x *GameState) GetMode() GameMode {
	if x != nil {
		return x.Mode
	}
	return GameMode_GAME_MODE_UNSPECIFIED
}

func (x *GameState) GetRunningRound() uint32 {
	if x != nil {
		return x.RunningRound
	}
	return 0
}

func (x *GameState) GetRunningRoundStart() *timestamppb.Timestamp {
	if x != nil {
		return x.RunningRoundStart
	}
	return nil
}

type GetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version *version.Version `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *GetRequest) Reset() {
	*x = GetRequest{}
	mi := &file_data_game_state_game_state_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRequest) ProtoMessage() {}

func (x *GetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_data_game_state_game_state_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRequest.ProtoReflect.Descriptor instead.
func (*GetRequest) Descriptor() ([]byte, []int) {
	return file_data_game_state_game_state_proto_rawDescGZIP(), []int{1}
}

func (x *GetRequest) GetVersion() *version.Version {
	if x != nil {
		return x.Version
	}
	return nil
}

type GetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GameState *GameState       `protobuf:"bytes,1,opt,name=game_state,json=gameState,proto3" json:"game_state,omitempty"`
	Version   *version.Version `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *GetResponse) Reset() {
	*x = GetResponse{}
	mi := &file_data_game_state_game_state_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetResponse) ProtoMessage() {}

func (x *GetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_data_game_state_game_state_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetResponse.ProtoReflect.Descriptor instead.
func (*GetResponse) Descriptor() ([]byte, []int) {
	return file_data_game_state_game_state_proto_rawDescGZIP(), []int{2}
}

func (x *GetResponse) GetGameState() *GameState {
	if x != nil {
		return x.GameState
	}
	return nil
}

func (x *GetResponse) GetVersion() *version.Version {
	if x != nil {
		return x.Version
	}
	return nil
}

type UpdateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartTime          *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime            *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	TotalRounds        uint32                 `protobuf:"varint,3,opt,name=total_rounds,json=totalRounds,proto3" json:"total_rounds,omitempty"`
	Paused             bool                   `protobuf:"varint,4,opt,name=paused,proto3" json:"paused,omitempty"`
	FlagLifetimeRounds uint32                 `protobuf:"varint,5,opt,name=flag_lifetime_rounds,json=flagLifetimeRounds,proto3" json:"flag_lifetime_rounds,omitempty"`
	RoundDuration      *durationpb.Duration   `protobuf:"bytes,6,opt,name=round_duration,json=roundDuration,proto3" json:"round_duration,omitempty"`
	Mode               GameMode               `protobuf:"varint,7,opt,name=mode,proto3,enum=data.game_state.GameMode" json:"mode,omitempty"`
}

func (x *UpdateRequest) Reset() {
	*x = UpdateRequest{}
	mi := &file_data_game_state_game_state_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateRequest) ProtoMessage() {}

func (x *UpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_data_game_state_game_state_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateRequest.ProtoReflect.Descriptor instead.
func (*UpdateRequest) Descriptor() ([]byte, []int) {
	return file_data_game_state_game_state_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateRequest) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *UpdateRequest) GetEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

func (x *UpdateRequest) GetTotalRounds() uint32 {
	if x != nil {
		return x.TotalRounds
	}
	return 0
}

func (x *UpdateRequest) GetPaused() bool {
	if x != nil {
		return x.Paused
	}
	return false
}

func (x *UpdateRequest) GetFlagLifetimeRounds() uint32 {
	if x != nil {
		return x.FlagLifetimeRounds
	}
	return 0
}

func (x *UpdateRequest) GetRoundDuration() *durationpb.Duration {
	if x != nil {
		return x.RoundDuration
	}
	return nil
}

func (x *UpdateRequest) GetMode() GameMode {
	if x != nil {
		return x.Mode
	}
	return GameMode_GAME_MODE_UNSPECIFIED
}

type UpdateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GameState *GameState       `protobuf:"bytes,1,opt,name=game_state,json=gameState,proto3" json:"game_state,omitempty"`
	Version   *version.Version `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *UpdateResponse) Reset() {
	*x = UpdateResponse{}
	mi := &file_data_game_state_game_state_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateResponse) ProtoMessage() {}

func (x *UpdateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_data_game_state_game_state_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateResponse.ProtoReflect.Descriptor instead.
func (*UpdateResponse) Descriptor() ([]byte, []int) {
	return file_data_game_state_game_state_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateResponse) GetGameState() *GameState {
	if x != nil {
		return x.GameState
	}
	return nil
}

func (x *UpdateResponse) GetVersion() *version.Version {
	if x != nil {
		return x.Version
	}
	return nil
}

type UpdateRoundRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RunningRound      uint32                 `protobuf:"varint,1,opt,name=running_round,json=runningRound,proto3" json:"running_round,omitempty"`
	RunningRoundStart *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=running_round_start,json=runningRoundStart,proto3" json:"running_round_start,omitempty"`
}

func (x *UpdateRoundRequest) Reset() {
	*x = UpdateRoundRequest{}
	mi := &file_data_game_state_game_state_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateRoundRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateRoundRequest) ProtoMessage() {}

func (x *UpdateRoundRequest) ProtoReflect() protoreflect.Message {
	mi := &file_data_game_state_game_state_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateRoundRequest.ProtoReflect.Descriptor instead.
func (*UpdateRoundRequest) Descriptor() ([]byte, []int) {
	return file_data_game_state_game_state_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateRoundRequest) GetRunningRound() uint32 {
	if x != nil {
		return x.RunningRound
	}
	return 0
}

func (x *UpdateRoundRequest) GetRunningRoundStart() *timestamppb.Timestamp {
	if x != nil {
		return x.RunningRoundStart
	}
	return nil
}

type UpdateRoundResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GameState *GameState       `protobuf:"bytes,1,opt,name=game_state,json=gameState,proto3" json:"game_state,omitempty"`
	Version   *version.Version `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *UpdateRoundResponse) Reset() {
	*x = UpdateRoundResponse{}
	mi := &file_data_game_state_game_state_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateRoundResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateRoundResponse) ProtoMessage() {}

func (x *UpdateRoundResponse) ProtoReflect() protoreflect.Message {
	mi := &file_data_game_state_game_state_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateRoundResponse.ProtoReflect.Descriptor instead.
func (*UpdateRoundResponse) Descriptor() ([]byte, []int) {
	return file_data_game_state_game_state_proto_rawDescGZIP(), []int{6}
}

func (x *UpdateRoundResponse) GetGameState() *GameState {
	if x != nil {
		return x.GameState
	}
	return nil
}

func (x *UpdateRoundResponse) GetVersion() *version.Version {
	if x != nil {
		return x.Version
	}
	return nil
}

var File_data_game_state_game_state_proto protoreflect.FileDescriptor

var file_data_game_state_game_state_proto_rawDesc = []byte{
	0x0a, 0x20, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0f, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xcc, 0x03, 0x0a, 0x09, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x39,
	0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x65, 0x6e, 0x64,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x21, 0x0a, 0x0c, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x73,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x52, 0x6f, 0x75,
	0x6e, 0x64, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x61, 0x75, 0x73, 0x65, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x06, 0x70, 0x61, 0x75, 0x73, 0x65, 0x64, 0x12, 0x30, 0x0a, 0x14, 0x66,
	0x6c, 0x61, 0x67, 0x5f, 0x6c, 0x69, 0x66, 0x65, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x72, 0x6f, 0x75,
	0x6e, 0x64, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x12, 0x66, 0x6c, 0x61, 0x67, 0x4c,
	0x69, 0x66, 0x65, 0x74, 0x69, 0x6d, 0x65, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x73, 0x12, 0x40, 0x0a,
	0x0e, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x5f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x0d, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x2d, 0x0a, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e,
	0x64, 0x61, 0x74, 0x61, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e,
	0x47, 0x61, 0x6d, 0x65, 0x4d, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x12, 0x23,
	0x0a, 0x0d, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x52, 0x6f,
	0x75, 0x6e, 0x64, 0x12, 0x4a, 0x0a, 0x13, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x72,
	0x6f, 0x75, 0x6e, 0x64, 0x5f, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x11, 0x72, 0x75,
	0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x53, 0x74, 0x61, 0x72, 0x74, 0x22,
	0x3d, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2f, 0x0a,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15,
	0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x79,
	0x0a, 0x0b, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a,
	0x0a, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x09, 0x67,
	0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x2f, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x64, 0x61, 0x74, 0x61,
	0x2e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0xdf, 0x02, 0x0a, 0x0d, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x21, 0x0a,
	0x0c, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x73,
	0x12, 0x16, 0x0a, 0x06, 0x70, 0x61, 0x75, 0x73, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x06, 0x70, 0x61, 0x75, 0x73, 0x65, 0x64, 0x12, 0x30, 0x0a, 0x14, 0x66, 0x6c, 0x61, 0x67,
	0x5f, 0x6c, 0x69, 0x66, 0x65, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x73,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x12, 0x66, 0x6c, 0x61, 0x67, 0x4c, 0x69, 0x66, 0x65,
	0x74, 0x69, 0x6d, 0x65, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x73, 0x12, 0x40, 0x0a, 0x0e, 0x72, 0x6f,
	0x75, 0x6e, 0x64, 0x5f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x72,
	0x6f, 0x75, 0x6e, 0x64, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2d, 0x0a, 0x04,
	0x6d, 0x6f, 0x64, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x64, 0x61, 0x74,
	0x61, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x47, 0x61, 0x6d,
	0x65, 0x4d, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x22, 0x7c, 0x0a, 0x0e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a,
	0x0a, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x09, 0x67,
	0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x2f, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x64, 0x61, 0x74, 0x61,
	0x2e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x85, 0x01, 0x0a, 0x12, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x23, 0x0a, 0x0d, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x72, 0x6f, 0x75, 0x6e,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67,
	0x52, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x4a, 0x0a, 0x13, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67,
	0x5f, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x5f, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x11,
	0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x22, 0x81, 0x01, 0x0a, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x75, 0x6e,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x67, 0x61, 0x6d,
	0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x64, 0x61, 0x74, 0x61, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e,
	0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x09, 0x67, 0x61, 0x6d, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x12, 0x2f, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2a, 0x3c, 0x0a, 0x08, 0x47, 0x61, 0x6d, 0x65, 0x4d, 0x6f, 0x64,
	0x65, 0x12, 0x19, 0x0a, 0x15, 0x47, 0x41, 0x4d, 0x45, 0x5f, 0x4d, 0x4f, 0x44, 0x45, 0x5f, 0x55,
	0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11,
	0x47, 0x41, 0x4d, 0x45, 0x5f, 0x4d, 0x4f, 0x44, 0x45, 0x5f, 0x43, 0x4c, 0x41, 0x53, 0x53, 0x49,
	0x43, 0x10, 0x01, 0x42, 0xb7, 0x01, 0x0a, 0x13, 0x63, 0x6f, 0x6d, 0x2e, 0x64, 0x61, 0x74, 0x61,
	0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x42, 0x0e, 0x47, 0x61, 0x6d,
	0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x37, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x34, 0x74, 0x2d, 0x62, 0x75,
	0x74, 0x2d, 0x73, 0x34, 0x64, 0x2f, 0x66, 0x61, 0x73, 0x74, 0x61, 0x64, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x67, 0x61, 0x6d, 0x65,
	0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0xa2, 0x02, 0x03, 0x44, 0x47, 0x58, 0xaa, 0x02, 0x0e, 0x44,
	0x61, 0x74, 0x61, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0xca, 0x02, 0x0e,
	0x44, 0x61, 0x74, 0x61, 0x5c, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0xe2, 0x02,
	0x1a, 0x44, 0x61, 0x74, 0x61, 0x5c, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x5c,
	0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0f, 0x44, 0x61,
	0x74, 0x61, 0x3a, 0x3a, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_data_game_state_game_state_proto_rawDescOnce sync.Once
	file_data_game_state_game_state_proto_rawDescData = file_data_game_state_game_state_proto_rawDesc
)

func file_data_game_state_game_state_proto_rawDescGZIP() []byte {
	file_data_game_state_game_state_proto_rawDescOnce.Do(func() {
		file_data_game_state_game_state_proto_rawDescData = protoimpl.X.CompressGZIP(file_data_game_state_game_state_proto_rawDescData)
	})
	return file_data_game_state_game_state_proto_rawDescData
}

var file_data_game_state_game_state_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_data_game_state_game_state_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_data_game_state_game_state_proto_goTypes = []any{
	(GameMode)(0),                 // 0: data.game_state.GameMode
	(*GameState)(nil),             // 1: data.game_state.GameState
	(*GetRequest)(nil),            // 2: data.game_state.GetRequest
	(*GetResponse)(nil),           // 3: data.game_state.GetResponse
	(*UpdateRequest)(nil),         // 4: data.game_state.UpdateRequest
	(*UpdateResponse)(nil),        // 5: data.game_state.UpdateResponse
	(*UpdateRoundRequest)(nil),    // 6: data.game_state.UpdateRoundRequest
	(*UpdateRoundResponse)(nil),   // 7: data.game_state.UpdateRoundResponse
	(*timestamppb.Timestamp)(nil), // 8: google.protobuf.Timestamp
	(*durationpb.Duration)(nil),   // 9: google.protobuf.Duration
	(*version.Version)(nil),       // 10: data.version.Version
}
var file_data_game_state_game_state_proto_depIdxs = []int32{
	8,  // 0: data.game_state.GameState.start_time:type_name -> google.protobuf.Timestamp
	8,  // 1: data.game_state.GameState.end_time:type_name -> google.protobuf.Timestamp
	9,  // 2: data.game_state.GameState.round_duration:type_name -> google.protobuf.Duration
	0,  // 3: data.game_state.GameState.mode:type_name -> data.game_state.GameMode
	8,  // 4: data.game_state.GameState.running_round_start:type_name -> google.protobuf.Timestamp
	10, // 5: data.game_state.GetRequest.version:type_name -> data.version.Version
	1,  // 6: data.game_state.GetResponse.game_state:type_name -> data.game_state.GameState
	10, // 7: data.game_state.GetResponse.version:type_name -> data.version.Version
	8,  // 8: data.game_state.UpdateRequest.start_time:type_name -> google.protobuf.Timestamp
	8,  // 9: data.game_state.UpdateRequest.end_time:type_name -> google.protobuf.Timestamp
	9,  // 10: data.game_state.UpdateRequest.round_duration:type_name -> google.protobuf.Duration
	0,  // 11: data.game_state.UpdateRequest.mode:type_name -> data.game_state.GameMode
	1,  // 12: data.game_state.UpdateResponse.game_state:type_name -> data.game_state.GameState
	10, // 13: data.game_state.UpdateResponse.version:type_name -> data.version.Version
	8,  // 14: data.game_state.UpdateRoundRequest.running_round_start:type_name -> google.protobuf.Timestamp
	1,  // 15: data.game_state.UpdateRoundResponse.game_state:type_name -> data.game_state.GameState
	10, // 16: data.game_state.UpdateRoundResponse.version:type_name -> data.version.Version
	17, // [17:17] is the sub-list for method output_type
	17, // [17:17] is the sub-list for method input_type
	17, // [17:17] is the sub-list for extension type_name
	17, // [17:17] is the sub-list for extension extendee
	0,  // [0:17] is the sub-list for field type_name
}

func init() { file_data_game_state_game_state_proto_init() }
func file_data_game_state_game_state_proto_init() {
	if File_data_game_state_game_state_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_data_game_state_game_state_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_data_game_state_game_state_proto_goTypes,
		DependencyIndexes: file_data_game_state_game_state_proto_depIdxs,
		EnumInfos:         file_data_game_state_game_state_proto_enumTypes,
		MessageInfos:      file_data_game_state_game_state_proto_msgTypes,
	}.Build()
	File_data_game_state_game_state_proto = out.File
	file_data_game_state_game_state_proto_rawDesc = nil
	file_data_game_state_game_state_proto_goTypes = nil
	file_data_game_state_game_state_proto_depIdxs = nil
}
