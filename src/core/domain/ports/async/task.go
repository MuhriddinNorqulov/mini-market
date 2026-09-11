package async

import "mini-market/src/core/domain/entity/enum"

type Task struct {
	TaskID   string
	TaskType enum.TaskType
	Payload  []byte
}

func NewTask(taskType enum.TaskType, payload []byte) *Task {
	return &Task{TaskType: taskType, Payload: payload}
}
