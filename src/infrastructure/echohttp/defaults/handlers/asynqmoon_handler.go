package handlers

import (
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
)

// @inject
func NewAsynqmonHandler(redisClient asynq.RedisClientOpt) *asynqmon.HTTPHandler {
	asynqmonHandler := asynqmon.New(asynqmon.Options{
		RootPath:     "/dev/monitoring/tasks",
		RedisConnOpt: redisClient,
	})
	return asynqmonHandler
}
