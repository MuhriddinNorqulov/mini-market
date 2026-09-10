package mapper

import (
	"mini-market/src/core/domain/ports/async"

	"github.com/hibiken/asynq"
)

func AsynqTaskMiddlewareMapper(middleware async.AsyncTaskMiddleware) asynq.MiddlewareFunc {
	return func(next asynq.Handler) asynq.Handler {
		return ToAsynQTask(middleware(ToAsynCTask(next)))
	}
}
