package redis

import (
	"context"
	"fmt"
	"time"

	"mini-market/src/core/application/response"
	"mini-market/src/core/domain/ports/ratelimit"
	"mini-market/src/core/utils"

	goredis "github.com/redis/go-redis/v9"
)

type RateLimiterImpl struct {
	client *goredis.Client
}

// @inject
func NewRateLimiterImpl(client *goredis.Client) ratelimit.Limiter {
	return &RateLimiterImpl{client: client}
}

func (this *RateLimiterImpl) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	count, err := this.client.Incr(ctx, key).Result()
	if err != nil {
		return false, response.NewSafeError(
			response.CodeGatewayError,
			fmt.Errorf("[RateLimiter] INCR %q: %w", key, err),
			utils.CallerPath(1),
		)
	}

	if count == 1 {
		if err := this.client.Expire(ctx, key, window).Err(); err != nil {
			return false, response.NewSafeError(
				response.CodeGatewayError,
				fmt.Errorf("[RateLimiter] EXPIRE %q: %w", key, err),
				utils.CallerPath(1),
			)
		}
	}

	return count <= int64(limit), nil
}
