package tasks

import (
	"context"
	"strconv"
	"time"

	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/async"
)

type OrderAutoCancelTask struct {
	pub async.TaskPublisher
}

// @inject
func NewOrderAutoCancelTask(pub async.TaskPublisher) *OrderAutoCancelTask {
	return &OrderAutoCancelTask{pub: pub}
}

func (this *OrderAutoCancelTask) Schedule(ctx context.Context, orderID uint, processAt time.Time) error {
	payload := []byte(strconv.FormatUint(uint64(orderID), 10))
	task := async.NewTask(enum.TaskTypeOrderAutoCancel, payload)
	return this.pub.Publish(ctx, task, async.Option{ProcessAt: &processAt})
}
