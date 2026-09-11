package cache

import (
	"context"

	"mini-market/src/core/domain/entity"
)

type OrderListCache interface {
	Get(ctx context.Context) (*entity.PagingEntity[*entity.OrderEntity], bool)
	Set(ctx context.Context, page *entity.PagingEntity[*entity.OrderEntity])
	InvalidateAll(ctx context.Context)
}
