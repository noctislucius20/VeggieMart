package ratelimit

import (
	"api-gateway/response"
	"api-gateway/utils"
	"api-gateway/utils/helper"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

type middlewareRateLimit struct {
	redisClient *redis.Client
}

type MiddlewareJWTInterface interface {
	MiddlewareRateLimiter() echo.MiddlewareFunc

	checkRateLimit(ctx context.Context, ip string, requestPerSecond, windowSize int64) (bool, error)
}

func NewRedisRateLimiter(redisClient *redis.Client) MiddlewareJWTInterface {
	return &middlewareRateLimit{
		redisClient: redisClient,
	}
}

// MiddlewareRateLimiter implements [MiddlewareJWTInterface].
func (rl *middlewareRateLimit) MiddlewareRateLimiter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var (
				requestPerSecond = helper.GetEnvAsFloat("RATE_LIMIT_RPS", 10.0)
				windowSize       = helper.GetEnvAsInt("RATE_LIMIT_WINDOW_SECONDS", 20)
			)

			if c.Request().URL.Path == "/health" {
				return next(c)
			}

			ip := c.RealIP()
			if ip == "" {
				ip = c.Request().RemoteAddr
			}

			allowed, err := rl.checkRateLimit(c.Request().Context(), ip, int64(requestPerSecond), windowSize)
			if err != nil {
				return next(c)
			}

			if !allowed {
				return c.JSON(http.StatusTooManyRequests, response.ResponseFailed(utils.TOO_MANY_REQUESTS))
			}

			return next(c)
		}
	}
}

func (rl *middlewareRateLimit) checkRateLimit(ctx context.Context, ip string, requestPerSecond, windowSize int64) (bool, error) {
	key := fmt.Sprintf("rate_limit:%s", ip)

	now := time.Now().Unix()
	windowStart := fmt.Sprintf("%d", now-int64(windowSize))

	pipe := rl.redisClient.Pipeline()

	pipe.ZRemRangeByScore(ctx, key, "0", windowStart)

	pipe.ZAdd(ctx, key, &redis.Z{
		Score:  float64(now),
		Member: time.Now().UnixNano(),
	})

	countCmd := pipe.ZCard(ctx, key)

	pipe.Expire(ctx, key, time.Duration(windowSize)*time.Second)

	if _, err := pipe.Exec(ctx); err != nil {
		return false, err
	}

	count := countCmd.Val()
	limit := int64(float64(requestPerSecond) * float64(windowSize))

	if count > limit {
		return false, nil
	}

	return true, nil
}
