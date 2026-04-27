package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"order-service/internal/adapter/repository"
	"order-service/internal/core/domain/entity"
	"order-service/utils"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
)

type OrderCacheInterface interface {
	GetOrderById(ctx context.Context, orderId int64, userId int64) (*entity.OrderEntity, error)
	GetRoleById(ctx context.Context, id int64) (*entity.RoleEntity, error)
	GetOrderByOrderCode(ctx context.Context, orderCode string, userId int64) (*entity.OrderEntity, error)
	DeleteOrderCache(ctx context.Context, id int64, orderCode string) error
}

type orderCache struct {
	redisClient *redis.Client
	repoOrder   repository.OrderRepositoryInterface
	logger      *log.Logger
}

func NewOrderCache(redisClient *redis.Client, repoOrder repository.OrderRepositoryInterface, logger *log.Logger) OrderCacheInterface {
	return &orderCache{
		redisClient: redisClient,
		repoOrder:   repoOrder,
		logger:      logger,
	}
}

// GetRoleById implements [OrderCacheInterface].
func (o *orderCache) GetRoleById(ctx context.Context, id int64) (*entity.RoleEntity, error) {
	var roleEntity entity.RoleEntity

	keyRolePermission := fmt.Sprintf("role:id:%d", id)
	rolePermission, err := o.redisClient.Get(ctx, keyRolePermission).Result()
	if err != nil {
		o.logger.Errorf("[OrderCache-1] GetRoleById: %v", err)
		if errors.Is(err, redis.Nil) {
			return nil, errors.New(utils.RELATION_DATA_NOT_FOUND)
		}
		return nil, err
	}

	json.Unmarshal([]byte(rolePermission), &roleEntity)

	return &roleEntity, nil
}

// GetOrderByOrderCode implements [OrderCacheInterface].
func (o *orderCache) GetOrderByOrderCode(ctx context.Context, orderCode string, userId int64) (*entity.OrderEntity, error) {
	var (
		order entity.OrderEntity
		key   = fmt.Sprintf("order:ordercode:%s", orderCode)
	)

	// Check redis if data exists.
	val, err := o.redisClient.Get(ctx, key).Result()
	if err == nil {
		// if key exists but value null, return data not found error
		if val == "null" {
			err := errors.New(utils.DATA_NOT_FOUND)
			o.logger.Errorf("[OrderCache-1] GetOrderByOrderCode: %v", err)
			return nil, err
		}

		json.Unmarshal([]byte(val), &order)

		return &order, nil
	}

	orderEntity, err := o.repoOrder.GetOrderByOrderCode(ctx, orderCode, userId)
	if err != nil {
		if err.Error() == utils.DATA_NOT_FOUND {
			if err := o.redisClient.Set(ctx, key, "null", 10*time.Minute); err != nil {
				o.logger.Errorf("[OrderCache-2] GetOrderByOrderCode: %v", err)
			}
		}

		o.logger.Errorf("[OrderCache-3] GetOrderByOrderCode: %v", err)
		return nil, err
	}

	order = *orderEntity

	// Save to redis
	jsonData, _ := json.Marshal(order)
	ttl := 10*time.Minute + time.Duration(rand.Intn(120))*time.Second
	if err := o.redisClient.Set(ctx, key, jsonData, ttl).Err(); err != nil {
		o.logger.Errorf("[OrderCache-4] GetOrderByOrderCode: %v", err)
	}

	return &order, nil
}

// DeleteOrderCache implements [OrderCacheInterface].
func (o *orderCache) DeleteOrderCache(ctx context.Context, id int64, orderCode string) error {
	var (
		delKeys = []string{
			fmt.Sprintf("order:id:%d", id),
			fmt.Sprintf("order:ordercode:%s", orderCode),
		}
	)

	if err := o.redisClient.Del(ctx, delKeys...).Err(); err != nil {
		o.logger.Errorf("[OrderCache-1] DeleteOrderCache: %v", err)
		return err
	}

	return nil
}

// GetOrderById implements [OrderCacheInterface].
func (o *orderCache) GetOrderById(ctx context.Context, orderId int64, userId int64) (*entity.OrderEntity, error) {
	var (
		order entity.OrderEntity
		key   = fmt.Sprintf("order:id:%d", orderId)
	)

	// Check redis if data exists.
	val, err := o.redisClient.Get(ctx, key).Result()
	if err == nil {
		// if key exists but value null, return data not found error
		if val == "null" {
			err := errors.New(utils.DATA_NOT_FOUND)
			o.logger.Errorf("[OrderCache-1] GetOrderById: %v", err)
			return nil, err
		}

		json.Unmarshal([]byte(val), &order)

		return &order, nil
	}

	orderEntity, err := o.repoOrder.GetOrderById(ctx, orderId, userId)
	if err != nil {
		if err.Error() == utils.DATA_NOT_FOUND {
			if err := o.redisClient.Set(ctx, key, "null", 10*time.Minute); err != nil {
				o.logger.Errorf("[OrderCache-2] GetOrderById: %v", err)
			}
		}

		o.logger.Errorf("[OrderCache-3] GetOrderById: %v", err)
		return nil, err
	}

	order = *orderEntity

	// Save to redis
	jsonData, _ := json.Marshal(order)
	ttl := 10*time.Minute + time.Duration(rand.Intn(120))*time.Second
	if err := o.redisClient.Set(ctx, key, jsonData, ttl).Err(); err != nil {
		o.logger.Errorf("[OrderCache-4] GetOrderById: %v", err)
	}

	return &order, nil
}
