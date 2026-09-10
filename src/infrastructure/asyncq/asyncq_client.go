package async

import (
	"github.com/hibiken/asynq"
)

// @inject
func NewtAsynqClient(redis asynq.RedisClientOpt) *asynq.Client {
	return asynq.NewClient(redis)
}
