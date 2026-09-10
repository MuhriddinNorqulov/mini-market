package async

import (
	"context"
)

type AsyncTaskHandler interface {
	ProcessTask(ctx context.Context, task *Task) error
}

type AsyncTaskHandlerFunc func(ctx context.Context, task *Task) error

func (f AsyncTaskHandlerFunc) ProcessTask(ctx context.Context, task *Task) error {
	return f(ctx, task)
}
