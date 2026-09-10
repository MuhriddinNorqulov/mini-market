package logger

import (
	"context"

	"mini-market/src/core/domain/ports/reqctx"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func Ctx(ctx context.Context) zap.Field {
	if !trace.SpanContextFromContext(ctx).IsValid() {
		return zap.Skip()
	}
	return zap.Any("ctx", ctx)
}

func traceFields(ctx context.Context) []zap.Field {
	fields := []zap.Field{Ctx(ctx)}
	if id := reqctx.GetRequestID(ctx); id != "" {
		fields = append(fields, zap.String("request_id", id))
	}
	if id := reqctx.GetTraceID(ctx); id != "" {
		fields = append(fields, zap.String("trace_id", id))
	}
	return fields
}
