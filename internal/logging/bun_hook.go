package logging

import (
	"context"
	"time"

	"github.com/uptrace/bun"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	OperationFieldName     = "operation"
	OperationTimeFieldName = "operation_time_ms"
	SlowQueryThreshold     = time.Millisecond * 100
)

type QueryHook struct {
	bun.QueryHook

	logger *zap.Logger
}

func NewQueryHook(logger *zap.Logger) QueryHook {
	return QueryHook{logger: logger.WithOptions(zap.AddCallerSkip(6))}
}

func (qh QueryHook) BeforeQuery(ctx context.Context, _ *bun.QueryEvent) context.Context {
	return ctx
}

func (qh QueryHook) AfterQuery(_ context.Context, event *bun.QueryEvent) {
	queryDuration := time.Since(event.StartTime)
	fields := []zapcore.Field{
		zap.String(OperationFieldName, event.Operation()),
		zap.Int64(OperationTimeFieldName, queryDuration.Milliseconds()),
	}

	// Errors will always be logged
	if event.Err != nil {
		fields = append(fields, zap.Error(event.Err))
		qh.logger.Error(event.Query, fields...)
		return
	}

	// Queries over a slow time duration
	// will be logged as debug
	if queryDuration >= SlowQueryThreshold {
		qh.logger.Warn(event.Query, fields...)
	} else {
		qh.logger.Debug(event.Query, fields...)
	}
}

func AddBunQueryHook(db *bun.DB) {
	db.AddQueryHook(NewQueryHook(zap.L()))
}
