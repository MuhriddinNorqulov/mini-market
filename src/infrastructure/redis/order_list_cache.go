package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/ports/cache"

	goredis "github.com/redis/go-redis/v9"
)

const orderListCacheTTL = 30 * time.Second
const orderListVersionKey = "order:list:version"

type OrderListCacheImpl struct {
	client *goredis.Client
}

// @inject
func NewOrderListCacheImpl(client *goredis.Client) cache.OrderListCache {
	return &OrderListCacheImpl{client: client}
}

func (this *OrderListCacheImpl) key(ctx context.Context) string {
	version, err := this.client.Get(ctx, orderListVersionKey).Int64()
	if err != nil {
		version = 0
	}
	return fmt.Sprintf("order:list:v%d", version)
}

func (this *OrderListCacheImpl) Get(ctx context.Context) (*entity.PagingEntity[*entity.OrderEntity], bool) {
	raw, err := this.client.Get(ctx, this.key(ctx)).Bytes()
	if err != nil {
		return nil, false
	}
	var page entity.PagingEntity[*entity.OrderEntity]
	if err := json.Unmarshal(raw, &page); err != nil {
		return nil, false
	}
	return &page, true
}

func (this *OrderListCacheImpl) Set(ctx context.Context, page *entity.PagingEntity[*entity.OrderEntity]) {
	raw, err := json.Marshal(page)
	if err != nil {
		return
	}
	this.client.Set(ctx, this.key(ctx), raw, orderListCacheTTL)
}

func (this *OrderListCacheImpl) InvalidateAll(ctx context.Context) {
	this.client.Incr(ctx, orderListVersionKey)
}
