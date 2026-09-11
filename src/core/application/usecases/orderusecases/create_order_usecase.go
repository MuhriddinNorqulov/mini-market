package orderusecases

import (
	"context"
	"fmt"
	"sort"
	"time"

	"mini-market/src/core/application/dto"
	"mini-market/src/core/application/response"
	"mini-market/src/core/application/tasks"
	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/cache"
	"mini-market/src/core/domain/ports/repository"
	"mini-market/src/core/domain/ports/unitofwork"
)

const orderExpiryDuration = 1 * time.Minute

type CreateOrderUseCase struct {
	orderRepo       repository.OrderRepository
	productRepo     repository.ProductRepository
	atomic          unitofwork.Atomic
	orderAutoCancel *tasks.OrderAutoCancelTask
	cache           cache.ProductCache
	listCache       cache.OrderListCache
	myOrdersCache   cache.UserOrderListCache
}

// @inject
func NewCreateOrderUseCase(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	atomic unitofwork.Atomic,
	orderAutoCancel *tasks.OrderAutoCancelTask,
	cache cache.ProductCache,
	listCache cache.OrderListCache,
	myOrdersCache cache.UserOrderListCache,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepo:       orderRepo,
		productRepo:     productRepo,
		atomic:          atomic,
		orderAutoCancel: orderAutoCancel,
		cache:           cache,
		listCache:       listCache,
		myOrdersCache:   myOrdersCache,
	}
}

func (this *CreateOrderUseCase) Invoke(ctx context.Context, userID uint, idempotencyKey string, items []dto.CreateOrderItemInput) (*entity.OrderEntity, error) {
	if len(items) == 0 {
		return nil, response.NewResponse(response.CodeBadRequest, false, nil, "order must contain at least one item")
	}

	if existing, err := this.orderRepo.GetByIdempotencyKey(ctx, userID, idempotencyKey); err == nil {
		return existing, nil
	}

	reservations := aggregateAndSortReservations(items)

	if err := this.precheckAvailability(ctx, reservations); err != nil {
		return nil, err
	}

	var order *entity.OrderEntity
	err := this.atomic.Transaction(func(tx unitofwork.Tx) error {
		if existing, err := this.orderRepo.GetByIdempotencyKeyTx(ctx, tx, userID, idempotencyKey); err == nil {
			order = existing
			return nil
		}

		prices, failedIDs, err := this.productRepo.ReserveStockBatchTx(ctx, tx, reservations)
		if err != nil {
			return err
		}
		if len(failedIDs) > 0 {
			return insufficientStockError(failedIDs)
		}

		var totalPrice int64
		orderItems := make([]*entity.OrderItemEntity, 0, len(reservations))
		for _, r := range reservations {
			price := prices[r.ProductID]
			orderItems = append(orderItems, &entity.OrderItemEntity{
				ProductID: r.ProductID,
				Quantity:  r.Quantity,
				UnitPrice: price,
			})
			totalPrice += price * r.Quantity
		}

		created, err := this.orderRepo.CreateTx(ctx, tx, &entity.OrderEntity{
			UserID:         userID,
			IdempotencyKey: idempotencyKey,
			TotalPrice:     totalPrice,
			Items:          orderItems,
			Status:         enum.OrderStatusPending,
			ExpiresAt:      time.Now().Add(orderExpiryDuration),
		})
		if err != nil {
			return err
		}

		order = created
		return nil
	})
	if err != nil {
		return nil, err
	}

	for _, item := range order.Items {
		this.cache.Invalidate(ctx, item.ProductID)
	}
	this.listCache.InvalidateAll(ctx)
	this.myOrdersCache.InvalidateAll(ctx, userID)

	_ = this.orderAutoCancel.Schedule(ctx, order.ID, order.ExpiresAt)

	return order, nil
}

func (this *CreateOrderUseCase) precheckAvailability(ctx context.Context, reservations []repository.StockReservation) error {
	ids := make([]uint, len(reservations))
	for i, r := range reservations {
		ids[i] = r.ProductID
	}

	products, err := this.productRepo.GetByIDs(ctx, ids)
	if err != nil {
		return err
	}

	productByID := make(map[uint]*entity.ProductEntity, len(products))
	for _, p := range products {
		productByID[p.ID] = p
	}

	var missingIDs, insufficientIDs []uint
	for _, r := range reservations {
		product, ok := productByID[r.ProductID]
		if !ok {
			missingIDs = append(missingIDs, r.ProductID)
			continue
		}
		if product.StockQuantity < r.Quantity {
			insufficientIDs = append(insufficientIDs, r.ProductID)
		}
	}

	if len(missingIDs) > 0 {
		return response.NewResponse(response.CodNotFound, false, missingIDs, fmt.Sprintf("products not found: %v", missingIDs))
	}
	if len(insufficientIDs) > 0 {
		return insufficientStockError(insufficientIDs)
	}
	return nil
}

func aggregateAndSortReservations(items []dto.CreateOrderItemInput) []repository.StockReservation {
	quantityByProduct := make(map[uint]int64, len(items))
	for _, item := range items {
		quantityByProduct[item.ProductID] += item.Quantity
	}

	reservations := make([]repository.StockReservation, 0, len(quantityByProduct))
	for productID, quantity := range quantityByProduct {
		reservations = append(reservations, repository.StockReservation{ProductID: productID, Quantity: quantity})
	}

	sort.Slice(reservations, func(i, j int) bool { return reservations[i].ProductID < reservations[j].ProductID })
	return reservations
}

func insufficientStockError(productIDs []uint) error {
	return response.NewResponse(response.CodeInsufficientStock, false, productIDs, fmt.Sprintf("insufficient stock for products: %v", productIDs))
}
