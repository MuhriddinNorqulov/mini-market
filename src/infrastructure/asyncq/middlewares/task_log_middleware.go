package middlewares

import (
	"context"
	"fmt"
	"mini-market/src/core/domain/ports/async"
	"mini-market/src/core/domain/ports/reqctx"
	"mini-market/src/infrastructure/logger"
	"mini-market/src/infrastructure/telemetry"
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

type TaskLogMiddleware struct {
	consoleLog *logger.ConsoleLogger
	fileLog    *logger.AsyncLogger
}

// @inject
func NewTaskLogMiddleware(consoleLog *logger.ConsoleLogger, fileLog *logger.AsyncLogger) *TaskLogMiddleware {
	return &TaskLogMiddleware{consoleLog: consoleLog, fileLog: fileLog}
}

func (this *TaskLogMiddleware) Wrap(next async.AsyncTaskHandler) async.AsyncTaskHandler {
	return async.AsyncTaskHandlerFunc(func(ctx context.Context, t *async.Task) error {
		start := time.Now()

		traceID := reqctx.NewTraceID()
		if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
			traceID = sc.TraceID().String()
		}
		ctx = reqctx.WithIDs(ctx, t.TaskID, traceID)

		err := next.ProcessTask(ctx, t)

		latency := time.Since(start)
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		payload := string(t.Payload)

		status := colorWrap(colorGreen, "OK")
		errMsg := ""
		if err != nil {
			status = colorWrap(colorRed, "FAIL")
			errMsg = fmt.Sprintf(" | %serr=%v%s", colorRed, err, colorReset)
		}

		msg := fmt.Sprintf("%s | %s | %s | %s | %s%s",
			status,
			colorWrap(colorBlue, string(t.TaskType)),
			colorWrap(colorGray, t.TaskID),
			colorLatencyTask(latency),
			colorWrap(colorCyan, timestamp),
			errMsg,
		)
		if payload != "" && payload != "{}" && payload != "[]" {
			msg += fmt.Sprintf(" | %spayload=%s%s", colorYellow, payload, colorReset)
		}

		this.consoleLog.Info(msg)

		fields := []zap.Field{
			zap.String(telemetry.AttrTaskType, string(t.TaskType)),
			zap.String(telemetry.AttrTaskID, t.TaskID),
			zap.String("payload", payload),
			zap.Int64("latency_ms", latency.Milliseconds()),
			logger.Ctx(ctx),
			zap.String("trace_id", traceID),
		}

		if err != nil {
			fields = append(fields, zap.String("task_status", "FAIL"), zap.Error(err))
			this.fileLog.Error("task_completed", fields...)
		} else {
			fields = append(fields, zap.String("task_status", "OK"))
			this.fileLog.Info("task_completed", fields...)
		}

		return err
	})
}

func colorWrap(color, s string) string {
	return color + s + colorReset
}

func colorLatencyTask(d time.Duration) string {
	ms := d.Milliseconds()
	switch {
	case ms < 500:
		return colorWrap(colorGreen, fmt.Sprintf("%dms", ms))
	case ms < 3000:
		return colorWrap(colorYellow, fmt.Sprintf("%dms", ms))
	default:
		return colorWrap(colorRed, fmt.Sprintf("%.1fs", d.Seconds()))
	}
}
