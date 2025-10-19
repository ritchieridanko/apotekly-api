package logger

import (
	"log"
	"os"
	"strings"

	"github.com/ritchieridanko/apotekly-api/auth/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewProvider(cfg *configs.App) *zap.Logger {
	encoderCfg := zapcore.EncoderConfig{
		LevelKey:       "level",
		MessageKey:     "message",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var level zapcore.Level
	if env := strings.ToLower(strings.TrimSpace(cfg.Env)); env == "production" {
		level = zapcore.InfoLevel
	} else {
		level = zapcore.DebugLevel
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		level,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	log.Println("✅ initialized logger")
	return logger
}
