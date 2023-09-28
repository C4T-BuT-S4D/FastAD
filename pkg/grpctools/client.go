package grpctools

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
)

func Dial(address, userAgent string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(
		opts,
		grpc.WithUserAgent(userAgent),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.UseCompressor(gzip.Name),
		),
	)
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("dialing %s: %w", address, err)
	}
	return conn, nil
}
