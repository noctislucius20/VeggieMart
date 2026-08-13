package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"product-service/internal/adapter/repository"
	"product-service/internal/core/domain/entity"
	"product-service/utils"
	"product-service/utils/helper"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
)

type CategoryCacheInterface interface {
	GetCategoryById(ctx context.Context, id int64) (*entity.CategoryEntity, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*entity.CategoryEntity, error)
	DeleteCategoryCache(ctx context.Context, id int64) error
}

type categoryCache struct {
	redisClient  *redis.Client
	repoCategory repository.CategoryRepositoryInterface
	logger       *log.Logger
}

// DeleteCategoryCache implements [CategoryCacheInterface].
func (c *categoryCache) DeleteCategoryCache(ctx context.Context, id int64) error {
	var (
		category entity.CategoryEntity
		keys     = []string{
			fmt.Sprintf("category:id:%d:slug", id),
		}
		delKeys = []string{
			fmt.Sprintf("category:id:%d", id),
		}
	)

	if err := helper.RedisBulkGet(ctx, c.redisClient, keys, &category); err != nil {
		c.logger.Errorf("[CategoryCache] DeleteCategoryCache: %v", err)
		return err
	}

	delKeys = helper.AppendKeyIfNotEmpty(delKeys, "category:slug:%s", category.Slug)
	delKeys = append(delKeys, keys...)

	if err := c.redisClient.Del(ctx, delKeys...).Err(); err != nil {
		c.logger.Errorf("[CategoryCache] DeleteCategoryCache: %v", err)
		return err
	}

	return nil
}

// GetCategoryById implements [CategoryCacheInterface].
func (c *categoryCache) GetCategoryById(ctx context.Context, id int64) (*entity.CategoryEntity, error) {
	var (
		category *entity.CategoryEntity
		key      = fmt.Sprintf("category:id:%d", id)
	)

	// Check redis if data exists.
	val, err := c.redisClient.Get(ctx, key).Result()
	if err == nil {
		// if key exists but value null, return data not found error
		if val == "null" {
			err := utils.ErrDataNotFound
			c.logger.Errorf("[CategoryCache] GetCategoryById: %v", err)
			return nil, err
		}

		json.Unmarshal([]byte(val), &category)
		return category, nil
	}

	categoryEntity, err := c.repoCategory.GetCategoryByIdOrSlug(ctx, id, "")
	if err != nil {
		// Save to redis (create key with null value if data not found)
		if errors.Is(err, utils.ErrDataNotFound) {
			if err := c.redisClient.Set(ctx, key, "null", 1*time.Minute).Err(); err != nil {
				c.logger.Errorf("[CategoryCache] GetCategoryById: %v", err)
			}
		}

		c.logger.Errorf("[CategoryCache] GetCategoryById: %v", err)
		return nil, err
	}

	category = categoryEntity

	// Save to redis
	jsonData, _ := json.Marshal(category)
	if err := c.redisClient.Set(ctx, key, jsonData, 1*time.Hour).Err(); err != nil {
		c.logger.Errorf("[CategoryCache] GetCategoryById: %v", err)
	}

	return category, nil
}

// GetCategoryBySlug implements [CategoryCacheInterface].
func (c *categoryCache) GetCategoryBySlug(ctx context.Context, slug string) (*entity.CategoryEntity, error) {
	var (
		category entity.CategoryEntity
		key      = fmt.Sprintf("category:slug:%s", slug)
	)

	// Check redis if data exists.
	if valKey, err := c.redisClient.Get(ctx, key).Result(); err == nil {
		// if key exists but value null, return data not found error
		if valKey == "null" {
			err := utils.ErrDataNotFound
			c.logger.Errorf("[CategoryCache] GetCategoryBySlug: %v", err)

			return nil, err
		}
		if val, err := c.redisClient.Get(ctx, fmt.Sprintf("category:id:%s:slug", valKey)).Result(); err == nil {
			json.Unmarshal([]byte(val), &category)

			return &category, nil
		}
	}

	categoryEntity, err := c.repoCategory.GetCategoryByIdOrSlug(ctx, 0, slug)
	if err != nil {
		// Save to redis (create key with null value if data not found)
		if errors.Is(err, utils.ErrDataNotFound) {
			if err := c.redisClient.Set(ctx, key, "null", 1*time.Minute).Err(); err != nil {
				c.logger.Errorf("[CategoryCache] GetCategoryBySlug: %v", err)
			}
		}

		c.logger.Errorf("[CategoryCache] GetCategoryBySlug: %v", err)
		return nil, err
	}

	category = *categoryEntity

	// Save to redis
	jsonData, _ := json.Marshal(category)
	ttl := 1*time.Hour + time.Duration(rand.Intn(120))*time.Second
	pipe := c.redisClient.Pipeline()
	if err := pipe.Set(ctx, key, category.ID, ttl).Err(); err != nil {
		c.logger.Errorf("[CategoryCache] GetCategoryBySlug: %v", err)
	}
	if err := pipe.Set(ctx, fmt.Sprintf("category:id:%d:slug", category.ID), jsonData, ttl).Err(); err != nil {
		c.logger.Errorf("[CategoryCache] GetCategoryBySlug: %v", err)
	}
	if _, err = pipe.Exec(ctx); err != nil {
		c.logger.Errorf("[CategoryCache] GetCategoryBySlug: %v", err)
	}

	return &category, nil
}

func NewCategoryCache(redisClient *redis.Client, repoCategory repository.CategoryRepositoryInterface, logger *log.Logger) CategoryCacheInterface {
	return &categoryCache{
		redisClient:  redisClient,
		repoCategory: repoCategory,
		logger:       logger,
	}
}
