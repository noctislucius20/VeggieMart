package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"
	"user-service/internal/adapter/repository"
	"user-service/internal/core/domain/entity"
	"user-service/utils"
	"user-service/utils/helper"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
)

type RoleCacheInterface interface {
	GetRoleById(ctx context.Context, id int64) (*entity.RoleEntity, error)
	GetRoleByName(ctx context.Context, name string) (*entity.RoleEntity, error)
	DeleteRoleCache(ctx context.Context, id int64) error
}

type roleCache struct {
	redisClient *redis.Client
	repoRole    repository.RoleRepositoryInterface
	logger      *log.Logger
}

func NewRoleCache(redisClient *redis.Client, repoRole repository.RoleRepositoryInterface, logger *log.Logger) RoleCacheInterface {
	return &roleCache{
		redisClient: redisClient,
		repoRole:    repoRole,
		logger:      logger,
	}
}

// DeleteRoleCache implements [RoleCacheInterface].
func (r *roleCache) DeleteRoleCache(ctx context.Context, id int64) error {
	var (
		role entity.RoleEntity
		keys = []string{
			fmt.Sprintf("role:id:%d:name", id),
		}
		delKeys = []string{
			fmt.Sprintf("role:id:%d", id),
		}
	)

	if err := helper.RedisBulkGet(ctx, r.redisClient, keys, &role); err != nil {
		r.logger.Errorf("[RoleCache-1] DeleteRoleCache: %v", err)
		return err
	}

	delKeys = helper.AppendKeyIfNotEmpty(delKeys, "role:name:%s", role.Name)
	delKeys = append(delKeys, keys...)

	if err := r.redisClient.Del(ctx, delKeys...).Err(); err != nil {
		r.logger.Errorf("[RoleCache-2] DeleteRoleCache: %v", err)
		return err
	}

	return nil
}

// GetRoleByName implements [RoleCacheInterface].
func (r *roleCache) GetRoleByName(ctx context.Context, name string) (*entity.RoleEntity, error) {
	var (
		role entity.RoleEntity
		key  = fmt.Sprintf("role:name:%s", name)
	)

	// Check redis if data exists.
	if valKey, err := r.redisClient.Get(ctx, key).Result(); err == nil {
		// if key exists but value null, return data not found error
		if valKey == "null" {
			err := errors.New(utils.DATA_NOT_FOUND)
			r.logger.Errorf("[RoleCache-1] GetRoleByName: %v", err)

			return nil, err
		}
		if val, err := r.redisClient.Get(ctx, fmt.Sprintf("role:id:%s:name", valKey)).Result(); err == nil {
			json.Unmarshal([]byte(val), &role)

			return &role, nil
		}
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

	role = *roleEntity

	// Save to redis
	jsonData, _ := json.Marshal(role)
	ttl := 1*time.Hour + time.Duration(rand.Intn(120))*time.Second
	pipe := r.redisClient.Pipeline()
	if err := pipe.Set(ctx, key, role.ID, ttl).Err(); err != nil {
		r.logger.Errorf("[RoleCache-4] GetRoleByName: %v", err)
	}
	if err := pipe.Set(ctx, fmt.Sprintf("role:id:%d:name", role.ID), jsonData, ttl).Err(); err != nil {
		r.logger.Errorf("[RoleCache-5] GetRoleByName: %v", err)
	}
	if _, err = pipe.Exec(ctx); err != nil {
		r.logger.Errorf("[RoleCache-6] GetRoleByName: %v", err)
	}

	return &role, nil
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
	if err := r.redisClient.Set(ctx, key, jsonData, 0).Err(); err != nil {
		r.logger.Errorf("[RoleCache-4] GetRoleById: %v", err)
	}

	return role, nil
}
