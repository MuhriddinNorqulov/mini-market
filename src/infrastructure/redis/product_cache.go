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

const productCacheTTL = 60 * time.Second

type ProductCache struct {
	client *goredis.Client
}

// @inject
func NewProductCache(client *goredis.Client) cache.ProductCache {
	return &ProductCache{client: client}
}

func (this *ProductCache) key(id uint) string {
	return fmt.Sprintf("product:%d", id)
}

func (this *ProductCache) Get(ctx context.Context, id uint) (*entity.ProductEntity, bool) {
	raw, err := this.client.Get(ctx, this.key(id)).Bytes()
	if err != nil {
		return nil, false
	}
	var product entity.ProductEntity
	if err := json.Unmarshal(raw, &product); err != nil {
		return nil, false
	}
	return &product, true
}

func (this *ProductCache) Set(ctx context.Context, product *entity.ProductEntity) {
	raw, err := json.Marshal(product)
	if err != nil {
		return
	}
	this.client.Set(ctx, this.key(product.ID), raw, productCacheTTL)
}

func (this *ProductCache) Invalidate(ctx context.Context, id uint) {
	this.client.Del(ctx, this.key(id))
}
