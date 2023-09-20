// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: receiver/receiver_service.proto

package receiver

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	ReceiverService_SubmitFlags_FullMethodName = "/receiver.ReceiverService/SubmitFlags"
)

// ReceiverServiceClient is the client API for ReceiverService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReceiverServiceClient interface {
	SubmitFlags(ctx context.Context, in *SubmitFlagsRequest, opts ...grpc.CallOption) (*SubmitFlagsResponse, error)
}

type receiverServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewReceiverServiceClient(cc grpc.ClientConnInterface) ReceiverServiceClient {
	return &receiverServiceClient{cc}
}

func (c *receiverServiceClient) SubmitFlags(ctx context.Context, in *SubmitFlagsRequest, opts ...grpc.CallOption) (*SubmitFlagsResponse, error) {
	out := new(SubmitFlagsResponse)
	err := c.cc.Invoke(ctx, ReceiverService_SubmitFlags_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReceiverServiceServer is the server API for ReceiverService service.
// All implementations must embed UnimplementedReceiverServiceServer
// for forward compatibility
type ReceiverServiceServer interface {
	SubmitFlags(context.Context, *SubmitFlagsRequest) (*SubmitFlagsResponse, error)
	mustEmbedUnimplementedReceiverServiceServer()
}

// UnimplementedReceiverServiceServer must be embedded to have forward compatible implementations.
type UnimplementedReceiverServiceServer struct {
}

func (UnimplementedReceiverServiceServer) SubmitFlags(context.Context, *SubmitFlagsRequest) (*SubmitFlagsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitFlags not implemented")
}
func (UnimplementedReceiverServiceServer) mustEmbedUnimplementedReceiverServiceServer() {}

// UnsafeReceiverServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReceiverServiceServer will
// result in compilation errors.
type UnsafeReceiverServiceServer interface {
	mustEmbedUnimplementedReceiverServiceServer()
}

func RegisterReceiverServiceServer(s grpc.ServiceRegistrar, srv ReceiverServiceServer) {
	s.RegisterService(&ReceiverService_ServiceDesc, srv)
}

func _ReceiverService_SubmitFlags_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitFlagsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReceiverServiceServer).SubmitFlags(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReceiverService_SubmitFlags_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReceiverServiceServer).SubmitFlags(ctx, req.(*SubmitFlagsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ReceiverService_ServiceDesc is the grpc.ServiceDesc for ReceiverService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReceiverService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "receiver.ReceiverService",
	HandlerType: (*ReceiverServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubmitFlags",
			Handler:    _ReceiverService_SubmitFlags_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "receiver/receiver_service.proto",
}
