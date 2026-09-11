package async

import (
	"context"
	"time"
)

type Option struct {
	ProcessIn     *time.Duration
	ProcessAt     *time.Time
	MaxRetryCount int
	TimeOut       *time.Duration
}

type TaskPublisher interface {
	Publish(ctx context.Context, task *Task, option Option) error
}
