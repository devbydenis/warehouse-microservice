package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"micro-warehouse/merchant-service/configs"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(cfg configs.Config) *RedisClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
			Password: "",
			DB:       0,
		}),
	}
}

func (rc *RedisClient) Ping(ctx context.Context) error {
	_, err := rc.client.Ping(ctx).Result()
	if err != nil {
		log.Errorf("[RedisClient] Ping - 1: %v", err)
		return err
	}

	return nil
}

func (rc *RedisClient) Get(ctx context.Context, key string, value interface{}) error {
	jsonData, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		log.Errorf("[RedisClient] Get - 1: %v", err)
		return err
	}

	err = json.Unmarshal([]byte(jsonData), value)
	if err != nil {
		log.Errorf("[RedisClient] Get - 2: %v", err)
		return err
	}

	return nil
}

func (rc *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		log.Errorf("[RedisClient] Set - 1: %v", err)
		return err
	}

	err = rc.client.Set(ctx, key, jsonData, expiration).Err()
	if err != nil {
		log.Errorf("[RedisClient] Set - 2: %v", err)
		return err
	}

	return nil
}

func (rc *RedisClient) Delete(ctx context.Context, key string) error {
	err := rc.client.Del(ctx, key).Err()
	if err != nil {
		log.Errorf("[RedisClient] Delete - 1: %v", err)
		return err
	}

	return nil
}

func (rc *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := rc.client.Exists(ctx, key).Result()
	if err != nil {
		log.Errorf("[RedisClient] Exists - 1: %v", err)
		return false, err
	}

	return exists > 0, nil
}

func (rc *RedisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := rc.client.TTL(ctx, key).Result()
	if err != nil {
		log.Errorf("[RedisClient] TTL - 1: %v", err)
		return 0, err
	}

	return ttl, nil
}

func (rc *RedisClient) Close(ctx context.Context) error {
	return rc.client.Close()
}

func (rc *RedisClient) FlushAll(ctx context.Context) error {
	err := rc.client.FlushAll(ctx).Err()
	if err != nil {
		log.Errorf("[RedisClient] FlushAll - 1: %v", err)
		return err
	}

	return nil
}