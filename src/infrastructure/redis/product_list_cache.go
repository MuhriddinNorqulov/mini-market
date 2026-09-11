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

const productListCacheTTL = 30 * time.Second
const productListVersionKey = "product:list:version"

type ProductListCacheImpl struct {
	client *goredis.Client
}

// @inject
func NewProductListCacheImpl(client *goredis.Client) cache.ProductListCache {
	return &ProductListCacheImpl{client: client}
}

func (this *ProductListCacheImpl) key(ctx context.Context) string {
	version, err := this.client.Get(ctx, productListVersionKey).Int64()
	if err != nil {
		version = 0
	}
	return fmt.Sprintf("product:list:v%d", version)
}

func (this *ProductListCacheImpl) Get(ctx context.Context) (*entity.PagingEntity[*entity.ProductEntity], bool) {
	raw, err := this.client.Get(ctx, this.key(ctx)).Bytes()
	if err != nil {
		return nil, false
	}
	var page entity.PagingEntity[*entity.ProductEntity]
	if err := json.Unmarshal(raw, &page); err != nil {
		return nil, false
	}
	return &page, true
}

func (this *ProductListCacheImpl) Set(ctx context.Context, page *entity.PagingEntity[*entity.ProductEntity]) {
	raw, err := json.Marshal(page)
	if err != nil {
		return
	}
	this.client.Set(ctx, this.key(ctx), raw, productListCacheTTL)
}

func (this *ProductListCacheImpl) InvalidateAll(ctx context.Context) {
	this.client.Incr(ctx, productListVersionKey)
}
