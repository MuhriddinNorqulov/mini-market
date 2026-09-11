package orderusecases

import (
	"context"

	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/ports/cache"
	"mini-market/src/core/domain/ports/repository"
	"mini-market/src/core/domain/query"
)

type ListMyOrdersUseCase struct {
	orderRepo repository.OrderRepository
	listCache cache.UserOrderListCache
}

// @inject
func NewListMyOrdersUseCase(orderRepo repository.OrderRepository, listCache cache.UserOrderListCache) *ListMyOrdersUseCase {
	return &ListMyOrdersUseCase{orderRepo: orderRepo, listCache: listCache}
}

func (this *ListMyOrdersUseCase) Invoke(ctx context.Context, caller *entity.UserEntity, filter query.OrderListFilter) (*entity.PagingEntity[*entity.OrderEntity], error) {
	filter.Normalize()

	cacheable := filter.Page == 1 && filter.Limit == query.DefaultPageLimit
	if cacheable {
		if cached, ok := this.listCache.Get(ctx, caller.ID); ok {
			return cached, nil
		}
	}

	result, err := this.orderRepo.ListByUserID(ctx, caller.ID, filter)
	if err != nil {
		return nil, err
	}

	if cacheable {
		this.listCache.Set(ctx, caller.ID, result)
	}
	return result, nil
}
