package async

import (
	"context"
	"mini-market/src/core/domain/entity/enum"
)

type AsyncServer interface {
	Init()
	Run() error
	Shutdown(ctx context.Context) error

	Use(middlewares ...AsyncTaskMiddleware)
	HandlerFunc(taskName enum.TaskType, handler AsyncTaskHandler)
	ProcessTask(ctx context.Context, task *Task) error
}
