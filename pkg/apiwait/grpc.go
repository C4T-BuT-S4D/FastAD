package apiwait

import (
	"context"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
)

// GRPC waits for a gRPC server to become available at the given address.
// It attempts to establish a TCP connection to verify the server is listening.
func GRPC(ctx context.Context, address string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	zap.L().Debug("waiting for gRPC server", zap.String("address", address))

	dialer := net.Dialer{}
	for {
		if ctx.Err() != nil {
			return fmt.Errorf("waiting for gRPC server at %s: %w", address, ctx.Err())
		}

		conn, err := dialer.DialContext(ctx, "tcp", address)
		if err != nil {
			zap.L().Debug(
				"gRPC server not ready",
				zap.String("address", address),
				zap.Error(err),
			)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		_ = conn.Close()
		zap.L().Debug("gRPC server is ready", zap.String("address", address))
		return nil
	}
}
