package grpcext

import (
	"context"
	"errors"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnwrapStatusUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		resp, err = handler(ctx, req)
		return resp, unwrapError(err)
	}
}

func UnwrapStatusStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return unwrapError(handler(srv, ss))
	}
}

func UnwrapStatusUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return unwrapError(invoker(ctx, method, req, reply, cc, opts...))
	}
}

func UnwrapStatusStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		stream := &unwrappingStream{clientStream}
		return stream, unwrapError(err)
	}
}

type unwrappingStream struct {
	grpc.ClientStream
}

func (s *unwrappingStream) SendMsg(m any) error {
	return unwrapError(s.ClientStream.SendMsg(m))
}

func (s *unwrappingStream) RecvMsg(m any) error {
	err := s.ClientStream.RecvMsg(m)
	if err == nil || errors.Is(err, io.EOF) {
		//nolint:wrapcheck // We don't want to wrap EOF.
		return err
	}
	return unwrapError(err)
}

func unwrapError(err error) error {
	if err == nil {
		return nil
	}
	if st, ok := status.FromError(err); ok {
		return st.Err()
	}
	return status.FromContextError(err).Err()
}
