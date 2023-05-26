package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
)

func New(level string) error {
	zapConfig := &zap.Config{
		Level:    getLogLevel(level),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",
			LevelKey:   "level",
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	log, err := zapConfig.Build()
	if err != nil {
		return err
	}

	logger = log
	return nil
}

func GetLogger() *zap.Logger {
	return logger
}

func getLogLevel(level string) zap.AtomicLevel {
	switch level {
	case "INFO", "info":
		return zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "DEBUG", "debug":
		return zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "WARN", "warn":
		return zap.NewAtomicLevelAt(zapcore.WarnLevel)
	default:
		return zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
}
