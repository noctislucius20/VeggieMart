package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"notification-service/internal/core/domain/entity"
	"notification-service/utils"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
)

type NotificationCacheInterface interface {
	GetRoleById(ctx context.Context, id int64) (*entity.RoleEntity, error)
}

type notificationCache struct {
	redisClient *redis.Client
	logger      *log.Logger
}

func NewNotificationCache(redisClient *redis.Client, logger *log.Logger) NotificationCacheInterface {
	return &notificationCache{
		redisClient: redisClient,
		logger:      logger,
	}
}

// GetRoleById implements [NotificationCacheInterface].
func (o *notificationCache) GetRoleById(ctx context.Context, id int64) (*entity.RoleEntity, error) {
	var roleEntity entity.RoleEntity

	keyRolePermission := fmt.Sprintf("role:id:%d", id)
	rolePermission, err := o.redisClient.Get(ctx, keyRolePermission).Result()
	if err != nil {
		o.logger.Errorf("[NotificationCache] GetRoleById: %v", err)
		if errors.Is(err, redis.Nil) {
			return nil, utils.ErrRelationDataNotFound
		}
		return nil, err
	}

	json.Unmarshal([]byte(rolePermission), &roleEntity)

	return &roleEntity, nil
}
