package defaults

import (
	"fmt"
	"mini-market/src/core/domain/ports/reqctx"
	"mini-market/src/infrastructure/logger"
	"mini-market/src/infrastructure/logredact"
	"mini-market/src/infrastructure/sentry"
	"mini-market/src/infrastructure/telemetry"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type RecoveryMiddleware struct {
	logger *logger.HttpLogger
}

// @inject
func NewRecoveryMiddleware(logger *logger.HttpLogger) *RecoveryMiddleware {
	return &RecoveryMiddleware{logger: logger}
}

func (m *RecoveryMiddleware) Wrap(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		start := time.Now()

		defer func() {
			if r := recover(); r != nil {
				req := c.Request()
				ctx := req.Context()

				panicErr, ok := r.(error)
				if !ok {
					panicErr = fmt.Errorf("panic: %v", r)
				}

				panicError := telemetry.NewPanicError(panicErr, 4)

				span := trace.SpanFromContext(ctx)
				span.SetAttributes(attribute.String(telemetry.AttrErrorKind, telemetry.ErrorKindPanic))
				span.RecordError(panicErr)
				span.SetStatus(codes.Error, panicErr.Error())

				m.logger.Error("panic "+requestMessage(c, req),
					zap.String("request_id", reqctx.GetRequestID(ctx)),
					zap.String("trace_id", reqctx.GetTraceID(ctx)),
					zap.String("method", req.Method),
					zap.String("path", req.URL.Path),
					zap.String("query", logredact.SanitizeQuery(req.URL.RawQuery)),
					zap.String("ip", c.RealIP()),
					zap.String("user_agent", req.UserAgent()),
					zap.Int("status", http.StatusInternalServerError),
					zap.Duration("latency", time.Since(start)),

					zap.String("panic", fmt.Sprintf("%v", r)),

					zap.String(telemetry.AttrErrorKind, telemetry.ErrorKindPanic),
					zap.Error(panicErr),

					sentry.Cause(panicError),
					logger.Ctx(ctx),
				)

				err = echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
			}
		}()

		return next(c)
	}
}
