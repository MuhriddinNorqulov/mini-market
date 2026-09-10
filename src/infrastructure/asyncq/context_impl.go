package async

import (
	"context"
	"mini-market/src/core/domain/ports/async"

	"github.com/hibiken/asynq"
)

type AsyncContextImpl struct {
}

// @inject
func NewAsyncContextImpl() async.AsyncContext {
	return &AsyncContextImpl{}
}

func (this *AsyncContextImpl) GetTaskID(c context.Context) string {
	id, ok := asynq.GetTaskID(c)
	if !ok {
		return ""
	}
	return id
}

func (this *AsyncContextImpl) GetRetryCount(c context.Context) int {
	count, ok := asynq.GetRetryCount(c)
	if !ok {
		return 0
	}
	return count
}

func (this *AsyncContextImpl) GetMaxRetry(c context.Context) int {
	count, ok := asynq.GetMaxRetry(c)
	if !ok {
		return 0
	}
	return count
}
