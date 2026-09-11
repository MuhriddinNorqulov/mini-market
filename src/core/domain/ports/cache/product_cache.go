package cache

import (
	"context"

	"mini-market/src/core/domain/entity"
)

type ProductCache interface {
	Get(ctx context.Context, id uint) (*entity.ProductEntity, bool)
	Set(ctx context.Context, product *entity.ProductEntity)
	Invalidate(ctx context.Context, id uint)
}
