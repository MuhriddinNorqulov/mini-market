package productusecases

import (
	"context"

	"mini-market/src/core/application/dto"
	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/ports/cache"
	"mini-market/src/core/domain/ports/repository"
)

type CreateProductUseCase struct {
	productRepo repository.ProductRepository
	listCache   cache.ProductListCache
}

// @inject
func NewCreateProductUseCase(productRepo repository.ProductRepository, listCache cache.ProductListCache) *CreateProductUseCase {
	return &CreateProductUseCase{productRepo: productRepo, listCache: listCache}
}

func (this *CreateProductUseCase) Invoke(ctx context.Context, input dto.CreateProductInput) (*entity.ProductEntity, error) {
	product, err := this.productRepo.Create(ctx, &entity.ProductEntity{Name: input.Name, Price: input.Price, StockQuantity: input.StockQuantity})
	if err != nil {
		return nil, err
	}
	this.listCache.InvalidateAll(ctx)
	return product, nil
}
