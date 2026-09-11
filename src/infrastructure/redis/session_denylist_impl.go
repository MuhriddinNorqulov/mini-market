package redis

import (
	"context"
	"errors"
	"fmt"
	"mini-market/src/core/application/response"
	"mini-market/src/core/domain/ports/security"
	"mini-market/src/core/utils"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type SessionDenylistImpl struct {
	client *goredis.Client
}

// @inject
func NewSessionDenylistImpl(client *goredis.Client) security.SessionDenylist {
	return &SessionDenylistImpl{client: client}
}

func denylistKey(sessionID string) string {
	return "auth:revoked_sid:" + sessionID
}

func (this *SessionDenylistImpl) Revoke(ctx context.Context, sessionID string, ttl time.Duration) error {
	if err := this.client.Set(ctx, denylistKey(sessionID), "1", ttl).Err(); err != nil {
		return response.NewSafeError(
			response.CodeGatewayError,
			fmt.Errorf("[SessionDenylist] SET %q: %w", sessionID, err),
			utils.CallerPath(1),
		)
	}
	return nil
}

func (this *SessionDenylistImpl) IsRevoked(ctx context.Context, sessionID string) (bool, error) {
	err := this.client.Get(ctx, denylistKey(sessionID)).Err()

	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, goredis.Nil):
		return false, nil
	default:
		return false, response.NewSafeError(
			response.CodeGatewayError,
			fmt.Errorf("[SessionDenylist] GET %q: %w", sessionID, err),
			utils.CallerPath(1),
		)
	}
}
