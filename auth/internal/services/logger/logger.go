package logger

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/constants"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{logger}
}

func (l *Logger) Log(ctx *gin.Context, level zapcore.Level, message string, statusCode int, additionalFields ...zap.Field) {
	traceID := trace.SpanFromContext(ctx.Request.Context()).SpanContext().TraceID().String()
	requestID := fmt.Sprintf("%s", ctx.Request.Context().Value(constants.CtxKeyRequestID))

	fields := []zap.Field{
		zap.String("timestamp", time.Now().UTC().Format(time.RFC3339)),
		zap.String("request_id", requestID),
		zap.String("trace_id", traceID),
		zap.String("client_ip", ctx.ClientIP()),
		zap.String("user_agent", ctx.Request.UserAgent()),
		zap.String("method", ctx.Request.Method),
		zap.String("path", ctx.Request.URL.Path),
		zap.Int("status", statusCode),
	}

	fields = append(fields, additionalFields...)
	l.logger.Log(level, message, fields...)
}
