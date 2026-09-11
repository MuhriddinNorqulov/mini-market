package logger

import (
	"os"

	"mini-market/src/infrastructure/sentry"
	"mini-market/src/infrastructure/telemetry"

	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewOtelZapLogger(tel *telemetry.Telemetry, gt *sentry.Client, scope string) *zap.Logger {
	var core zapcore.Core

	if tel.Enabled {
		core = otelzap.NewCore(scope, otelzap.WithLoggerProvider(tel.LoggerProvider()))
	} else {
		encCfg := zap.NewProductionEncoderConfig()
		encCfg.TimeKey = "ts"
		encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encCfg),
			zapcore.AddSync(os.Stdout),
			zapcore.InfoLevel,
		)
	}

	if gt != nil && gt.Enabled {
		core = zapcore.NewTee(core, sentry.NewCore(gt))
	}

	return zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
}
