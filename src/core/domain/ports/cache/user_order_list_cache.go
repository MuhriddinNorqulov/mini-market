package cache

import (
	"context"

	"mini-market/src/core/domain/entity"
)

type UserOrderListCache interface {
	Get(ctx context.Context, userID uint) (*entity.PagingEntity[*entity.OrderEntity], bool)
	Set(ctx context.Context, userID uint, page *entity.PagingEntity[*entity.OrderEntity])
	InvalidateAll(ctx context.Context, userID uint)
}
