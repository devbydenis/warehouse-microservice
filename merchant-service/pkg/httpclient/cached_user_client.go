package httpclient

import (
	"context"
	"fmt"
	"micro-warehouse/merchant-service/pkg/redis"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type CachedUserClient struct {
	client UserClientInterface
	redis  *redis.RedisClient
	ttl    time.Duration
}

func NewCachedUserClient(UserClient UserClientInterface, redis *redis.RedisClient) *CachedUserClient {
	return &CachedUserClient{
		client: UserClient,
		redis:  redis,
		ttl:    1 * time.Hour,
	}
}

func (cuc *CachedUserClient) generateCacheKey(prefix string, id uint) string {
	return fmt.Sprintf("user:%s:%d", prefix, id)
}

func (cuc *CachedUserClient) GetUserByID(ctx context.Context, userID uint) (*UserResponse, error) {
	cacheKey := cuc.generateCacheKey("single", userID)

	var cachedUser UserResponse
	if err := cuc.redis.Get(ctx, cacheKey, &cachedUser); err == nil {
		log.Infof("[CachedUserClient] GetUserByID - 1: %v", cachedUser)
		return &cachedUser, nil
	}

	user, err := cuc.client.GetUserByID(ctx, userID)
	if err != nil {
		log.Errorf("[CachedUserClient] GetUserByID - 2: %v", err)
		return nil, err
	}

	err = cuc.redis.Set(ctx, cacheKey, &cachedUser, cuc.ttl)
	if err != nil {
		log.Errorf("[CachedUserClient] GetUserByID - 3: %v", err)
		return nil, err
	}

	return user, nil
}
