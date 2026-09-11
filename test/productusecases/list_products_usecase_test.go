package productusecases_test

import (
	"context"
	"testing"

	"mini-market/src/core/application/usecases/productusecases"
	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/ports/repository"
	"mini-market/src/core/domain/ports/unitofwork"
	"mini-market/src/core/domain/query"
)

type fakeProductRepository struct {
	listCalls int
}

func (this *fakeProductRepository) Create(ctx context.Context, product *entity.ProductEntity) (*entity.ProductEntity, error) {
	return product, nil
}

func (this *fakeProductRepository) GetByID(ctx context.Context, id uint) (*entity.ProductEntity, error) {
	panic("not implemented")
}

func (this *fakeProductRepository) GetByIDs(ctx context.Context, ids []uint) ([]*entity.ProductEntity, error) {
	panic("not implemented")
}

func (this *fakeProductRepository) List(ctx context.Context, filter query.ProductListFilter) (*entity.PagingEntity[*entity.ProductEntity], error) {
	this.listCalls++
	items := []*entity.ProductEntity{{ID: uint(this.listCalls), Name: "p"}}
	return entity.NewPagingEntity(filter.Page, filter.Limit, 1, items), nil
}

func (this *fakeProductRepository) ReserveStockBatchTx(ctx context.Context, tx unitofwork.Tx, reservations []repository.StockReservation) (map[uint]int64, []uint, error) {
	panic("not implemented")
}

func (this *fakeProductRepository) ReleaseStockBatchTx(ctx context.Context, tx unitofwork.Tx, reservations []repository.StockReservation) error {
	panic("not implemented")
}

type fakeProductListCache struct {
	cached  *entity.PagingEntity[*entity.ProductEntity]
	getHits int
}

func (this *fakeProductListCache) Get(ctx context.Context) (*entity.PagingEntity[*entity.ProductEntity], bool) {
	this.getHits++
	if this.cached == nil {
		return nil, false
	}
	return this.cached, true
}

func (this *fakeProductListCache) Set(ctx context.Context, page *entity.PagingEntity[*entity.ProductEntity]) {
	this.cached = page
}

func (this *fakeProductListCache) InvalidateAll(ctx context.Context) {
	this.cached = nil
}

func TestListProductsUseCase_DefaultPageIsCached(t *testing.T) {
	repo := &fakeProductRepository{}
	listCache := &fakeProductListCache{}
	uc := productusecases.NewListProductsUseCase(repo, listCache)

	first, err := uc.Invoke(context.Background(), query.ProductListFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.listCalls != 1 {
		t.Fatalf("expected repository List called once, got %d", repo.listCalls)
	}

	second, err := uc.Invoke(context.Background(), query.ProductListFilter{Page: 1, Limit: query.DefaultPageLimit})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.listCalls != 1 {
		t.Fatalf("expected repository List still called once (second call served from cache), got %d", repo.listCalls)
	}
	if first.Items[0].ID != second.Items[0].ID {
		t.Fatalf("expected cached result identical to first, got %d vs %d", first.Items[0].ID, second.Items[0].ID)
	}
}

func TestListProductsUseCase_NonDefaultPageBypassesCache(t *testing.T) {
	repo := &fakeProductRepository{}
	listCache := &fakeProductListCache{}
	uc := productusecases.NewListProductsUseCase(repo, listCache)

	if _, err := uc.Invoke(context.Background(), query.ProductListFilter{Page: 2, Limit: query.DefaultPageLimit}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if listCache.getHits != 0 {
		t.Fatalf("expected cache never consulted for a non-default page, got %d hits", listCache.getHits)
	}
	if repo.listCalls != 1 {
		t.Fatalf("expected repository List called, got %d", repo.listCalls)
	}

	if _, err := uc.Invoke(context.Background(), query.ProductListFilter{Page: 1, Limit: 5}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if listCache.getHits != 0 {
		t.Fatalf("expected cache never consulted for a non-default limit, got %d hits", listCache.getHits)
	}
	if repo.listCalls != 2 {
		t.Fatalf("expected repository List called again, got %d", repo.listCalls)
	}
}
