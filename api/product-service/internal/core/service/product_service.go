package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"product-service/config"
	"product-service/internal/adapter/repository"
	"product-service/internal/core/domain/entity"
	"product-service/internal/core/service/transaction"
	"product-service/utils"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
)

type ProductServiceInterface interface {
	GetAllProducts(ctx context.Context, query entity.QueryStringProduct) ([]entity.ProductEntity, int64, int64, error)
	GetBatchProducts(ctx context.Context, productIds []int64) ([]entity.ProductEntity, error)
	GetProductById(ctx context.Context, productId int64) (*entity.ProductEntity, error)
	CreateProduct(ctx context.Context, req entity.ProductEntity) (int64, error)
	UpdateProduct(ctx context.Context, req entity.ProductEntity) error
	DeleteProduct(ctx context.Context, productId int64) error
}

type productService struct {
	repo            repository.ProductRepositoryInterface
	redisClient     *redis.Client
	categoryService CategoryServiceInterface
	repoOutbox      repository.OutboxEventInterface
	txManager       transaction.TransactionManager
	logger          *log.Logger
	cfg             *config.Config
}

func NewProductService(cfg *config.Config, repo repository.ProductRepositoryInterface, redisClient *redis.Client, txManager transaction.TransactionManager, categoryService CategoryServiceInterface, repoOutbox repository.OutboxEventInterface, logger *log.Logger) ProductServiceInterface {
	return &productService{
		cfg:             cfg,
		repo:            repo,
		redisClient:     redisClient,
		txManager:       txManager,
		categoryService: categoryService,
		repoOutbox:      repoOutbox,
		logger:          logger,
	}
}

// GetBatchProducts implements [ProductServiceInterface].
func (p *productService) GetBatchProducts(ctx context.Context, productIds []int64) ([]entity.ProductEntity, error) {
	var products []entity.ProductEntity

	if err := p.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		productEntities, err := p.repo.GetBatchProducts(txCtx, productIds)
		if err != nil {
			return err
		}

		products = productEntities

		return nil
	}); err != nil {
		p.logger.Errorf("[ProductService-1] GetBatchProducts: %v", err)
		return nil, err
	}

	return products, nil
}

// CreateProduct implements ProductServiceInterface.
func (p *productService) CreateProduct(ctx context.Context, req entity.ProductEntity) (int64, error) {
	var (
		productId   int64
		publishName = p.cfg.PublisherName.ProductCreate
	)

	if err := p.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		categoryEntity, err := p.categoryService.GetCategoryByIdOrSlug(txCtx, 0, req.CategorySlug)
		if err != nil {
			return err
		}

		req.CategoryID = categoryEntity.ID

		id, err := p.repo.CreateProduct(txCtx, req)
		if err != nil {
			return err
		}

		productEntity, err := p.repo.GetProductById(txCtx, id)
		if err != nil {
			return err
		}

		if err := p.repoOutbox.CreateEvent(txCtx, publishName, productEntity, &productEntity.ID); err != nil {
			return err
		}

		productId = productEntity.ID

		return nil
	}); err != nil {
		p.logger.Errorf("[ProductService-1] CreateProduct: %v", err)
		return 0, err
	}

	// redis delete key
	key := fmt.Sprintf("product:%d", productId)
	if err := p.redisClient.Del(ctx, key); err != nil {
		p.logger.Errorf("[ProductService-2] CreateProduct: %v", err)
	}

	return productId, nil
}

// DeleteProduct implements ProductServiceInterface.
func (p *productService) DeleteProduct(ctx context.Context, productId int64) error {
	var publishName = p.cfg.PublisherName.ProductDelete

	if err := p.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := p.repo.DeleteProduct(txCtx, productId); err != nil {
			return err
		}

		productDeletePayload := map[string]any{
			"product_id": productId,
		}
		if err := p.repoOutbox.CreateEvent(ctx, publishName, productDeletePayload, &productId); err != nil {
			return err
		}

		return nil
	}); err != nil {
		p.logger.Errorf("[ProductService-1] DeleteProduct: %v", err)
		return err
	}

	// redis delete key
	key := fmt.Sprintf("product:%d", productId)
	if err := p.redisClient.Del(ctx, key); err != nil {
		p.logger.Errorf("[ProductService-2] DeleteProduct: %v", err)
	}

	return nil
}

// GetAllProducts implements ProductServiceInterface.
func (p *productService) GetAllProducts(ctx context.Context, query entity.QueryStringProduct) ([]entity.ProductEntity, int64, int64, error) {
	products, countData, totalPages, err := p.repo.SearchProducts(ctx, query)
	if err == nil {
		return products, countData, totalPages, nil
	}

	if err := p.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		productEntities, count, pages, err := p.repo.GetAllProducts(txCtx, query)
		if err != nil {
			return err
		}

		products, countData, totalPages = productEntities, count, pages

		return nil
	}); err != nil {
		p.logger.Errorf("[ProductService-1] GetAllProducts: %v", err)
		return nil, 0, 0, err
	}

	return products, countData, totalPages, nil
}

// GetProductById implements ProductServiceInterface.
func (p *productService) GetProductById(ctx context.Context, productId int64) (*entity.ProductEntity, error) {
	var (
		product *entity.ProductEntity
		key     = fmt.Sprintf("product:%d", productId)
	)

	// Check redis if data exists.
	val, err := p.redisClient.Get(ctx, key).Result()
	if err == nil {
		// if key exists but value null, return data not found error
		if val == "null" {
			err := errors.New(utils.DATA_NOT_FOUND)
			p.logger.Errorf("[ProductService-1] GetProductById: %v", err)
			return nil, err
		}

		json.Unmarshal([]byte(val), &product)
		return product, nil
	}

	// Query DB
	if err := p.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		productEntity, err := p.repo.GetProductById(txCtx, productId)
		if err != nil {
			return err
		}

		product = productEntity

		return nil
	}); err != nil {
		// Save to redis (create key with null value if data not found)
		if err.Error() == utils.DATA_NOT_FOUND {
			if err := p.redisClient.Set(ctx, key, "null", 1*time.Minute); err != nil {
				p.logger.Errorf("[ProductService-2] GetProductById: %v", err)
			}
		}

		p.logger.Errorf("[ProductService-3] GetProductById: %v", err)
		return nil, err
	}

	// Save to redis
	jsonData, _ := json.Marshal(product)
	if err := p.redisClient.Set(ctx, key, jsonData, 1*time.Hour).Err(); err != nil {
		p.logger.Errorf("[ProductService-4] GetProductById: %v", err)
	}

	return product, nil
}

// UpdateProduct implements ProductServiceInterface.
func (p *productService) UpdateProduct(ctx context.Context, req entity.ProductEntity) error {
	publishName := p.cfg.PublisherName.ProductUpdate

	if err := p.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		categoryEntity, err := p.categoryService.GetCategoryByIdOrSlug(txCtx, 0, req.CategorySlug)
		if err != nil {
			return err
		}

		req.CategoryID = categoryEntity.ID

		if err := p.repo.UpdateProduct(ctx, req); err != nil {
			return err
		}

		productEntity, err := p.repo.GetProductById(txCtx, req.ID)
		if err != nil {
			return err
		}

		if err := p.repoOutbox.CreateEvent(ctx, publishName, productEntity, &productEntity.ID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		p.logger.Errorf("[ProductService-1] UpdateProduct: %v", err)
		return err
	}

	return nil
}
