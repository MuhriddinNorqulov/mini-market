package productusecases

import (
	"context"

	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/ports/cache"
	"mini-market/src/core/domain/ports/repository"
	"mini-market/src/core/domain/query"
)

type ListProductsUseCase struct {
	productRepo repository.ProductRepository
	listCache   cache.ProductListCache
}

// @inject
func NewListProductsUseCase(productRepo repository.ProductRepository, listCache cache.ProductListCache) *ListProductsUseCase {
	return &ListProductsUseCase{productRepo: productRepo, listCache: listCache}
}

func (this *ListProductsUseCase) Invoke(ctx context.Context, filter query.ProductListFilter) (*entity.PagingEntity[*entity.ProductEntity], error) {
	filter.Normalize()

	cacheable := filter.Page == 1 && filter.Limit == query.DefaultPageLimit
	if cacheable {
		if cached, ok := this.listCache.Get(ctx); ok {
			return cached, nil
		}
	}

	result, err := this.productRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	if cacheable {
		this.listCache.Set(ctx, result)
	}
	return result, nil
}
