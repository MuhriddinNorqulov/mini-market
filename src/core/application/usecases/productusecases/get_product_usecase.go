package productusecases

import (
	"context"

	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/ports/cache"
	"mini-market/src/core/domain/ports/repository"
)

type GetProductUseCase struct {
	productRepo repository.ProductRepository
	cache       cache.ProductCache
}

// @inject
func NewGetProductUseCase(productRepo repository.ProductRepository, cache cache.ProductCache) *GetProductUseCase {
	return &GetProductUseCase{productRepo: productRepo, cache: cache}
}

func (this *GetProductUseCase) Invoke(ctx context.Context, id uint) (*entity.ProductEntity, error) {
	if cached, ok := this.cache.Get(ctx, id); ok {
		return cached, nil
	}

	product, err := this.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	this.cache.Set(ctx, product)
	return product, nil
}
