package grpcext

import (
	"runtime/debug"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	// Enable gzip compression.
	_ "google.golang.org/grpc/encoding/gzip"

	"github.com/c4t-but-s4d/fastad/internal/pinger"
	"github.com/c4t-but-s4d/fastad/pkg/logging"
	pingerpb "github.com/c4t-but-s4d/fastad/pkg/proto/pinger"
)

type ServerConfig struct {
	installation string
}

func GetServerConfig(opts ...ServerOption) *ServerConfig {
	cfg := &ServerConfig{
		installation: "unknown",
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

type ServerOption func(*ServerConfig)

func WithServerInstallation(name string) ServerOption {
	return func(cfg *ServerConfig) {
		cfg.installation = name
	}
}

func NewServer(opts ...ServerOption) *grpc.Server {
	cfg := GetServerConfig(opts...)

	rpcLogger := zap.L().Named("grpc_server").With(
		zap.String("installation", cfg.installation),
	)

	var coounterOpts []grpcprom.CounterOption
	histogramOpts := []grpcprom.HistogramOption{
		grpcprom.WithHistogramBuckets(
			prometheus.ExponentialBucketsRange(0.001, 120, 30),
		),
	}
	if cfg.installation != "" {
		coounterOpts = append(
			coounterOpts,
			grpcprom.WithConstLabels(prometheus.Labels{"installation": cfg.installation}),
		)
		histogramOpts = append(
			histogramOpts,
			grpcprom.WithHistogramConstLabels(prometheus.Labels{"installation": cfg.installation}),
		)
	}

	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(histogramOpts...),
		grpcprom.WithServerCounterOptions(coounterOpts...),
	)
	prometheus.MustRegister(srvMetrics)

	panicsCounterOpts := prometheus.CounterOpts{
		Name: "grpc_server_panics_recovered_total",
		Help: "Total number of gRPC requests recovered from internal panic.",
	}
	if cfg.installation != "" {
		panicsCounterOpts.ConstLabels = map[string]string{"installation": cfg.installation}
	}

	panicsTotal := promauto.NewCounter(panicsCounterOpts)
	grpcPanicRecoveryHandler := func(p any) (err error) {
		panicsTotal.Inc()
		rpcLogger.Error("recovered from panic", zap.Any("panic", p), zap.ByteString("stack", debug.Stack()))
		return status.Errorf(codes.Internal, "%s", p)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			srvMetrics.UnaryServerInterceptor(), // TODO: exemplars.
			grpclog.UnaryServerInterceptor(
				logging.InterceptorLogger(rpcLogger),
				grpclog.WithLevels(logging.GRPCCodeToLevel),
			), // TODO: tracing.
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
			UnwrapStatusUnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			srvMetrics.StreamServerInterceptor(), // TODO: exemplars.
			grpclog.StreamServerInterceptor(
				logging.InterceptorLogger(rpcLogger),
				grpclog.WithLevels(logging.GRPCCodeToLevel),
			), // TODO: tracing.
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
			UnwrapStatusStreamServerInterceptor(),
		),
	)
	reflection.Register(s)
	pingerpb.RegisterPingerServiceServer(s, pinger.New())
	return s
}
