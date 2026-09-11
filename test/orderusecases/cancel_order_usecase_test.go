package orderusecases_test

import (
	"context"
	"testing"

	"mini-market/src/core/application/dto"
	"mini-market/src/core/application/response"
	"mini-market/src/core/application/tasks"
	"mini-market/src/core/application/usecases/orderusecases"
	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/async"
	"mini-market/src/core/domain/ports/cache"
	"mini-market/src/core/domain/ports/repository"
	"mini-market/src/core/domain/ports/unitofwork"
	"mini-market/src/core/domain/query"
)

type fakeOrderRepository struct {
	orders        map[uint]*entity.OrderEntity
	byIdempotency map[string]*entity.OrderEntity
	nextID        uint
	listCalls     int
}

func newFakeOrderRepository() *fakeOrderRepository {
	return &fakeOrderRepository{orders: map[uint]*entity.OrderEntity{}, byIdempotency: map[string]*entity.OrderEntity{}}
}

func (this *fakeOrderRepository) CreateTx(ctx context.Context, tx unitofwork.Tx, order *entity.OrderEntity) (*entity.OrderEntity, error) {
	this.nextID++
	order.ID = this.nextID
	this.orders[order.ID] = order
	this.byIdempotency[key(order.UserID, order.IdempotencyKey)] = order
	return order, nil
}

func (this *fakeOrderRepository) GetByIdempotencyKey(ctx context.Context, userID uint, idempotencyKey string) (*entity.OrderEntity, error) {
	if order, ok := this.byIdempotency[key(userID, idempotencyKey)]; ok {
		return order, nil
	}
	return nil, response.NewResponse(response.CodNotFound, false, nil, "not found")
}

func (this *fakeOrderRepository) GetByIdempotencyKeyTx(ctx context.Context, tx unitofwork.Tx, userID uint, idempotencyKey string) (*entity.OrderEntity, error) {
	return this.GetByIdempotencyKey(ctx, userID, idempotencyKey)
}

func (this *fakeOrderRepository) GetByID(ctx context.Context, id uint) (*entity.OrderEntity, error) {
	if order, ok := this.orders[id]; ok {
		return order, nil
	}
	return nil, response.NewResponse(response.CodNotFound, false, nil, "not found")
}

func (this *fakeOrderRepository) GetForUpdateTx(ctx context.Context, tx unitofwork.Tx, id uint) (*entity.OrderEntity, error) {
	return this.GetByID(ctx, id)
}

func (this *fakeOrderRepository) UpdateStatusTx(ctx context.Context, tx unitofwork.Tx, id uint, status enum.OrderStatus, cancelReason *enum.OrderCancelReason) error {
	order, ok := this.orders[id]
	if !ok {
		return response.NewResponse(response.CodNotFound, false, nil, "not found")
	}
	order.Status = status
	order.CancelReason = cancelReason
	return nil
}

func (this *fakeOrderRepository) List(ctx context.Context, filter query.OrderListFilter) (*entity.PagingEntity[*entity.OrderEntity], error) {
	this.listCalls++
	items := []*entity.OrderEntity{{ID: uint(this.listCalls)}}
	return entity.NewPagingEntity(filter.Page, filter.Limit, 1, items), nil
}

func (this *fakeOrderRepository) ListByUserID(ctx context.Context, userID uint, filter query.OrderListFilter) (*entity.PagingEntity[*entity.OrderEntity], error) {
	this.listCalls++
	items := []*entity.OrderEntity{{ID: uint(this.listCalls), UserID: userID}}
	return entity.NewPagingEntity(filter.Page, filter.Limit, 1, items), nil
}

func key(userID uint, idempotencyKey string) string {
	return string(rune(userID)) + "|" + idempotencyKey
}

type fakeProductRepository struct {
	stock map[uint]int64
	price map[uint]int64
}

func newFakeProductRepository() *fakeProductRepository {
	return &fakeProductRepository{stock: map[uint]int64{}, price: map[uint]int64{}}
}

func (this *fakeProductRepository) Create(ctx context.Context, product *entity.ProductEntity) (*entity.ProductEntity, error) {
	return product, nil
}

func (this *fakeProductRepository) GetByID(ctx context.Context, id uint) (*entity.ProductEntity, error) {
	price, ok := this.price[id]
	if !ok {
		return nil, response.NewResponse(response.CodNotFound, false, nil, "not found")
	}
	return &entity.ProductEntity{ID: id, Price: price, StockQuantity: this.stock[id]}, nil
}

func (this *fakeProductRepository) GetByIDs(ctx context.Context, ids []uint) ([]*entity.ProductEntity, error) {
	result := make([]*entity.ProductEntity, 0, len(ids))
	for _, id := range ids {
		if price, ok := this.price[id]; ok {
			result = append(result, &entity.ProductEntity{ID: id, Price: price, StockQuantity: this.stock[id]})
		}
	}
	return result, nil
}

func (this *fakeProductRepository) List(ctx context.Context, filter query.ProductListFilter) (*entity.PagingEntity[*entity.ProductEntity], error) {
	return entity.NewPagingEntity[*entity.ProductEntity](filter.Page, filter.Limit, 0, nil), nil
}

func (this *fakeProductRepository) ReserveStockBatchTx(ctx context.Context, tx unitofwork.Tx, reservations []repository.StockReservation) (map[uint]int64, []uint, error) {
	var failed []uint
	for _, r := range reservations {
		if this.stock[r.ProductID] < r.Quantity {
			failed = append(failed, r.ProductID)
		}
	}
	if len(failed) > 0 {
		return nil, failed, nil
	}

	prices := make(map[uint]int64, len(reservations))
	for _, r := range reservations {
		this.stock[r.ProductID] -= r.Quantity
		prices[r.ProductID] = this.price[r.ProductID]
	}
	return prices, nil, nil
}

func (this *fakeProductRepository) ReleaseStockBatchTx(ctx context.Context, tx unitofwork.Tx, reservations []repository.StockReservation) error {
	for _, r := range reservations {
		this.stock[r.ProductID] += r.Quantity
	}
	return nil
}

type raceProductRepository struct {
	*fakeProductRepository
	raceProductID uint
}

func (this *raceProductRepository) GetByIDs(ctx context.Context, ids []uint) ([]*entity.ProductEntity, error) {
	result, err := this.fakeProductRepository.GetByIDs(ctx, ids)
	this.fakeProductRepository.stock[this.raceProductID] = 0
	return result, err
}

type fakeAtomic struct{}

func (this *fakeAtomic) Transaction(fn func(tx unitofwork.Tx) error) error {
	return fn(struct{}{})
}

type fakeTaskPublisher struct{}

func (this *fakeTaskPublisher) Publish(ctx context.Context, task *async.Task, opt async.Option) error {
	return nil
}

func newTestOrderAutoCancelTask() *tasks.OrderAutoCancelTask {
	return tasks.NewOrderAutoCancelTask(&fakeTaskPublisher{})
}

type fakeProductCache struct {
	entries map[uint]*entity.ProductEntity
}

func newTestProductCache() cache.ProductCache {
	return &fakeProductCache{entries: map[uint]*entity.ProductEntity{}}
}

func (this *fakeProductCache) Get(ctx context.Context, id uint) (*entity.ProductEntity, bool) {
	product, ok := this.entries[id]
	return product, ok
}

func (this *fakeProductCache) Set(ctx context.Context, product *entity.ProductEntity) {
	this.entries[product.ID] = product
}

func (this *fakeProductCache) Invalidate(ctx context.Context, id uint) {
	delete(this.entries, id)
}

type fakeOrderListCache struct {
	cached  *entity.PagingEntity[*entity.OrderEntity]
	getHits int
}

func newTestOrderListCache() cache.OrderListCache {
	return &fakeOrderListCache{}
}

func (this *fakeOrderListCache) Get(ctx context.Context) (*entity.PagingEntity[*entity.OrderEntity], bool) {
	this.getHits++
	if this.cached == nil {
		return nil, false
	}
	return this.cached, true
}

func (this *fakeOrderListCache) Set(ctx context.Context, page *entity.PagingEntity[*entity.OrderEntity]) {
	this.cached = page
}

func (this *fakeOrderListCache) InvalidateAll(ctx context.Context) {
	this.cached = nil
}

func TestCreateOrderUseCase_DuplicateIdempotencyKeyReturnsExistingOrder(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	productRepo := newFakeProductRepository()
	productRepo.stock[1] = 10
	productRepo.price[1] = 100

	uc := orderusecases.NewCreateOrderUseCase(orderRepo, productRepo, &fakeAtomic{}, newTestOrderAutoCancelTask(), newTestProductCache(), newTestOrderListCache(), newFakeUserOrderListCache())
	ctx := context.Background()

	first, err := uc.Invoke(ctx, 1, "same-key", []dto.CreateOrderItemInput{{ProductID: 1, Quantity: 1}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second, err := uc.Invoke(ctx, 1, "same-key", []dto.CreateOrderItemInput{{ProductID: 1, Quantity: 1}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first.ID != second.ID {
		t.Fatalf("expected same order returned, got %d and %d", first.ID, second.ID)
	}
	if productRepo.stock[1] != 9 {
		t.Fatalf("expected stock decremented exactly once, got %d", productRepo.stock[1])
	}
}

func TestCreateOrderUseCase_InsufficientStockReturnsConflictCode(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	productRepo := newFakeProductRepository()
	productRepo.stock[1] = 0
	productRepo.price[1] = 100

	uc := orderusecases.NewCreateOrderUseCase(orderRepo, productRepo, &fakeAtomic{}, newTestOrderAutoCancelTask(), newTestProductCache(), newTestOrderListCache(), newFakeUserOrderListCache())

	_, err := uc.Invoke(context.Background(), 1, "key", []dto.CreateOrderItemInput{{ProductID: 1, Quantity: 1}})
	if !response.IsErrorCode(err, response.CodeInsufficientStock) {
		t.Fatalf("expected CodeInsufficientStock, got %v", err)
	}
}

func TestCreateOrderUseCase_MissingProductReturnsNotFoundWithIDs(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	productRepo := newFakeProductRepository()
	productRepo.stock[1] = 10
	productRepo.price[1] = 100

	uc := orderusecases.NewCreateOrderUseCase(orderRepo, productRepo, &fakeAtomic{}, newTestOrderAutoCancelTask(), newTestProductCache(), newTestOrderListCache(), newFakeUserOrderListCache())

	_, err := uc.Invoke(context.Background(), 1, "missing-key", []dto.CreateOrderItemInput{
		{ProductID: 1, Quantity: 1},
		{ProductID: 999, Quantity: 1},
	})

	resp, ok := err.(*response.Response)
	if !ok {
		t.Fatalf("expected *response.Response error, got %v", err)
	}
	if resp.Code != response.CodNotFound {
		t.Fatalf("expected CodNotFound, got %v", resp.Code)
	}
	ids, ok := resp.Payload.([]uint)
	if !ok || len(ids) != 1 || ids[0] != 999 {
		t.Fatalf("expected payload [999], got %v", resp.Payload)
	}
}

func TestCreateOrderUseCase_InsufficientStockErrorNamesFailingProductIDs(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	productRepo := newFakeProductRepository()
	productRepo.stock[1] = 10
	productRepo.price[1] = 100
	productRepo.stock[2] = 0
	productRepo.price[2] = 50

	uc := orderusecases.NewCreateOrderUseCase(orderRepo, productRepo, &fakeAtomic{}, newTestOrderAutoCancelTask(), newTestProductCache(), newTestOrderListCache(), newFakeUserOrderListCache())

	_, err := uc.Invoke(context.Background(), 1, "insuf-ids-key", []dto.CreateOrderItemInput{
		{ProductID: 1, Quantity: 1},
		{ProductID: 2, Quantity: 1},
	})

	resp, ok := err.(*response.Response)
	if !ok {
		t.Fatalf("expected *response.Response error, got %v", err)
	}
	if resp.Code != response.CodeInsufficientStock {
		t.Fatalf("expected CodeInsufficientStock, got %v", resp.Code)
	}
	ids, ok := resp.Payload.([]uint)
	if !ok || len(ids) != 1 || ids[0] != 2 {
		t.Fatalf("expected payload [2], got %v", resp.Payload)
	}
}

func TestCreateOrderUseCase_DuplicateProductIDsAreAggregated(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	productRepo := newFakeProductRepository()
	productRepo.stock[1] = 10
	productRepo.price[1] = 100

	uc := orderusecases.NewCreateOrderUseCase(orderRepo, productRepo, &fakeAtomic{}, newTestOrderAutoCancelTask(), newTestProductCache(), newTestOrderListCache(), newFakeUserOrderListCache())

	order, err := uc.Invoke(context.Background(), 1, "dup-key", []dto.CreateOrderItemInput{
		{ProductID: 1, Quantity: 3},
		{ProductID: 1, Quantity: 2},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order.Items) != 1 {
		t.Fatalf("expected duplicate product_id lines to collapse into one order item, got %d", len(order.Items))
	}
	if order.Items[0].Quantity != 5 {
		t.Fatalf("expected aggregated quantity 5, got %d", order.Items[0].Quantity)
	}
	if productRepo.stock[1] != 5 {
		t.Fatalf("expected stock decremented by combined quantity (10-5=5), got %d", productRepo.stock[1])
	}
	if order.TotalPrice != 500 {
		t.Fatalf("expected total price 500 (5 * 100), got %d", order.TotalPrice)
	}
}

func TestCreateOrderUseCase_MultiItemOrderReservesAllAndComputesTotalPrice(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	productRepo := newFakeProductRepository()
	productRepo.stock[1] = 10
	productRepo.price[1] = 100
	productRepo.stock[2] = 10
	productRepo.price[2] = 50

	uc := orderusecases.NewCreateOrderUseCase(orderRepo, productRepo, &fakeAtomic{}, newTestOrderAutoCancelTask(), newTestProductCache(), newTestOrderListCache(), newFakeUserOrderListCache())

	order, err := uc.Invoke(context.Background(), 1, "multi-key", []dto.CreateOrderItemInput{
		{ProductID: 1, Quantity: 2},
		{ProductID: 2, Quantity: 3},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if productRepo.stock[1] != 8 || productRepo.stock[2] != 7 {
		t.Fatalf("expected stock 8 and 7, got %d and %d", productRepo.stock[1], productRepo.stock[2])
	}
	if order.TotalPrice != 2*100+3*50 {
		t.Fatalf("expected total price %d, got %d", 2*100+3*50, order.TotalPrice)
	}
}

func TestCreateOrderUseCase_RaceBetweenPrecheckAndReserveStillPreventsNegativeStock(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	base := newFakeProductRepository()
	base.stock[1] = 5
	base.price[1] = 100
	productRepo := &raceProductRepository{fakeProductRepository: base, raceProductID: 1}

	uc := orderusecases.NewCreateOrderUseCase(orderRepo, productRepo, &fakeAtomic{}, newTestOrderAutoCancelTask(), newTestProductCache(), newTestOrderListCache(), newFakeUserOrderListCache())

	_, err := uc.Invoke(context.Background(), 1, "race-key", []dto.CreateOrderItemInput{{ProductID: 1, Quantity: 2}})
	if !response.IsErrorCode(err, response.CodeInsufficientStock) {
		t.Fatalf("expected CodeInsufficientStock from the in-transaction guard, got %v", err)
	}
	if base.stock[1] != 0 {
		t.Fatalf("expected stock to remain at the raced value (0), got %d", base.stock[1])
	}
}

func TestCancelOrderUseCase_ReleasesStockForAllItems(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	productRepo := newFakeProductRepository()
	productRepo.stock[1] = 5
	productRepo.stock[2] = 5

	orderRepo.orders[1] = &entity.OrderEntity{
		ID:     1,
		UserID: 1,
		Status: enum.OrderStatusPending,
		Items: []*entity.OrderItemEntity{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 3},
		},
	}

	uc := orderusecases.NewCancelOrderUseCase(orderRepo, productRepo, &fakeAtomic{}, newTestProductCache())

	if err := uc.Invoke(context.Background(), &entity.UserEntity{ID: 1, Role: enum.RoleUser}, 1, enum.OrderCancelReasonUserRequested); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if productRepo.stock[1] != 7 || productRepo.stock[2] != 8 {
		t.Fatalf("expected stock 7 and 8 after release, got %d and %d", productRepo.stock[1], productRepo.stock[2])
	}
}

func TestCancelOrderUseCase_NonPendingOrderReturnsInvalidStatus(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	productRepo := newFakeProductRepository()
	orderRepo.orders[1] = &entity.OrderEntity{ID: 1, UserID: 1, Status: enum.OrderStatusConfirmed}

	uc := orderusecases.NewCancelOrderUseCase(orderRepo, productRepo, &fakeAtomic{}, newTestProductCache())

	err := uc.Invoke(context.Background(), &entity.UserEntity{ID: 1, Role: enum.RoleUser}, 1, enum.OrderCancelReasonUserRequested)
	if !response.IsErrorCode(err, response.CodeInvalidOrderStatus) {
		t.Fatalf("expected CodeInvalidOrderStatus, got %v", err)
	}
}

func TestCancelOrderUseCase_NonOwnerNonAdminIsForbidden(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	productRepo := newFakeProductRepository()
	orderRepo.orders[1] = &entity.OrderEntity{ID: 1, UserID: 1, Status: enum.OrderStatusPending}

	uc := orderusecases.NewCancelOrderUseCase(orderRepo, productRepo, &fakeAtomic{}, newTestProductCache())

	err := uc.Invoke(context.Background(), &entity.UserEntity{ID: 2, Role: enum.RoleUser}, 1, enum.OrderCancelReasonUserRequested)
	if !response.IsErrorCode(err, response.CodeForbidden) {
		t.Fatalf("expected CodeForbidden, got %v", err)
	}
}

func TestCancelOrderUseCase_AdminCanCancelAnyOrder(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	productRepo := newFakeProductRepository()
	productRepo.stock[1] = 5
	orderRepo.orders[1] = &entity.OrderEntity{
		ID:     1,
		UserID: 1,
		Status: enum.OrderStatusPending,
		Items:  []*entity.OrderItemEntity{{ProductID: 1, Quantity: 2}},
	}

	uc := orderusecases.NewCancelOrderUseCase(orderRepo, productRepo, &fakeAtomic{}, newTestProductCache())

	err := uc.Invoke(context.Background(), &entity.UserEntity{ID: 2, Role: enum.RoleAdmin}, 1, enum.OrderCancelReasonUserRequested)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if orderRepo.orders[1].Status != enum.OrderStatusCancelled {
		t.Fatalf("expected order to be cancelled, got %s", orderRepo.orders[1].Status)
	}
}

func TestConfirmOrderUseCase_NonPendingOrderReturnsInvalidStatus(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	orderRepo.orders[1] = &entity.OrderEntity{ID: 1, Status: enum.OrderStatusCancelled}

	uc := orderusecases.NewConfirmOrderUseCase(orderRepo, &fakeAtomic{})

	err := uc.Invoke(context.Background(), 1)
	if !response.IsErrorCode(err, response.CodeInvalidOrderStatus) {
		t.Fatalf("expected CodeInvalidOrderStatus, got %v", err)
	}
}

func TestListOrdersUseCase_DefaultPageIsCached(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	listCache := &fakeOrderListCache{}
	uc := orderusecases.NewListOrdersUseCase(orderRepo, listCache)

	first, err := uc.Invoke(context.Background(), query.OrderListFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if orderRepo.listCalls != 1 {
		t.Fatalf("expected repository List called once, got %d", orderRepo.listCalls)
	}

	second, err := uc.Invoke(context.Background(), query.OrderListFilter{Page: 1, Limit: query.DefaultPageLimit})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if orderRepo.listCalls != 1 {
		t.Fatalf("expected repository List still called once (second call served from cache), got %d", orderRepo.listCalls)
	}
	if first.Items[0].ID != second.Items[0].ID {
		t.Fatalf("expected cached result identical to first, got %d vs %d", first.Items[0].ID, second.Items[0].ID)
	}
}

func TestListOrdersUseCase_NonDefaultPageBypassesCache(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	listCache := &fakeOrderListCache{}
	uc := orderusecases.NewListOrdersUseCase(orderRepo, listCache)

	if _, err := uc.Invoke(context.Background(), query.OrderListFilter{Page: 2, Limit: query.DefaultPageLimit}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if listCache.getHits != 0 {
		t.Fatalf("expected cache never consulted for a non-default page, got %d hits", listCache.getHits)
	}
	if orderRepo.listCalls != 1 {
		t.Fatalf("expected repository List called, got %d", orderRepo.listCalls)
	}
}
