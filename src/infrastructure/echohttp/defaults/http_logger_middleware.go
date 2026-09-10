package defaults

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mini-market/src/core/application/response"
	"mini-market/src/core/domain/ports/reqctx"
	"mini-market/src/infrastructure/logger"
	"mini-market/src/infrastructure/logredact"
	"mini-market/src/infrastructure/sentry"
	"mini-market/src/infrastructure/telemetry"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type HttpLoggerMiddleware struct {
	logger *logger.HttpLogger
}

// @inject
func NewHttpLoggerMiddleware(logger *logger.HttpLogger) *HttpLoggerMiddleware {
	return &HttpLoggerMiddleware{logger: logger}
}

func (m *HttpLoggerMiddleware) Wrap(next echo.HandlerFunc) echo.HandlerFunc {
	const maxBodyBytes = 64 * 1024

	return func(c echo.Context) error {
		req := c.Request()
		start := time.Now()

		reqID := req.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}

		traceID := reqctx.NewTraceID()
		if sc := trace.SpanContextFromContext(req.Context()); sc.IsValid() {
			traceID = sc.TraceID().String()
		}

		ctx := reqctx.WithIDs(req.Context(), reqID, traceID)
		c.SetRequest(req.WithContext(ctx))

		req = c.Request()

		lg := m.logger.With(zap.String("request_id", reqID), zap.String("trace_id", traceID))

		baseFields := func() []zap.Field {
			fields := []zap.Field{
				zap.String("request_id", reqID),
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.String("query", logredact.SanitizeQuery(req.URL.RawQuery)),
				zap.String("ip", c.RealIP()),
				zap.String("user_agent", req.UserAgent()),
				logger.Ctx(req.Context()),
			}

			if userID, ok := reqctx.GetUserID(req.Context()); ok {
				fields = append(fields, zap.Uint(telemetry.AttrUserID, userID))
			}
			return fields
		}

		var reqBody string
		contentType := req.Header.Get(echo.HeaderContentType)
		if req.Body != nil && req.Body != http.NoBody {
			if strings.HasPrefix(contentType, "multipart/") {
				reqBody = summarizeMultipart(req)
			} else {
				b, _ := io.ReadAll(io.LimitReader(req.Body, maxBodyBytes))
				reqBody = bodyForLog(contentType, b)
				req.Body = io.NopCloser(bytes.NewReader(b))
			}
		}

		rec := &responseBodyRecorder{
			ResponseWriter: c.Response().Writer,
			limit:          maxBodyBytes,
		}
		c.Response().Writer = rec

		err := next(c)

		req = c.Request()

		latency := time.Since(start)

		kind, failure := requestErrorKind(err)

		status := c.Response().Status
		if err != nil && !c.Response().Committed {
			if he, ok := errors.AsType[*echo.HTTPError](err); ok {
				status = he.Code
			} else if failure {
				status = http.StatusInternalServerError
			}
		}

		resBody := bodyForLog(c.Response().Header().Get(echo.HeaderContentType), rec.body.Bytes())

		fields := append(baseFields(),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("req_body", reqBody),
			zap.String("res_body", resBody),
		)

		isGet := req.Method == http.MethodGet
		hasError := failure || status >= 400

		if isGet && !hasError {
			return nil
		}

		switch {
		case kind != "":
			fields = append(fields, zap.String(telemetry.AttrErrorKind, kind))
		case status >= 400:
			fields = append(fields, zap.String(telemetry.AttrErrorKind, telemetry.ErrorKindBusiness))
		}

		if failure {
			if se, ok := errors.AsType[*response.SafeError](err); ok {
				fields = append(fields,
					zap.String(telemetry.AttrErrorCode, string(se.Code)),
					zap.String(telemetry.AttrCaller, se.Caller),

					zap.Error(se.Internal),

					sentry.Cause(err),
				)
			} else {
				fields = append(fields, zap.Error(err))
			}
		}

		msg := requestMessage(c, req)

		switch {
		case failure:
			lg.Error(msg, fields...)
		case status == http.StatusUnauthorized:
			lg.Info(msg, fields...)
		case status >= 400:
			lg.Warn(msg, fields...)
		default:
			lg.Info(msg, fields...)
		}

		if err != nil {
			if kind == "" {

				return nil
			}
			return err
		}

		return nil
	}
}

func requestMessage(c echo.Context, req *http.Request) string {
	route := c.Path()
	if route == "" {
		route = req.URL.Path
	}
	return req.Method + " " + route
}

type responseBodyRecorder struct {
	http.ResponseWriter
	body  bytes.Buffer
	limit int64
}

func (r *responseBodyRecorder) Write(b []byte) (int, error) {
	if int64(r.body.Len()) < r.limit {
		remain := r.limit - int64(r.body.Len())
		if int64(len(b)) > remain {
			r.body.Write(b[:remain])
		} else {
			r.body.Write(b)
		}
	}
	return r.ResponseWriter.Write(b)
}

func (r *responseBodyRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return h.Hijack()
}

const maxLoggedBodyChars = 2000

func bodyForLog(contentType string, raw []byte) string {
	if len(raw) == 0 {
		return ""
	}

	ct := strings.ToLower(contentType)

	if strings.Contains(ct, "application/x-www-form-urlencoded") {
		return truncate(logredact.SanitizeQuery(string(raw)), maxLoggedBodyChars)
	}

	if !strings.Contains(ct, "application/json") {
		return truncate(string(raw), maxLoggedBodyChars)
	}

	redacted, ok := logredact.RedactJSON(raw)
	if !ok {
		return truncate(string(raw), maxLoggedBodyChars)
	}
	return truncate(redacted, maxLoggedBodyChars)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "...(truncated)"
}

func summarizeMultipart(req *http.Request) string {
	ct := req.Header.Get("Content-Type")
	_, params, err := mime.ParseMediaType(ct)
	if err != nil {
		return "multipart (unparseable)"
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "multipart (read error)"
	}
	req.Body = io.NopCloser(bytes.NewReader(body))

	reader := multipart.NewReader(bytes.NewReader(body), params["boundary"])
	var parts []map[string]string
	for {
		part, err := reader.NextPart()
		if err != nil {
			break
		}
		info := map[string]string{
			"name": part.FormName(),
		}
		if fn := part.FileName(); fn != "" {
			info["filename"] = fn
			info["content_type"] = part.Header.Get("Content-Type")
			size, _ := io.Copy(io.Discard, part)
			info["size"] = fmt.Sprintf("%d", size)
		}
		_ = part.Close()
		parts = append(parts, info)
	}

	if len(parts) == 0 {
		return ""
	}

	b, err := json.Marshal(parts)
	if err != nil {
		return "multipart (unsummarizable)"
	}
	return truncate(string(b), maxLoggedBodyChars)
}
