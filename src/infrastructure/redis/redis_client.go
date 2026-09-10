package redis

import (
	"mini-market/src/infrastructure/env"

	goredis "github.com/redis/go-redis/v9"
)

// @inject
func NewRedisClient(cfg *env.Env) *goredis.Client {
	return goredis.NewClient(&goredis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       0,
	})
}
