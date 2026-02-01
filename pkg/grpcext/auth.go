package grpcext

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const authKey = "authorization"

type ServerTokenInterceptor struct {
	token string
}

func NewServerTokenInterceptor(token string) *ServerTokenInterceptor {
	return &ServerTokenInterceptor{token: token}
}

func (ti *ServerTokenInterceptor) authorize(ctx context.Context, _ string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}
	values := md[authKey]
	if len(values) != 1 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}
	if ti.token != values[0] {
		return status.Errorf(codes.Unauthenticated, "invalid token provided")
	}
	return nil
}

func (ti *ServerTokenInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if err := ti.authorize(ctx, info.FullMethod); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func (ti *ServerTokenInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv any,
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if err := ti.authorize(stream.Context(), info.FullMethod); err != nil {
			return err
		}
		return handler(srv, stream)
	}
}

type ClientTokenInterceptor struct {
	token string
}

func NewClientTokenInterceptor(token string) *ClientTokenInterceptor {
	return &ClientTokenInterceptor{token: token}
}

func (ti *ClientTokenInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req any,
		reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctx = metadata.AppendToOutgoingContext(ctx, authKey, ti.token)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (ti *ClientTokenInterceptor) Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		ctx = metadata.AppendToOutgoingContext(ctx, authKey, ti.token)
		return streamer(ctx, desc, cc, method, opts...)
	}
}
