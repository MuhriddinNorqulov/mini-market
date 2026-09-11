package orderusecases_test

import (
	"context"
	"testing"

	"mini-market/src/core/application/usecases/orderusecases"
	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/query"
)

type fakeUserOrderListCache struct {
	cached  map[uint]*entity.PagingEntity[*entity.OrderEntity]
	getHits int
}

func newFakeUserOrderListCache() *fakeUserOrderListCache {
	return &fakeUserOrderListCache{cached: map[uint]*entity.PagingEntity[*entity.OrderEntity]{}}
}

func (this *fakeUserOrderListCache) Get(ctx context.Context, userID uint) (*entity.PagingEntity[*entity.OrderEntity], bool) {
	this.getHits++
	page, ok := this.cached[userID]
	return page, ok
}

func (this *fakeUserOrderListCache) Set(ctx context.Context, userID uint, page *entity.PagingEntity[*entity.OrderEntity]) {
	this.cached[userID] = page
}

func (this *fakeUserOrderListCache) InvalidateAll(ctx context.Context, userID uint) {
	delete(this.cached, userID)
}

func TestListMyOrdersUseCase_DefaultPageIsCachedPerCaller(t *testing.T) {
	repo := newFakeOrderRepository()
	listCache := newFakeUserOrderListCache()
	uc := orderusecases.NewListMyOrdersUseCase(repo, listCache)
	caller := &entity.UserEntity{ID: 1, Role: enum.RoleUser}

	first, err := uc.Invoke(context.Background(), caller, query.OrderListFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.listCalls != 1 {
		t.Fatalf("expected repository called once, got %d", repo.listCalls)
	}

	second, err := uc.Invoke(context.Background(), caller, query.OrderListFilter{Page: 1, Limit: query.DefaultPageLimit})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.listCalls != 1 {
		t.Fatalf("expected repository still called once (second call served from cache), got %d", repo.listCalls)
	}
	if first.Items[0].ID != second.Items[0].ID {
		t.Fatalf("expected cached result identical to first, got %d vs %d", first.Items[0].ID, second.Items[0].ID)
	}
}

func TestListMyOrdersUseCase_NonDefaultPageBypassesCache(t *testing.T) {
	repo := newFakeOrderRepository()
	listCache := newFakeUserOrderListCache()
	uc := orderusecases.NewListMyOrdersUseCase(repo, listCache)
	caller := &entity.UserEntity{ID: 1, Role: enum.RoleUser}

	if _, err := uc.Invoke(context.Background(), caller, query.OrderListFilter{Page: 2, Limit: query.DefaultPageLimit}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if listCache.getHits != 0 {
		t.Fatalf("expected cache never consulted for a non-default page, got %d hits", listCache.getHits)
	}
	if repo.listCalls != 1 {
		t.Fatalf("expected repository called, got %d", repo.listCalls)
	}
}

func TestListMyOrdersUseCase_DifferentCallersHaveIndependentCaches(t *testing.T) {
	repo := newFakeOrderRepository()
	listCache := newFakeUserOrderListCache()
	uc := orderusecases.NewListMyOrdersUseCase(repo, listCache)

	callerA := &entity.UserEntity{ID: 1, Role: enum.RoleUser}
	callerB := &entity.UserEntity{ID: 2, Role: enum.RoleUser}

	if _, err := uc.Invoke(context.Background(), callerA, query.OrderListFilter{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := uc.Invoke(context.Background(), callerB, query.OrderListFilter{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.listCalls != 2 {
		t.Fatalf("expected repository called once per caller (2 total), got %d", repo.listCalls)
	}
}
