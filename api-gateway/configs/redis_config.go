package configs

import (
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

func LoadRedisConfig() RedisConfig {
	db, _ := strconv.Atoi(getRedisEnv("REDIS_DB", "0"))
	poolSize, _ := strconv.Atoi(getRedisEnv("REDIS_POOL_SIZE", "10"))

	return RedisConfig{
		Host:     getRedisEnv("REDIS_HOST", "localhost"),
		Port:     getRedisEnv("REDIS_PORT", "6379"),
		Password: getRedisEnv("REDIS_PASSWORD", ""),
		DB:       db,
		PoolSize: poolSize,
	}
}

func NewRedisClient(config RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})
}

func getRedisEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
