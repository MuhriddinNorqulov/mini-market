package reqctx

import (
	"context"

	"mini-market/src/core/domain/entity/enum"

	"github.com/google/uuid"
)

type ctxKey int

const (
	keyRequestID ctxKey = iota
	keyTraceID
	keyUserID
	keyUserRole
	keyLang
)

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyRequestID, id)
}

func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyTraceID, id)
}

func GetRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(keyRequestID).(string); ok {
		return v
	}
	return ""
}

func GetTraceID(ctx context.Context) string {
	if v, ok := ctx.Value(keyTraceID).(string); ok {
		return v
	}
	return ""
}

func WithIDs(ctx context.Context, requestID, traceID string) context.Context {
	ctx = WithRequestID(ctx, requestID)
	ctx = WithTraceID(ctx, traceID)
	return ctx
}

func NewTraceID() string {
	return uuid.NewString()
}

func WithUser(ctx context.Context, userID uint, role string) context.Context {
	ctx = context.WithValue(ctx, keyUserID, userID)
	return context.WithValue(ctx, keyUserRole, role)
}

func GetUserID(ctx context.Context) (uint, bool) {
	v, ok := ctx.Value(keyUserID).(uint)
	return v, ok
}

func GetUserRole(ctx context.Context) string {
	if v, ok := ctx.Value(keyUserRole).(string); ok {
		return v
	}
	return ""
}

func WithLang(ctx context.Context, lang enum.Lang) context.Context {
	return context.WithValue(ctx, keyLang, lang)
}

func GetLang(ctx context.Context) enum.Lang {
	if v, ok := ctx.Value(keyLang).(enum.Lang); ok {
		return v
	}
	return enum.DefaultLang
}
