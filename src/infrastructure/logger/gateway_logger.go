package logger

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mini-market/src/infrastructure/logredact"
	"mini-market/src/infrastructure/sentry"
	"mini-market/src/infrastructure/telemetry"

	"go.uber.org/zap"
)

type GatewayLogger struct {
	*zap.Logger
}

// @inject
func NewGatewayLogger(tel *telemetry.Telemetry, gt *sentry.Client) *GatewayLogger {
	return &GatewayLogger{Logger: NewOtelZapLogger(tel, gt, "gateway").With(zap.String("scope", "gateway"))}
}

func (this *GatewayLogger) LogRequest(ctx context.Context, method, url, contentType string, body []byte) {
	fields := append(traceFields(ctx),
		zap.String("method", method),
		zap.String("url", url),
		zap.String("content_type", contentType),
		zap.String("body", requestBody(contentType, body)),
	)
	this.Info("gateway request", fields...)
}

func (this *GatewayLogger) LogResponse(ctx context.Context, method, url string, statusCode int, contentType string, body []byte, duration time.Duration) {
	fields := append(traceFields(ctx),
		zap.String("method", method),
		zap.String("url", url),
		zap.Int("status", statusCode),
		zap.Duration("duration", duration),
		zap.String("body", responseBody(contentType, body)),
	)
	this.Info("gateway response", fields...)
}

func (this *GatewayLogger) LogError(ctx context.Context, method, url string, err error, duration time.Duration) {
	fields := append(traceFields(ctx),
		zap.String("method", method),
		zap.String("url", url),
		zap.Duration("duration", duration),
		zap.Error(err),
	)
	this.Error("gateway error", fields...)
}

const maxBodyLog = 512

func requestBody(contentType string, body []byte) string {
	if strings.HasPrefix(contentType, "multipart/") {
		return "[multipart form-data]"
	}
	if isBinary(contentType) {
		return "[binary]"
	}
	return truncate(redactBody(contentType, body))
}

func responseBody(contentType string, body []byte) string {
	if isBinary(contentType) {
		return "[binary " + bytesSize(len(body)) + "]"
	}
	return truncate(redactBody(contentType, body))
}

func redactBody(contentType string, body []byte) string {
	if strings.Contains(strings.ToLower(contentType), "application/json") {
		if redacted, ok := logredact.RedactJSON(body); ok {
			return redacted
		}
	}
	return string(body)
}

func truncate(s string) string {
	runes := []rune(s)
	if len(runes) <= maxBodyLog {
		return s
	}
	return string(runes[:maxBodyLog]) + fmt.Sprintf("... [+%d chars]", len(runes)-maxBodyLog)
}

func isBinary(contentType string) bool {
	for _, t := range []string{"application/pdf", "application/octet-stream", "image/"} {
		if strings.HasPrefix(contentType, t) {
			return true
		}
	}
	return false
}

func bytesSize(n int) string {
	switch {
	case n >= 1024*1024:
		return fmt.Sprintf("%d MB", n/(1024*1024))
	case n >= 1024:
		return fmt.Sprintf("%d KB", n/1024)
	default:
		return fmt.Sprintf("%d B", n)
	}
}
