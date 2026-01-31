package grpcext

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func RunServer(runCtx, shutdownCtx context.Context, server *grpc.Server, addr string) error {
	logger := zap.L().With(zap.String("address", addr))

	var lc net.ListenConfig
	lis, err := lc.Listen(runCtx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("listening: %w", err)
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		const stopTimeout = 5 * time.Second

		<-runCtx.Done()

		logger.Info("stopping gRPC server")
		ctx, cancel := context.WithTimeout(shutdownCtx, stopTimeout)
		defer cancel()

		stopped := make(chan struct{})
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(stopped)
			server.GracefulStop()
		}()

		select {
		case <-stopped:
			logger.Info("gRPC server stopped gracefully")
		case <-ctx.Done():
			server.Stop()
			logger.Warn("gRPC server stopped forcefully")
		}
	}()

	logger.Info("serving gRPC")
	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("serving: %w", err)
	}
	return nil
}
