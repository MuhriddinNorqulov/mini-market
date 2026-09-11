package security

import (
	"context"
	"time"
)

type SessionDenylist interface {
	Revoke(ctx context.Context, sessionID string, ttl time.Duration) error
	IsRevoked(ctx context.Context, sessionID string) (bool, error)
}
