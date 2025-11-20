package httpclient

import (
	"context"
	"fmt"
	"micro-warehouse/merchant-service/pkg/redis"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type CachedProductCLient struct {
	client ProductClientInterface
	redis  *redis.RedisClient
	ttl    time.Duration
}

func NewCachedProductClient(productClient ProductClientInterface, redisClient *redis.RedisClient, ttl time.Duration) *CachedProductCLient {
	return &CachedProductCLient{
		client: productClient,
		redis:  redisClient,
		ttl:    ttl,
	}
}

func (cpc *CachedProductCLient) generateCacheKey(prefix string, id uint) string {
	return fmt.Sprintf("product:%s:%d", prefix, id)
}

func (cpc *CachedProductCLient) generateCacheKeyMultiple(prefix string, ids []uint) string {
	key := fmt.Sprintf("product:%s", prefix)
	for _, id := range ids {
		key += fmt.Sprintf(":%d", id)
	}
	return key[:len(key)-1]
}

func (cpc *CachedProductCLient) GetProductByID(ctx context.Context, productID uint) (*ProductResponse, error) {
	cacheKey := cpc.generateCacheKey("single", productID)

	var cachedProduct ProductResponse
	if err := cpc.redis.Get(ctx, cacheKey, &cachedProduct); err != nil {
		log.Infof("[CachedProductCLient] GetProductByID - 1: %v", err)
		return &cachedProduct, nil
	}

	product, err := cpc.client.GetProductByID(ctx, productID)
	if err != nil {
		log.Errorf("[CachedProductCLient] GetProductByID - 2: %v", err)
		return nil, err
	}

	err = cpc.redis.Set(ctx, cacheKey, product, cpc.ttl)
	if err != nil {
		log.Errorf("[CachedProductCLient] GetProductByID - 3: %v", err)
		return nil, err
	}

	return product, nil
}

func (cpc *CachedProductCLient) GetProducts(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]ProductResponse, error) {
	return cpc.client.GetProducts(ctx, page, limit, search, sortBy, sortOrder)
}

func (cpc *CachedProductCLient) HealthCheck(ctx context.Context) error {
	return cpc.client.HealthCheck(ctx)
}
