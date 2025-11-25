package httpclient

import (
	"context"
	"fmt"
	"micro-warehouse/merchant-service/pkg/redis"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type CachedWarehouseClient struct {
	client WarehouseClientInterface
	redis *redis.RedisClient
	ttl time.Duration
}

func NewCachedWarehouseClient(warehouseClient WarehouseClientInterface, redisClient *redis.RedisClient, ttl time.Duration) *CachedWarehouseClient {
	return &CachedWarehouseClient{
		client: warehouseClient,
		redis:  redisClient,
		ttl:    1 * time.Hour,
	}
}

func (cwc *CachedWarehouseClient) generateCacheKey(prefix string, id uint) string {
	return fmt.Sprintf("warehouse:%s:%d", prefix, id)
}

func (cwc *CachedWarehouseClient) GetWarehouseByID(ctx context.Context, warehouseID uint) (*WarehouseResponse, error) {
	cacheKey := cwc.generateCacheKey("single", warehouseID)

	var cachedWarehouse WarehouseResponse
	if err := cwc.redis.Get(ctx, cacheKey, &cachedWarehouse); err == nil {
		log.Infof("[CachedWarehouseClient] GetWarehouseByID - 1: %v", cachedWarehouse)
		return &cachedWarehouse, nil
	}

	warehouse, err := cwc.client.GetWarehouseByID(ctx, warehouseID)
	if err != nil {
		log.Errorf("[CachedWarehouseClient] GetWarehouseByID - 2: %v", err)
		return nil, err
	}

	err = cwc.redis.Set(ctx, cacheKey, warehouse, cwc.ttl)
	if err != nil {
		log.Errorf("[CachedWarehouseClient] GetWarehouseByID - 3: %v", err)
		return nil, err
	}

	return warehouse, nil
}

func (cwc *CachedWarehouseClient) GetWarehouseProductStock(ctx context.Context, warehouseID uint, productID uint) (*WarehouseProductStockResponse, error) {
	cacheKey := cwc.generateCacheKey("single", warehouseID)

	var cachedWarehouseProductStock WarehouseProductStockResponse
	if err := cwc.redis.Get(ctx, cacheKey, &cachedWarehouseProductStock); err == nil {
		log.Infof("[CachedWarehouseClient] GetWarehouseProductStock - 1: %v", cachedWarehouseProductStock)
		return &cachedWarehouseProductStock, nil
	}

	warehouseProductStock, err := cwc.client.GetWarehouseProductStock(ctx, warehouseID, productID)
	if err != nil {
		log.Errorf("[CachedWarehouseClient] GetWarehouseProductStock - 2: %v", err)
		return nil, err
	}

	err = cwc.redis.Set(ctx, cacheKey, warehouseProductStock, cwc.ttl)
	if err != nil {
		log.Errorf("[CachedWArehouseClient] GetWarehouseProductStock - 3: %v", err)
		return nil, err
	}

	return warehouseProductStock, nil
}