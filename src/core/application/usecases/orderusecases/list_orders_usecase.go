package orderusecases

import (
	"context"

	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/ports/cache"
	"mini-market/src/core/domain/ports/repository"
	"mini-market/src/core/domain/query"
)

type ListOrdersUseCase struct {
	orderRepo repository.OrderRepository
	listCache cache.OrderListCache
}

// @inject
func NewListOrdersUseCase(orderRepo repository.OrderRepository, listCache cache.OrderListCache) *ListOrdersUseCase {
	return &ListOrdersUseCase{orderRepo: orderRepo, listCache: listCache}
}

func (this *ListOrdersUseCase) Invoke(ctx context.Context, filter query.OrderListFilter) (*entity.PagingEntity[*entity.OrderEntity], error) {
	filter.Normalize()

	cacheable := filter.Page == 1 && filter.Limit == query.DefaultPageLimit
	if cacheable {
		if cached, ok := this.listCache.Get(ctx); ok {
			return cached, nil
		}
	}

	result, err := this.orderRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	if cacheable {
		this.listCache.Set(ctx, result)
	}
	return result, nil
}
