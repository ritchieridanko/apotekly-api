package constants

import "go.uber.org/zap/zapcore"

const (
	LogLevelError zapcore.Level = zapcore.ErrorLevel
	LogLevelInfo  zapcore.Level = zapcore.InfoLevel
	LogLevelWarn  zapcore.Level = zapcore.WarnLevel
)
