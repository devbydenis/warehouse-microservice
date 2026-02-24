package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type RedisRateLimiterConfig struct {
	Max         int
	Expiration  time.Duration
	KeyPrefix   string
	RedisClient *redis.Client
}

func DefaultRateLimiterConfig() RedisRateLimiterConfig {
	return RedisRateLimiterConfig{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyPrefix:  "api_rate_limit",
		RedisClient: nil,
	}
}

func RedisRateLimiter(config RedisRateLimiterConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		key := fmt.Sprintf("%s:%s", config.KeyPrefix, c.IP())

		current, err := config.RedisClient.Get(ctx, key).Result()
		if err != nil && err != redis.Nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": "Failed to get rate limit",
			})
		}

		if err == redis.Nil {
			err = config.RedisClient.Set(ctx, key, 1, config.Expiration).Err()
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": "Failed to set rate limit",
				})
			}
			return c.Next()
		}

		currentCount, err := strconv.Atoi(current)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": "Failed to parse rate limit",
			})
		}

		if currentCount >= config.Max {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Too Many Requests",
				"message":     "Rate limit exceeded. Please try again later",
				"retry_after": config.Expiration.Seconds(),
				"limit":       config.Max,
				"current":     currentCount,
			})
		}

		err = config.RedisClient.Incr(ctx, key).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": "Failed to increment rate limit",
			})
		}

		if currentCount == 0 {
			err = config.RedisClient.Expire(ctx, key, config.Expiration).Err()
			if err != nil {
				fmt.Printf("Failed to set expiration for key %s: %v\n", key, err)
			}
		}

		c.Set("X-RateLimit-Limit", strconv.Itoa(config.Max))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(config.Max-currentCount-1))
		c.Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(config.Expiration).Unix(), 10))

		return c.Next()
	}
}

func RedisGlobalRateLimiter(config RedisRateLimiterConfig) fiber.Handler {
	config.KeyPrefix = "global"
	return RedisRateLimiter(config)
}

func RedisAuthRateLimiter(config RedisRateLimiterConfig) fiber.Handler {
	config.KeyPrefix = "auth"
	return RedisRateLimiter(config)
}

func RedisAPIRateLimiter(config RedisRateLimiterConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		userID := c.Locals("user_id")

		var key string
		if userID != nil {
			key = fmt.Sprintf("api:%s:%d", c.IP(), userID.(uint))
		} else {
			key = fmt.Sprintf("api:%s", c.IP())
		}

		current, err := config.RedisClient.Get(ctx, key).Result()
		if err != nil && err != redis.Nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": "Failed to get rate limit",
			})
		}

		if err == redis.Nil {
			err = config.RedisClient.Set(ctx, key, 1, config.Expiration).Err()
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": "Failed to set rate limit",
				})
			}
			return c.Next()
		}

		currentCount, err := strconv.Atoi(current)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": "Failed to parse rate limit",
			})
		}

		if currentCount >= config.Max {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Too Many Requests",
				"message":     "Rate limit exceeded. Please try again later",
				"retry_after": config.Expiration.Seconds(),
				"limit":       config.Max,
				"current":     currentCount,
			})
		}

		err = config.RedisClient.Incr(ctx, key).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": "Failed to increment rate limit",
			})
		}

		if currentCount == 0 {
			err = config.RedisClient.Expire(ctx, key, config.Expiration).Err()
			if err != nil {
				fmt.Printf("Failed to set expiration for key %s: %v\n", key, err)
			}
		}

		c.Set("X-RateLimit-Limit", strconv.Itoa(config.Max))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(config.Max-currentCount-1))
		c.Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(config.Expiration).Unix(), 10))

		return c.Next()
	}
}
