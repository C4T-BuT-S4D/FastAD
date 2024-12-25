package logging

import (
	"context"
	"fmt"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

// InterceptorLogger adapts zap logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *zap.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(_ context.Context, lvl grpclog.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			keyStr, ok := key.(string)
			if !ok || i+1 >= len(fields) {
				continue
			}
			value := fields[i+1]

			switch v := value.(type) {
			case string:
				f = append(f, zap.String(keyStr, v))
			case int:
				f = append(f, zap.Int(keyStr, v))
			case bool:
				f = append(f, zap.Bool(keyStr, v))
			default:
				f = append(f, zap.Any(keyStr, v))
			}
		}

		logger := l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch lvl {
		case grpclog.LevelDebug:
			logger.Debug(msg)
		case grpclog.LevelInfo:
			logger.Info(msg)
		case grpclog.LevelWarn:
			logger.Warn(msg)
		case grpclog.LevelError:
			logger.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

func GRPCCodeToLevel(code codes.Code) grpclog.Level {
	switch code {
	case codes.OK:
		return grpclog.LevelDebug

	case codes.NotFound, codes.Canceled, codes.AlreadyExists, codes.InvalidArgument, codes.Unauthenticated:
		return grpclog.LevelInfo

	case codes.DeadlineExceeded, codes.PermissionDenied, codes.ResourceExhausted, codes.FailedPrecondition, codes.Aborted,
		codes.OutOfRange, codes.Unavailable:
		return grpclog.LevelWarn

	case codes.Unknown, codes.Unimplemented, codes.Internal, codes.DataLoss:
		return grpclog.LevelError

	default:
		return grpclog.LevelError
	}
}
