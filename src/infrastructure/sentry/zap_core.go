package sentry

import (
	"context"
	"errors"

	"mini-market/src/infrastructure/telemetry"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const CauseKey = "_gt_cause"

func Cause(err error) zap.Field {
	return zapcore.Field{Key: CauseKey, Type: zapcore.SkipType, Interface: err}
}

var _ zapcore.Core = (*Core)(nil)

type Core struct {
	client *Client
	fields []zapcore.Field
}

func NewCore(client *Client) *Core {
	return &Core{client: client}
}

func (this *Core) Enabled(level zapcore.Level) bool {
	return this.client != nil && this.client.Enabled && level >= zapcore.ErrorLevel
}

func (this *Core) With(fields []zapcore.Field) zapcore.Core {
	merged := make([]zapcore.Field, len(this.fields), len(this.fields)+len(fields))
	copy(merged, this.fields)
	merged = append(merged, fields...)
	return &Core{client: this.client, fields: merged}
}

func (this *Core) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if this.Enabled(entry.Level) {
		return ce.AddCore(entry, this)
	}
	return ce
}

func (this *Core) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	all := make([]zapcore.Field, 0, len(this.fields)+len(fields))
	all = append(all, this.fields...)
	all = append(all, fields...)

	err := extractError(all)

	if extractKind(all, err) == telemetry.ErrorKindBusiness {
		return nil
	}

	if err != nil && errors.Is(err, context.Canceled) {
		return nil
	}

	ev := sentry.NewEvent()
	ev.Level = sentry.LevelError
	ev.Message = entry.Message
	if err != nil {

		ev.SetException(err, 10)
	}
	applyFields(ev, all, this.client.sendBodies)

	this.client.capture(ev)
	return nil
}

func (this *Core) Sync() error { return nil }

func extractError(fields []zapcore.Field) error {
	var fallback error

	for _, f := range fields {
		if f.Key == CauseKey {
			if err, ok := f.Interface.(error); ok && err != nil {
				return err
			}
			continue
		}
		if f.Type == zapcore.ErrorType && fallback == nil {
			if err, ok := f.Interface.(error); ok && err != nil {
				fallback = err
			}
		}
	}

	return fallback
}

func extractKind(fields []zapcore.Field, err error) string {
	for _, f := range fields {
		if f.Key == telemetry.AttrErrorKind && f.Type == zapcore.StringType {
			return f.String
		}
	}

	kind, _ := telemetry.ErrorKindOf(err)
	return kind
}
