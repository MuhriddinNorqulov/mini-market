package async

import (
	"github.com/hibiken/asynq"
)

// @inject
func NewAsynqConfig() asynq.Config {
	return asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 5,
			"default":  4,
			"low":      1,
		},
	}
}

// @inject
func NewAsynqServer(
	config asynq.Config,
	redisClient asynq.RedisClientOpt,
) *asynq.Server {
	return asynq.NewServer(redisClient, config)
}

// @inject
func NewAsynqServeMux() *asynq.ServeMux {
	return asynq.NewServeMux()
}

// @inject
func NewAsynqInspector(redisClient asynq.RedisClientOpt) *asynq.Inspector {
	return asynq.NewInspector(redisClient)
}
