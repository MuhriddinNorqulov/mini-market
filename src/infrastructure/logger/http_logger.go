package logger

import (
	"mini-market/src/infrastructure/sentry"
	"mini-market/src/infrastructure/telemetry"

	"go.uber.org/zap"
)

type HttpLogger struct {
	*zap.Logger
}

// @inject
func NewHttpLogger(tel *telemetry.Telemetry, gt *sentry.Client) *HttpLogger {
	return &HttpLogger{Logger: NewOtelZapLogger(tel, gt, "http").With(zap.String("scope", "http"))}
}
