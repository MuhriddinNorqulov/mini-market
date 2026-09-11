package async

import "context"

type AsyncContext interface {
	GetTaskID(c context.Context) string
	GetRetryCount(c context.Context) int
	GetMaxRetry(c context.Context) int
}
