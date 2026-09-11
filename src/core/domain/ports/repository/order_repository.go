package repository

import (
	"context"

	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/unitofwork"
	"mini-market/src/core/domain/query"
)

type OrderRepository interface {
	CreateTx(ctx context.Context, tx unitofwork.Tx, order *entity.OrderEntity) (*entity.OrderEntity, error)
	GetByIdempotencyKey(ctx context.Context, userID uint, idempotencyKey string) (*entity.OrderEntity, error)
	GetByIdempotencyKeyTx(ctx context.Context, tx unitofwork.Tx, userID uint, idempotencyKey string) (*entity.OrderEntity, error)
	GetByID(ctx context.Context, id uint) (*entity.OrderEntity, error)
	GetForUpdateTx(ctx context.Context, tx unitofwork.Tx, id uint) (*entity.OrderEntity, error)
	UpdateStatusTx(ctx context.Context, tx unitofwork.Tx, id uint, status enum.OrderStatus, cancelReason *enum.OrderCancelReason) error
	List(ctx context.Context, filter query.OrderListFilter) (*entity.PagingEntity[*entity.OrderEntity], error)
	ListByUserID(ctx context.Context, userID uint, filter query.OrderListFilter) (*entity.PagingEntity[*entity.OrderEntity], error)
}
