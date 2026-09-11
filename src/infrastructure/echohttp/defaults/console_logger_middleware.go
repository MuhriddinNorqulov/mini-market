package defaults

import (
	"errors"
	"fmt"
	"mini-market/src/core/application/response"
	"mini-market/src/infrastructure/logger"

	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	cyan   = "\033[36m"
)

type ConsoleLoggerMiddleware struct {
	log *logger.ConsoleLogger
}

// @inject
func NewConsoleLoggerMiddleware(log *logger.ConsoleLogger) *ConsoleLoggerMiddleware {
	return &ConsoleLoggerMiddleware{log: log}
}

func (m *ConsoleLoggerMiddleware) Wrap(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		req := c.Request()
		res := c.Response()

		err := next(c)

		latency := time.Since(start)

		method := req.Method
		path := req.URL.Path
		status := res.Status
		if err != nil && !res.Committed {
			if he, ok := errors.AsType[*echo.HTTPError](err); ok {
				status = he.Code
			} else {
				status = http.StatusInternalServerError
			}
		}
		timestamp := time.Now().Format("2006-01-02 15:04:05")

		msg := fmt.Sprintf(
			"%s | %s | %s | %s | %s%s",
			colorStatus(status),
			colorMethod(method),
			path,
			colorLatency(latency),
			colorTime(timestamp),
			errText(err),
		)

		m.log.Info(msg)

		return err
	}
}

func colorMethod(method string) string {
	switch method {
	case "GET":
		return wrap(green, method)
	case "POST":
		return wrap(blue, method)
	case "PUT":
		return wrap(yellow, method)
	case "DELETE":
		return wrap(red, method)
	default:
		return wrap(cyan, method)
	}
}

func colorStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return wrap(green, fmt.Sprintf("%d", code))
	case code >= 300 && code < 400:
		return wrap(cyan, fmt.Sprintf("%d", code))
	case code >= 400 && code < 500:
		return wrap(yellow, fmt.Sprintf("%d", code))
	default:
		return wrap(red, fmt.Sprintf("%d", code))
	}
}

func colorLatency(d time.Duration) string {
	ms := d.Milliseconds()
	switch {
	case ms < 100:
		return wrap(green, fmt.Sprintf("%dms", ms))
	case ms < 500:
		return wrap(yellow, fmt.Sprintf("%dms", ms))
	default:
		return wrap(red, fmt.Sprintf("%dms", ms))
	}
}

func wrap(color, s string) string {
	return color + s + reset
}

func errText(err error) string {
	if err == nil {
		return ""
	}
	if se, ok := errors.AsType[*response.SafeError](err); ok {
		return fmt.Sprintf(" | error_code=%s caller=%s error=%v", se.Code, se.Caller, se.Internal)
	}
	return fmt.Sprintf(" | error=%v", err)
}

func colorTime(t string) string {
	return wrap(cyan, t)
}
