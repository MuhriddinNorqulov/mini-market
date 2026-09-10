package sentry

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"mini-market/src/infrastructure/telemetry"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap/zapcore"
)

var tagKeys = map[string]bool{
	telemetry.AttrErrorKind: true,
	telemetry.AttrErrorCode: true,
	telemetry.AttrTaskType:  true,
	"scope":                 true,
	"status":                true,
	"method":                true,

	"trace_id": true,
}

var piiKeys = map[string]bool{
	"req_body":   true,
	"res_body":   true,
	"payload":    true,
	"user_agent": true,
	"ip":         true,
}

func applyFields(ev *sentry.Event, fields []zapcore.Field, sendBodies bool) {
	extra := make(map[string]any, len(fields))

	for _, f := range fields {

		if f.Type == zapcore.SkipType {
			continue
		}

		if _, ok := f.Interface.(context.Context); ok {
			continue
		}
		if piiKeys[f.Key] && !sendBodies {
			continue
		}

		if f.Key == telemetry.AttrUserID {
			ev.User = sentry.User{ID: fieldToString(f)}
			continue
		}
		if tagKeys[f.Key] {
			ev.Tags[f.Key] = fieldToString(f)
			continue
		}

		extra[f.Key] = fieldValue(f)
	}

	if len(extra) > 0 {

		ev.Contexts["log"] = extra
	}
}

func fieldValue(f zapcore.Field) any {
	switch f.Type {
	case zapcore.StringType:
		return f.String
	case zapcore.BoolType:
		return f.Integer == 1
	case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type:
		return f.Integer
	case zapcore.Uint64Type, zapcore.Uint32Type, zapcore.Uint16Type, zapcore.Uint8Type:
		return uint64(f.Integer)
	case zapcore.Float64Type:
		return math.Float64frombits(uint64(f.Integer))
	case zapcore.Float32Type:
		return float64(math.Float32frombits(uint32(f.Integer)))
	case zapcore.DurationType:
		return time.Duration(f.Integer).String()
	case zapcore.TimeFullType:
		if t, ok := f.Interface.(time.Time); ok {
			return t.Format(time.RFC3339Nano)
		}
		return fmt.Sprintf("%v", f.Interface)
	case zapcore.TimeType:
		loc := time.UTC
		if l, ok := f.Interface.(*time.Location); ok && l != nil {
			loc = l
		}
		return time.Unix(0, f.Integer).In(loc).Format(time.RFC3339Nano)
	case zapcore.ErrorType:
		if err, ok := f.Interface.(error); ok && err != nil {
			return err.Error()
		}
		return ""
	case zapcore.StringerType:
		if s, ok := f.Interface.(fmt.Stringer); ok {
			return s.String()
		}
		return fmt.Sprintf("%v", f.Interface)
	case zapcore.ByteStringType:
		if b, ok := f.Interface.([]byte); ok {
			return string(b)
		}
		return fmt.Sprintf("%v", f.Interface)
	default:
		if f.Interface != nil {
			return fmt.Sprintf("%+v", f.Interface)
		}
		return f.Integer
	}
}

func fieldToString(f zapcore.Field) string {
	switch v := fieldValue(f).(type) {
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}
