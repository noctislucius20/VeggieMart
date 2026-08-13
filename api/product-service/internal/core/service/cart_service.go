package service

import (
	"context"
	"product-service/internal/adapter/repository"
	"product-service/internal/core/domain/entity"

	"github.com/labstack/gommon/log"
)

type CartServiceInterface interface {
	AddToCart(ctx context.Context, userId int64, req entity.CartItem) error
	GetCart(ctx context.Context, userId int64) ([]entity.CartItem, error)
	RemoveFromCart(ctx context.Context, userId int64, productId int64) error
	RemoveAllFromCart(ctx context.Context, userId int64) error
}

type cartService struct {
	cartRepo       repository.CartRepositoryInterface
	productService ProductServiceInterface
	logger         *log.Logger
}

func NewCartService(cartRepo repository.CartRepositoryInterface, productService ProductServiceInterface, logger *log.Logger) CartServiceInterface {
	return &cartService{
		cartRepo:       cartRepo,
		productService: productService,
		logger:         logger,
	}
}

// AddToCart implements [CartServiceInterface].
func (c *cartService) AddToCart(ctx context.Context, userId int64, req entity.CartItem) error {
	return c.cartRepo.AddToCart(ctx, userId, req)
}

// GetCart implements [CartServiceInterface].
func (c *cartService) GetCart(ctx context.Context, userId int64) ([]entity.CartItem, error) {
	var (
		productIds []int64
		productMap = make(map[int64]entity.ProductEntity)
	)

	cartEntities, err := c.cartRepo.GetCart(ctx, userId)
	if err != nil {
		c.logger.Errorf("[CartService] GetCart: %v", err)
		return nil, err
	}

	if len(cartEntities) == 0 {
		return nil, nil
	}

	for _, c := range cartEntities {
		productIds = append(productIds, c.ProductID)
	}

	productEntities, err := c.productService.GetBatchProducts(ctx, productIds)
	if err != nil {
		c.logger.Errorf("[CartService] GetCart: %v", err)
		return nil, err
	}

	for _, p := range productEntities {
		productMap[p.ID] = p
	}

	for i, c := range cartEntities {
		if productDetail, ok := productMap[c.ProductID]; ok {
			cartEntities[i].ProductDetail = productDetail
		}
	}

	return cartEntities, nil
}

// RemoveFromCart implements [CartServiceInterface].
func (c *cartService) RemoveFromCart(ctx context.Context, userId int64, productId int64) error {
	return c.cartRepo.RemoveFromCart(ctx, userId, productId)
}

// RemoveAllFromCart implements [CartServiceInterface].
func (c *cartService) RemoveAllFromCart(ctx context.Context, userId int64) error {
	return c.cartRepo.RemoveAllFromCart(ctx, userId)
}
