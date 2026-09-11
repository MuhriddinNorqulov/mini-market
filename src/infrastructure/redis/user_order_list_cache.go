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

const userOrderListCacheTTL = 30 * time.Second

type UserOrderListCacheImpl struct {
	client *goredis.Client
}

// @inject
func NewUserOrderListCacheImpl(client *goredis.Client) cache.UserOrderListCache {
	return &UserOrderListCacheImpl{client: client}
}

func (this *UserOrderListCacheImpl) versionKey(userID uint) string {
	return fmt.Sprintf("order:list:user:%d:version", userID)
}

func (this *UserOrderListCacheImpl) key(ctx context.Context, userID uint) string {
	version, err := this.client.Get(ctx, this.versionKey(userID)).Int64()
	if err != nil {
		version = 0
	}
	return fmt.Sprintf("order:list:user:%d:v%d", userID, version)
}

func (this *UserOrderListCacheImpl) Get(ctx context.Context, userID uint) (*entity.PagingEntity[*entity.OrderEntity], bool) {
	raw, err := this.client.Get(ctx, this.key(ctx, userID)).Bytes()
	if err != nil {
		return nil, false
	}
	var page entity.PagingEntity[*entity.OrderEntity]
	if err := json.Unmarshal(raw, &page); err != nil {
		return nil, false
	}
	return &page, true
}

func (this *UserOrderListCacheImpl) Set(ctx context.Context, userID uint, page *entity.PagingEntity[*entity.OrderEntity]) {
	raw, err := json.Marshal(page)
	if err != nil {
		return
	}
	this.client.Set(ctx, this.key(ctx, userID), raw, userOrderListCacheTTL)
}

func (this *UserOrderListCacheImpl) InvalidateAll(ctx context.Context, userID uint) {
	this.client.Incr(ctx, this.versionKey(userID))
}
