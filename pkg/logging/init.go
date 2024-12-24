package logging

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/c4t-but-s4d/fastad/pkg/baseconfig"
)

var initOnce sync.Once

type CheckedCloser interface {
	Close()
}

type checkedCloserImpl func()

func (c checkedCloserImpl) Close() {
	c()
}

func Init() CheckedCloser {
	initOnce.Do(func() {
		cfg := baseconfig.MustSetupAll(&Config{}, baseconfig.WithEnvPrefix("FASTAD_LOG"))
		level := parseLogLevel(cfg.Level, zap.DebugLevel)

		devEncoder := zap.NewDevelopmentEncoderConfig()
		devEncoder.EncodeTime = zapcore.ISO8601TimeEncoder
		devEncoder.EncodeLevel = zapcore.CapitalColorLevelEncoder

		consoleEncoder := zapcore.NewConsoleEncoder(devEncoder)

		stderr := zapcore.Lock(os.Stderr)
		core := zapcore.NewCore(consoleEncoder, stderr, level)

		zap.ReplaceGlobals(zap.New(core, zap.WithCaller(true)))
	})

	return checkedCloserImpl(func() {
		if err := zap.L().Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
			fmt.Printf("failed to sync logger: %v\n", err)
		}
	})
}

func parseLogLevel(levelStr string, defaultLevel zapcore.Level) zapcore.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn", "warning":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "panic", "fatal":
		return zap.FatalLevel
	default:
		return defaultLevel
	}
}
