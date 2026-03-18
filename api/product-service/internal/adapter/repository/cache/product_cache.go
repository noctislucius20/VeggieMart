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
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
)

type ProductCacheInterface interface {
	GetProductById(ctx context.Context, id int64) (*entity.ProductEntity, error)
	DeleteProductCache(ctx context.Context, id int64) error
}

type productCache struct {
	redisClient *redis.Client
	repoProduct repository.ProductRepositoryInterface
	logger      *log.Logger
}

func NewProductCache(redisClient *redis.Client, repoProduct repository.ProductRepositoryInterface, logger *log.Logger) ProductCacheInterface {
	return &productCache{
		redisClient: redisClient,
		repoProduct: repoProduct,
		logger:      logger,
	}
}

// DeleteProductCache implements [ProductCacheInterface].
func (p *productCache) DeleteProductCache(ctx context.Context, id int64) error {
	var (
		delKeys = []string{
			fmt.Sprintf("product:id:%d", id),
		}
	)

	if err := p.redisClient.Del(ctx, delKeys...).Err(); err != nil {
		p.logger.Errorf("[ProductCache-1] DeleteProductCache: %v", err)
		return err
	}

	return nil
}

// GetProductById implements [ProductCacheInterface].
func (p *productCache) GetProductById(ctx context.Context, id int64) (*entity.ProductEntity, error) {
	var (
		product entity.ProductEntity
		key     = fmt.Sprintf("product:id:%d", id)
	)

	// Check redis if data exists.
	val, err := p.redisClient.Get(ctx, key).Result()
	if err == nil {
		// if key exists but value null, return data not found error
		if val == "null" {
			err := errors.New(utils.DATA_NOT_FOUND)
			p.logger.Errorf("[ProductCache-1] GetProductById: %v", err)
			return nil, err
		}

		json.Unmarshal([]byte(val), &product)

		return &product, nil
	}

	productEntity, err := p.repoProduct.GetProductById(ctx, id)
	if err != nil {
		if err.Error() == utils.DATA_NOT_FOUND {
			if err := p.redisClient.Set(ctx, key, "null", 10*time.Minute); err != nil {
				p.logger.Errorf("[ProductCache-2] GetProductById: %v", err)
			}
		}

		p.logger.Errorf("[ProductCache-3] GetProductById: %v", err)
		return nil, err
	}

	product = *productEntity

	// Save to redis
	jsonData, _ := json.Marshal(product)
	ttl := 10*time.Minute + time.Duration(rand.Intn(120))*time.Second
	if err := p.redisClient.Set(ctx, key, jsonData, ttl).Err(); err != nil {
		p.logger.Errorf("[ProductCache-4] GetProductById: %v", err)
	}

	return &product, nil
}
