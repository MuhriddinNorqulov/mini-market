package orderusecases

import (
	"context"
	"sort"

	"mini-market/src/core/application/response"
	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/cache"
	"mini-market/src/core/domain/ports/repository"
	"mini-market/src/core/domain/ports/unitofwork"
)

type CancelOrderUseCase struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	atomic      unitofwork.Atomic
	cache       cache.ProductCache
}

// @inject
func NewCancelOrderUseCase(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	atomic unitofwork.Atomic,
	cache cache.ProductCache,
) *CancelOrderUseCase {
	return &CancelOrderUseCase{orderRepo: orderRepo, productRepo: productRepo, atomic: atomic, cache: cache}
}

func (this *CancelOrderUseCase) Invoke(ctx context.Context, caller *entity.UserEntity, id uint, reason enum.OrderCancelReason) error {
	var order *entity.OrderEntity

	err := this.atomic.Transaction(func(tx unitofwork.Tx) error {
		o, err := this.orderRepo.GetForUpdateTx(ctx, tx, id)
		if err != nil {
			return err
		}
		if o.UserID != caller.ID && caller.Role != enum.RoleAdmin {
			return response.NewResponse(response.CodeForbidden, false, nil, "not your order")
		}
		if o.Status != enum.OrderStatusPending {
			return response.NewResponse(response.CodeInvalidOrderStatus, false, nil, "order is not pending")
		}

		reservations := make([]repository.StockReservation, 0, len(o.Items))
		for _, item := range o.Items {
			reservations = append(reservations, repository.StockReservation{ProductID: item.ProductID, Quantity: item.Quantity})
		}
		sort.Slice(reservations, func(i, j int) bool { return reservations[i].ProductID < reservations[j].ProductID })

		if err := this.productRepo.ReleaseStockBatchTx(ctx, tx, reservations); err != nil {
			return err
		}

		order = o
		return this.orderRepo.UpdateStatusTx(ctx, tx, id, enum.OrderStatusCancelled, &reason)
	})
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		this.cache.Invalidate(ctx, item.ProductID)
	}
	return nil
}
