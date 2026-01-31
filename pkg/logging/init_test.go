//nolint:testpackage // Testing unexported function parseLogLevel
package logging

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name         string
		levelStr     string
		defaultLevel zapcore.Level
		want         zapcore.Level
	}{
		// Valid level strings
		{
			name:         "debug level",
			levelStr:     "debug",
			defaultLevel: zap.InfoLevel,
			want:         zap.DebugLevel,
		},
		{
			name:         "DEBUG level uppercase",
			levelStr:     "DEBUG",
			defaultLevel: zap.InfoLevel,
			want:         zap.DebugLevel,
		},
		{
			name:         "Debug level mixed case",
			levelStr:     "Debug",
			defaultLevel: zap.InfoLevel,
			want:         zap.DebugLevel,
		},
		{
			name:         "info level",
			levelStr:     "info",
			defaultLevel: zap.DebugLevel,
			want:         zap.InfoLevel,
		},
		{
			name:         "warn level",
			levelStr:     "warn",
			defaultLevel: zap.InfoLevel,
			want:         zap.WarnLevel,
		},
		{
			name:         "warning level alias",
			levelStr:     "warning",
			defaultLevel: zap.InfoLevel,
			want:         zap.WarnLevel,
		},
		{
			name:         "error level",
			levelStr:     "error",
			defaultLevel: zap.InfoLevel,
			want:         zap.ErrorLevel,
		},
		{
			name:         "panic level maps to fatal",
			levelStr:     "panic",
			defaultLevel: zap.InfoLevel,
			want:         zap.FatalLevel,
		},
		{
			name:         "fatal level",
			levelStr:     "fatal",
			defaultLevel: zap.InfoLevel,
			want:         zap.FatalLevel,
		},
		// Invalid level strings - should return default
		{
			name:         "empty string returns default",
			levelStr:     "",
			defaultLevel: zap.InfoLevel,
			want:         zap.InfoLevel,
		},
		{
			name:         "unknown level returns default",
			levelStr:     "unknown",
			defaultLevel: zap.WarnLevel,
			want:         zap.WarnLevel,
		},
		{
			name:         "typo level returns default",
			levelStr:     "debuf",
			defaultLevel: zap.ErrorLevel,
			want:         zap.ErrorLevel,
		},
		{
			name:         "numeric string returns default",
			levelStr:     "123",
			defaultLevel: zap.InfoLevel,
			want:         zap.InfoLevel,
		},
		{
			name:         "whitespace returns default",
			levelStr:     "   ",
			defaultLevel: zap.DebugLevel,
			want:         zap.DebugLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseLogLevel(tt.levelStr, tt.defaultLevel)
			if got != tt.want {
				t.Errorf("parseLogLevel(%q, %v) = %v, want %v",
					tt.levelStr, tt.defaultLevel, got, tt.want)
			}
		})
	}
}

func TestCheckedCloserImpl(t *testing.T) {
	// Test that the closer implementation works
	called := false
	closer := checkedCloserImpl(func() {
		called = true
	})

	closer.Close()

	if !called {
		t.Error("expected closer function to be called")
	}
}

func TestCheckedCloser_Interface(_ *testing.T) {
	var _ CheckedCloser = checkedCloserImpl(func() {})
}
