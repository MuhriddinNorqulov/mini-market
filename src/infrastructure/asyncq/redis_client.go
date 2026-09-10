package async

import (
	"mini-market/src/infrastructure/env"

	"github.com/hibiken/asynq"
)

// @inject
func NewAsynqRedisClientOpt0(cfg *env.Env) asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       0,
	}
}
