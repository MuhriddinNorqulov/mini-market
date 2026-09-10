package reqctx

import "context"

type deviceIDKey struct{}

func WithDeviceID(ctx context.Context, deviceID string) context.Context {
	return context.WithValue(ctx, deviceIDKey{}, deviceID)
}

func DeviceID(ctx context.Context) string {
	v, _ := ctx.Value(deviceIDKey{}).(string)
	return v
}

type sessionIDKey struct{}

func WithSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, sessionIDKey{}, sessionID)
}

func SessionID(ctx context.Context) string {
	v, _ := ctx.Value(sessionIDKey{}).(string)
	return v
}
