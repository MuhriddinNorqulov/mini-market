package orderusecases

import (
	"context"

	"mini-market/src/core/application/response"
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/repository"
	"mini-market/src/core/domain/ports/unitofwork"
)

type ConfirmOrderUseCase struct {
	orderRepo repository.OrderRepository
	atomic    unitofwork.Atomic
}

// @inject
func NewConfirmOrderUseCase(orderRepo repository.OrderRepository, atomic unitofwork.Atomic) *ConfirmOrderUseCase {
	return &ConfirmOrderUseCase{orderRepo: orderRepo, atomic: atomic}
}

func (this *ConfirmOrderUseCase) Invoke(ctx context.Context, id uint) error {
	return this.atomic.Transaction(func(tx unitofwork.Tx) error {
		order, err := this.orderRepo.GetForUpdateTx(ctx, tx, id)
		if err != nil {
			return err
		}
		if order.Status != enum.OrderStatusPending {
			return response.NewResponse(response.CodeInvalidOrderStatus, false, nil, "order is not pending")
		}
		return this.orderRepo.UpdateStatusTx(ctx, tx, id, enum.OrderStatusConfirmed, nil)
	})
}
