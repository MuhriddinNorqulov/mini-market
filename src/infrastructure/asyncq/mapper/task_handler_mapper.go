package mapper

import (
	"context"
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/async"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type _asyncHandler struct {
	handler async.AsyncTaskHandler
}

func (this *_asyncHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	t := async.NewTask(enum.TaskType(task.Type()), task.Payload())
	taskID, ok := asynq.GetTaskID(ctx)
	if !ok {
		taskID = uuid.NewString()
	}
	t.TaskID = taskID
	return this.handler.ProcessTask(ctx, t)
}

func ToAsynQTask(handler async.AsyncTaskHandler) asynq.Handler {
	return &_asyncHandler{handler: handler}
}

type _asynqHandler struct {
	handler asynq.Handler
}

func (this *_asynqHandler) ProcessTask(ctx context.Context, task *async.Task) error {
	t := asynq.NewTask(string(task.TaskType), task.Payload)
	return this.handler.ProcessTask(ctx, t)
}

func ToAsynCTask(handler asynq.Handler) async.AsyncTaskHandler {
	return &_asynqHandler{handler: handler}
}
