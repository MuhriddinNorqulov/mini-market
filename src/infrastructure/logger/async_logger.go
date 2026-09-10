package logger

import (
	"mini-market/src/infrastructure/sentry"
	"mini-market/src/infrastructure/telemetry"

	"go.uber.org/zap"
)

type AsyncLogger struct {
	*zap.Logger
}

// @inject
func NewAsyncLogger(tel *telemetry.Telemetry, gt *sentry.Client) *AsyncLogger {
	return &AsyncLogger{Logger: NewOtelZapLogger(tel, gt, "async").With(zap.String("scope", "async"))}
}
