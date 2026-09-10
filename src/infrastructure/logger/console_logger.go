package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ConsoleLogger struct {
	*zap.Logger
}

// @inject
func NewConsoleLogger() *ConsoleLogger {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: true,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeDuration: zapcore.StringDurationEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	l, _ := cfg.Build()

	return &ConsoleLogger{Logger: l}
}
