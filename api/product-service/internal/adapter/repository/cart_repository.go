package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"product-service/internal/core/domain/entity"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
)

type CartRepositoryInterface interface {
	AddToCart(ctx context.Context, userId int64, req entity.CartItem) error
	GetCart(ctx context.Context, userId int64) ([]entity.CartItem, error)
	RemoveFromCart(ctx context.Context, userId int64, productId int64) error
}

type cartRepository struct {
	redisClient *redis.Client
	logger      *log.Logger
}

func NewCartRepository(redisClient *redis.Client, logger *log.Logger) CartRepositoryInterface {
	return &cartRepository{
		redisClient: redisClient,
		logger:      logger,
	}
}

// AddToCart implements [CartRepositoryInterface].
func (c *cartRepository) AddToCart(ctx context.Context, userId int64, req entity.CartItem) error {
	var (
		key = fmt.Sprintf("user:id:%d:cart:%d", userId, req.ProductID)
	)

	jsonData, err := json.Marshal(req)
	if err != nil {
		c.logger.Errorf("[CartRepository-1] AddToCart: %v", err)
		return err
	}

	if err := c.redisClient.Set(ctx, key, jsonData, 0).Err(); err != nil {
		c.logger.Errorf("[CartRepository-2] AddToCart: %v", err)
		return err
	}

	return nil
}

// GetCart implements [CartRepositoryInterface].
func (c *cartRepository) GetCart(ctx context.Context, userId int64) ([]entity.CartItem, error) {
	var (
		items []entity.CartItem
		key   = fmt.Sprintf("user:id:%d:cart:*", userId)
	)

	val, err := c.redisClient.Get(ctx, key).Result()
	if err != nil {
		c.logger.Errorf("[CartRepository-1] GetCart: %v", err)
		return nil, err
	}

	if err := json.Unmarshal([]byte(val), &items); err != nil {
		c.logger.Errorf("[CartRepository-2] GetCart: %v", err)
		return nil, err
	}

	return items, nil
}

// RemoveFromCart implements [CartRepositoryInterface].
func (c *cartRepository) RemoveFromCart(ctx context.Context, userId int64, productId int64) error {
	var (
		key = fmt.Sprintf("user:id:%d:cart:%d", userId, productId)
	)

	if err := c.redisClient.Del(ctx, key).Err(); err != nil {
		c.logger.Errorf("[CartRepository-1] RemoveFromCart: %v", err)
		return err
	}

	return nil
}
