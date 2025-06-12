package config

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerConfig struct {
	IsProduction bool
	LogLevel     string
	LogFile      string
	MaxSize      int
	MaxBackups   int
	MaxAge       int
}

func InitLogger(cfg LoggerConfig) (*zap.Logger, error) {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		level = zap.InfoLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	if !cfg.IsProduction {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	}
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	logDir := filepath.Dir(cfg.LogFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	fileSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.LogFile,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
	})

	syncers := []zapcore.WriteSyncer{zapcore.AddSync(os.Stdout), fileSyncer}
	if cfg.IsProduction {
		syncers = []zapcore.WriteSyncer{fileSyncer}
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.NewMultiWriteSyncer(syncers...),
		zap.NewAtomicLevelAt(level),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	zap.AddStacktrace(zap.ErrorLevel)
	return logger, nil
}

func DefaultConfig(isProduction bool) LoggerConfig {
	return LoggerConfig{
		IsProduction: isProduction,
		LogLevel:     "info",
		LogFile:      "logs/app.log",
		MaxSize:      10,
		MaxBackups:   5,
		MaxAge:       30,
	}
}

func WithRequestID(logger *zap.Logger, requestID string) *zap.Logger {
	return logger.With(zap.String("request_id", requestID))
}
