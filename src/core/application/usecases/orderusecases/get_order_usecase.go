package orderusecases

import (
	"context"

	"mini-market/src/core/application/response"
	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/repository"
)

type GetOrderUseCase struct {
	orderRepo repository.OrderRepository
}

// @inject
func NewGetOrderUseCase(orderRepo repository.OrderRepository) *GetOrderUseCase {
	return &GetOrderUseCase{orderRepo: orderRepo}
}

func (this *GetOrderUseCase) Invoke(ctx context.Context, caller *entity.UserEntity, id uint) (*entity.OrderEntity, error) {
	order, err := this.orderRepo.GetByID(ctx, id)
	if err != nil {
		if response.IsErrorCode(err, response.CodNotFound) {
			return nil, response.NewResponse(response.CodNotFound, false, nil, "order not found")
		}
		return nil, err
	}
	if order.UserID != caller.ID && caller.Role != enum.RoleAdmin {
		return nil, response.NewResponse(response.CodeForbidden, false, nil, "not your order")
	}
	return order, nil
}
