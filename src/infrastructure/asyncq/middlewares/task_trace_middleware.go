package middlewares

import (
	"context"
	"fmt"

	"mini-market/src/core/domain/ports/async"
	"mini-market/src/infrastructure/asyncq/envelope"
	"mini-market/src/infrastructure/logger"
	"mini-market/src/infrastructure/sentry"
	"mini-market/src/infrastructure/telemetry"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type TaskTraceMiddleware struct {
	log *logger.AsyncLogger
}

// @inject
func NewTaskTraceMiddleware(log *logger.AsyncLogger) *TaskTraceMiddleware {
	return &TaskTraceMiddleware{log: log}
}

func (this *TaskTraceMiddleware) Wrap(next async.AsyncTaskHandler) async.AsyncTaskHandler {
	return async.AsyncTaskHandlerFunc(func(ctx context.Context, t *async.Task) (returnErr error) {

		ctx, t.Payload = envelope.Unwrap(ctx, t.Payload)

		ctx, span := otel.Tracer("asynq").Start(ctx, string(t.TaskType),
			trace.WithSpanKind(trace.SpanKindConsumer),
			trace.WithAttributes(
				attribute.String(telemetry.AttrTaskType, string(t.TaskType)),
				attribute.String(telemetry.AttrTaskID, t.TaskID),
			),
		)
		defer span.End()

		defer func() {
			if r := recover(); r != nil {

				panicErr, ok := r.(error)
				if !ok {
					panicErr = fmt.Errorf("panic: %v", r)
				}

				panicError := telemetry.NewPanicError(panicErr, 4)

				span.SetAttributes(attribute.String(telemetry.AttrErrorKind, telemetry.ErrorKindPanic))
				span.RecordError(panicErr)
				span.SetStatus(codes.Error, panicErr.Error())

				this.log.Error("task panic",
					zap.String(telemetry.AttrTaskType, string(t.TaskType)),
					zap.String(telemetry.AttrTaskID, t.TaskID),
					zap.String("panic", fmt.Sprintf("%v", r)),
					zap.Error(panicErr),
					zap.String(telemetry.AttrErrorKind, telemetry.ErrorKindPanic),
					sentry.Cause(panicError),
					logger.Ctx(ctx),
				)

				span.End()

				panic(r)
			}
		}()

		err := next.ProcessTask(ctx, t)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		return err
	})
}
