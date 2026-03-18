package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisMGetter interface {
	MGet(ctx context.Context, keys ...string) *redis.SliceCmd
}

func RedisBulkGet(ctx context.Context, client RedisMGetter, keys []string, dest ...any) error {
	if len(keys) != len(dest) {
		return errors.New("keys and destination length mismatch")
	}

	vals, err := client.MGet(ctx, keys...).Result()
	if err != nil {
		return err
	}

	for i, v := range vals {
		if v == nil {
			continue
		}

		if err := json.Unmarshal([]byte(v.(string)), &dest[i]); err != nil {
			return err
		}
	}

	return nil
}

func AppendKeyIfNotEmpty(keys []string, format string, value string) []string {
	if value == "" {
		return keys
	}

	return append(keys, fmt.Sprintf(format, value))
}
