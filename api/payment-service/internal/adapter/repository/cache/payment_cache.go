package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"payment-service/internal/adapter/repository"
	"payment-service/internal/core/domain/entity"
	"payment-service/utils"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
)

type PaymentCacheInterface interface {
	GetPaymentById(ctx context.Context, paymentId uint, userId uint) (*entity.PaymentEntity, error)
	DeletePaymentCache(ctx context.Context, id int64) error
}

type paymentCache struct {
	redisClient *redis.Client
	repoPayment repository.PaymentRepositoryInterface
	logger      *log.Logger
}

func NewPaymentCache(redisClient *redis.Client, repoPayment repository.PaymentRepositoryInterface, logger *log.Logger) PaymentCacheInterface {
	return &paymentCache{
		redisClient: redisClient,
		repoPayment: repoPayment,
		logger:      logger,
	}
}

// GetPaymentById implements [PaymentCacheInterface].
func (p *paymentCache) GetPaymentById(ctx context.Context, paymentId uint, userId uint) (*entity.PaymentEntity, error) {
	var (
		payment entity.PaymentEntity
		key     = fmt.Sprintf("payment:id:%d", paymentId)
	)

	// Check redis if data exists.
	val, err := p.redisClient.Get(ctx, key).Result()
	if err == nil {
		// if key exists but value null, return data not found error
		if val == "null" {
			err := errors.New(utils.DATA_NOT_FOUND)
			p.logger.Errorf("[PaymentCache-1] GetPaymentById: %v", err)
			return nil, err
		}

		json.Unmarshal([]byte(val), &payment)

		return &payment, nil
	}

	paymentEntity, err := p.repoPayment.GetPaymentById(ctx, paymentId, userId)
	if err != nil {
		if err.Error() == utils.DATA_NOT_FOUND {
			if err := p.redisClient.Set(ctx, key, "null", 1*time.Minute); err != nil {
				p.logger.Errorf("[PaymentCache-2] GetPaymentById: %v", err)
			}
		}

		p.logger.Errorf("[PaymentCache-3] GetPaymentById: %v", err)
		return nil, err
	}

	payment = *paymentEntity

	// Save to redis
	jsonData, _ := json.Marshal(payment)
	ttl := 10*time.Minute + time.Duration(rand.Intn(120))*time.Second
	if err := p.redisClient.Set(ctx, key, jsonData, ttl).Err(); err != nil {
		p.logger.Errorf("[PaymentCache-4] GetPaymentById: %v", err)
	}

	return &payment, nil
}

// DeletePaymentCache implements [PaymentCacheInterface].
func (p *paymentCache) DeletePaymentCache(ctx context.Context, id int64) error {
	var (
		delKeys = []string{
			fmt.Sprintf("payment:id:%d", id),
		}
	)

	if err := p.redisClient.Del(ctx, delKeys...).Err(); err != nil {
		p.logger.Errorf("[PaymentCache-1] DeletePaymentCache: %v", err)
		return err
	}

	return nil
}
