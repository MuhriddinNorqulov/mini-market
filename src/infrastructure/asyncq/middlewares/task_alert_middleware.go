package middlewares

import (
	"context"
	"fmt"
	"mini-market/src/core/domain/ports/async"
	"mini-market/src/core/domain/ports/notification"
)

type TaskAlertMiddleware struct {
	asyncCtx      async.AsyncContext
	alertNotifier notification.AlertNotifier
}

// @inject
func NewTaskAlertMiddleware(asyncCtx async.AsyncContext, alertNotifier notification.AlertNotifier) *TaskAlertMiddleware {
	return &TaskAlertMiddleware{asyncCtx: asyncCtx, alertNotifier: alertNotifier}
}

func (this *TaskAlertMiddleware) Wrap(next async.AsyncTaskHandler) async.AsyncTaskHandler {
	return async.AsyncTaskHandlerFunc(func(ctx context.Context, t *async.Task) error {
		err := next.ProcessTask(ctx, t)
		if err != nil && this.asyncCtx.GetRetryCount(ctx) >= this.asyncCtx.GetMaxRetry(ctx) {
			msg := fmt.Sprintf(
				"<b>Task Exhausted</b>\n\nTask: %s\nTask ID: %s\nPayload: %s\nError: %s",
				t.TaskType, t.TaskID, string(t.Payload), err.Error(),
			)
			this.alertNotifier.Send(ctx, msg)
		}
		return err
	})
}
