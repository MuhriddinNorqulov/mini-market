package handlers

import (
	"context"

	"mini-market/src/core/application/response"
	"mini-market/src/core/application/usecases/orderusecases"
	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/async"
	"mini-market/src/core/utils"
)

type OrderAutoCancelHandler struct {
	uc *orderusecases.CancelOrderUseCase
}

// @inject
func NewOrderAutoCancelHandler(uc *orderusecases.CancelOrderUseCase) *OrderAutoCancelHandler {
	return &OrderAutoCancelHandler{uc: uc}
}

func (this *OrderAutoCancelHandler) ProcessTask(ctx context.Context, task *async.Task) error {
	orderID, err := utils.JsonUnmarshal[uint](task.Payload)
	if err != nil {
		return err
	}

	systemCaller := &entity.UserEntity{Role: enum.RoleAdmin}
	err = this.uc.Invoke(ctx, systemCaller, orderID, enum.OrderCancelReasonTimeout)
	if response.IsErrorCode(err, response.CodeInvalidOrderStatus) {
		return nil
	}
	return err
}
