package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"user-service/internal/adapter/repository"
	"user-service/internal/core/domain/entity"
	"user-service/utils"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
)

type RoleCacheInterface interface {
	GetRoleById(ctx context.Context, id int64) (*entity.RoleEntity, error)
	GetRoleByName(ctx context.Context, name string) (*entity.RoleEntity, error)
}

type roleCache struct {
	redisClient *redis.Client
	repoRole    repository.RoleRepositoryInterface
	logger      *log.Logger
}

// GetRoleByName implements [RoleCacheInterface].
func (r *roleCache) GetRoleByName(ctx context.Context, name string) (*entity.RoleEntity, error) {
	var (
		role *entity.RoleEntity
		key  = fmt.Sprintf("role:name:%s", name)
	)

	// Check redis if data exists.
	val, err := r.redisClient.Get(ctx, key).Result()
	if err == nil {
		// if key exists but value null, return data not found error
		if val == "null" {
			err := errors.New(utils.DATA_NOT_FOUND)
			r.logger.Errorf("[RoleCache-1] GetRoleByName: %v", err)
			return nil, err
		}

		json.Unmarshal([]byte(val), &role)
		return role, nil
	}

	roleEntity, err := r.repoRole.GetRoleByIdOrName(ctx, 0, name)
	if err != nil {
		// Save to redis (create key with null value if data not found)
		if err.Error() == utils.DATA_NOT_FOUND {
			if err := r.redisClient.Set(ctx, key, "null", 1*time.Minute).Err(); err != nil {
				r.logger.Errorf("[RoleCache-2] GetRoleByName: %v", err)
			}
		}

		r.logger.Errorf("[RoleCache-3] GetRoleByName: %v", err)
		return nil, err
	}

	role = roleEntity

	// Save to redis
	jsonData, _ := json.Marshal(role)
	if err := r.redisClient.Set(ctx, key, jsonData, 1*time.Hour).Err(); err != nil {
		r.logger.Errorf("[RoleCache-4] GetRoleByName: %v", err)
	}

	return role, nil
}

// GetRoleById implements [RoleCacheInterface].
func (r *roleCache) GetRoleById(ctx context.Context, id int64) (*entity.RoleEntity, error) {
	var (
		role *entity.RoleEntity
		key  = fmt.Sprintf("role:id:%d", id)
	)

	// Check redis if data exists.
	val, err := r.redisClient.Get(ctx, key).Result()
	if err == nil {
		// if key exists but value null, return data not found error
		if val == "null" {
			err := errors.New(utils.DATA_NOT_FOUND)
			r.logger.Errorf("[RoleCache-1] GetRoleById: %v", err)
			return nil, err
		}

		json.Unmarshal([]byte(val), &role)
		return role, nil
	}

	roleEntity, err := r.repoRole.GetRoleByIdOrName(ctx, id, "")
	if err != nil {
		// Save to redis (create key with null value if data not found)
		if err.Error() == utils.DATA_NOT_FOUND {
			if err := r.redisClient.Set(ctx, key, "null", 1*time.Minute).Err(); err != nil {
				r.logger.Errorf("[RoleCache-2] GetRoleById: %v", err)
			}
		}

		r.logger.Errorf("[RoleCache-3] GetRoleById: %v", err)
		return nil, err
	}

	role = roleEntity

	// Save to redis
	jsonData, _ := json.Marshal(role)
	if err := r.redisClient.Set(ctx, key, jsonData, 1*time.Hour).Err(); err != nil {
		r.logger.Errorf("[RoleCache-4] GetRoleById: %v", err)
	}

	return role, nil
}

func NewRoleCache(redisClient *redis.Client, repoRole repository.RoleRepositoryInterface, logger *log.Logger) RoleCacheInterface {
	return &roleCache{
		redisClient: redisClient,
		repoRole:    repoRole,
		logger:      logger,
	}
}
