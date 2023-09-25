// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: data/services/services_service.proto

package services

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
	ServicesService_List_FullMethodName = "/data.services.ServicesService/List"
)

// ServicesServiceClient is the client API for ServicesService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServicesServiceClient interface {
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
}

type servicesServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewServicesServiceClient(cc grpc.ClientConnInterface) ServicesServiceClient {
	return &servicesServiceClient{cc}
}

func (c *servicesServiceClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, ServicesService_List_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServicesServiceServer is the server API for ServicesService service.
// All implementations must embed UnimplementedServicesServiceServer
// for forward compatibility
type ServicesServiceServer interface {
	List(context.Context, *ListRequest) (*ListResponse, error)
	mustEmbedUnimplementedServicesServiceServer()
}

// UnimplementedServicesServiceServer must be embedded to have forward compatible implementations.
type UnimplementedServicesServiceServer struct {
}

func (UnimplementedServicesServiceServer) List(context.Context, *ListRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedServicesServiceServer) mustEmbedUnimplementedServicesServiceServer() {}

// UnsafeServicesServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServicesServiceServer will
// result in compilation errors.
type UnsafeServicesServiceServer interface {
	mustEmbedUnimplementedServicesServiceServer()
}

func RegisterServicesServiceServer(s grpc.ServiceRegistrar, srv ServicesServiceServer) {
	s.RegisterService(&ServicesService_ServiceDesc, srv)
}

func _ServicesService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ServicesService_List_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServiceServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ServicesService_ServiceDesc is the grpc.ServiceDesc for ServicesService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ServicesService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "data.services.ServicesService",
	HandlerType: (*ServicesServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "List",
			Handler:    _ServicesService_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "data/services/services_service.proto",
}
