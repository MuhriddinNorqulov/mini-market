package defaults

import (
	"errors"

	"mini-market/src/core/application/response"
	"mini-market/src/infrastructure/telemetry"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ErrorRecorderMiddleware struct{}

// @inject
func NewErrorRecorderMiddleware() *ErrorRecorderMiddleware {
	return &ErrorRecorderMiddleware{}
}

func (this *ErrorRecorderMiddleware) Wrap(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		span := trace.SpanFromContext(c.Request().Context())
		kind, failure := requestErrorKind(err)

		if kind != "" {
			span.SetAttributes(attribute.String(telemetry.AttrErrorKind, kind))
		}
		if !failure {
			return err
		}

		if se, ok := errors.AsType[*response.SafeError](err); ok {
			span.SetAttributes(
				attribute.String(telemetry.AttrErrorCode, string(se.Code)),
				attribute.String(telemetry.AttrCaller, se.Caller),
			)
			span.RecordError(se.Internal)
			span.SetStatus(codes.Error, string(se.Code))
			return err
		}

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
}

func requestErrorKind(err error) (kind string, failure bool) {
	kind, failure = telemetry.ErrorKindOf(err)
	if he, ok := errors.AsType[*echo.HTTPError](err); ok && he.Code >= 400 && he.Code < 500 {
		return telemetry.ErrorKindBusiness, false
	}
	return kind, failure
}
