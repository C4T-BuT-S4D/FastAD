package logging

import (
	"go.temporal.io/sdk/log"
	"go.uber.org/zap"
)

var (
	_ log.Logger          = TemporalAdapter{}
	_ log.WithSkipCallers = TemporalAdapter{}
	_ log.WithLogger      = TemporalAdapter{}
)

type TemporalAdapter struct {
	impl *zap.SugaredLogger
}

func NewTemporalAdapter(logger *zap.Logger) TemporalAdapter {
	return TemporalAdapter{impl: logger.Sugar()}
}

func (a TemporalAdapter) With(keyvals ...any) log.Logger {
	return TemporalAdapter{impl: a.impl.With(keyvals...)}
}

func (a TemporalAdapter) WithCallerSkip(i int) log.Logger {
	return TemporalAdapter{impl: a.impl.WithOptions(zap.AddCallerSkip(i))}
}

func (a TemporalAdapter) Debug(msg string, keyvals ...any) {
	a.impl.Debugw(msg, keyvals...)
}

func (a TemporalAdapter) Info(msg string, keyvals ...any) {
	a.impl.Infow(msg, keyvals...)
}

func (a TemporalAdapter) Warn(msg string, keyvals ...any) {
	a.impl.Warnw(msg, keyvals...)
}

func (a TemporalAdapter) Error(msg string, keyvals ...any) {
	a.impl.Errorw(msg, keyvals...)
}
