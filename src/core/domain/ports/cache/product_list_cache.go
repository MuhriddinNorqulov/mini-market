package cache

import (
	"context"

	"mini-market/src/core/domain/entity"
)

type ProductListCache interface {
	Get(ctx context.Context) (*entity.PagingEntity[*entity.ProductEntity], bool)
	Set(ctx context.Context, page *entity.PagingEntity[*entity.ProductEntity])
	InvalidateAll(ctx context.Context)
}
