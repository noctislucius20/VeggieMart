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
	RemoveAllFromCart(ctx context.Context, userId int64) error
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
		c.logger.Errorf("[CartRepository] AddToCart: %v", err)
		return err
	}

	if err := c.redisClient.Set(ctx, key, jsonData, 0).Err(); err != nil {
		c.logger.Errorf("[CartRepository] AddToCart: %v", err)
		return err
	}

	return nil
}

// GetCart implements [CartRepositoryInterface].
func (c *cartRepository) GetCart(ctx context.Context, userId int64) ([]entity.CartItem, error) {
	var (
		items        []entity.CartItem
		matchPattern = fmt.Sprintf("user:id:%d:cart:*", userId)
		cursor       uint64
		allKeys      []string
	)

	// 1. Lakukan Scan secara bertahap sampai semua key yang cocok terkumpul
	for {
		// Menggunakan SCAN dengan pattern dan count (limit per batch fetch)
		keys, nextCursor, err := c.redisClient.Scan(ctx, cursor, matchPattern, 100).Result()
		if err != nil {
			c.logger.Errorf("[CartRepository] GetCart: %v", err)
			return nil, err
		}

		allKeys = append(allKeys, keys...)
		cursor = nextCursor

		// Jika cursor kembali ke 0, artinya scan sudah selesai
		if cursor == 0 {
			break
		}
	}

	// Jika tidak ada key yang ditemukan, kembalikan slice kosong (bukan error)
	if len(allKeys) == 0 {
		return []entity.CartItem{}, nil
	}

	// 2. Ambil semua data dari list key yang ditemukan menggunakan MGet
	values, err := c.redisClient.MGet(ctx, allKeys...).Result()
	if err != nil {
		c.logger.Errorf("[CartRepository] GetCart: %v", err)
		return nil, err
	}

	// 3. Iterasi hasil MGet dan masukkan ke dalam struct slice
	for _, val := range values {
		// MGet mengembalikan tipe interface{}, pastikan nilainya ada dan tipenya string
		if valStr, ok := val.(string); ok && valStr != "" {
			var item entity.CartItem
			if err := json.Unmarshal([]byte(valStr), &item); err != nil {
				c.logger.Errorf("[CartRepository] GetCart: %v", err)
				continue // lanjut ke item berikutnya jika ada satu yang corrupt
			}
			items = append(items, item)
		}
	}

	return items, nil
}

// RemoveFromCart implements [CartRepositoryInterface].
func (c *cartRepository) RemoveFromCart(ctx context.Context, userId int64, productId int64) error {
	var (
		key = fmt.Sprintf("user:id:%d:cart:%d", userId, productId)
	)

	if err := c.redisClient.Del(ctx, key).Err(); err != nil {
		c.logger.Errorf("[CartRepository] RemoveFromCart: %v", err)
		return err
	}

	return nil
}

// RemoveAllFromCart implements [CartRepositoryInterface].
func (c *cartRepository) RemoveAllFromCart(ctx context.Context, userId int64) error {
	var (
		matchPattern = fmt.Sprintf("user:id:%d:cart:*", userId)
		cursor       uint64
		hasKeys      bool
	)

	// 1. Inisialisasi Pipeline
	pipe := c.redisClient.Pipeline()

	// 2. Scan semua key dan masukkan ke pipeline
	for {
		keys, nextCursor, err := c.redisClient.Scan(ctx, cursor, matchPattern, 100).Result()
		if err != nil {
			c.logger.Errorf("[CartRepository] RemoveAllFromCart: %v", err)
			return err
		}

		if len(keys) > 0 {
			hasKeys = true
			// Masukkan perintah DEL untuk setiap key yang ditemukan ke dalam pipeline
			for _, key := range keys {
				pipe.Del(ctx, key)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	// 3. Jika ada key yang dimasukkan ke pipeline, eksekusi semuanya sekaligus
	if hasKeys {
		_, err := pipe.Exec(ctx)
		if err != nil {
			c.logger.Errorf("[CartRepository] RemoveAllFromCart: %v", err)
			return err
		}
	}

	return nil
}
