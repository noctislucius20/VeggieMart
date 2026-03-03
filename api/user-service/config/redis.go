package config

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func (cfg Config) NewRedisClient(ctx context.Context) (*redis.Client, error) {
	connect := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	client := redis.NewClient(&redis.Options{
		Addr: connect,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
