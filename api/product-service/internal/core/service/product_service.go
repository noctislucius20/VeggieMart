package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"product-service/config"
	"product-service/internal/adapter/repository"
	"product-service/internal/adapter/repository/cache"
	"product-service/internal/core/domain/entity"
	"product-service/internal/core/service/transaction"
	"product-service/utils"

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
	UpdateStockProduct(ctx context.Context, products []entity.ProductUpdateStockEntity) error

	getAllProductsCategory(ctx context.Context, products []entity.ProductEntity) error
}

type productService struct {
	repo            repository.ProductRepositoryInterface
	redisClient     *redis.Client
	categoryService CategoryServiceInterface
	repoOutbox      repository.OutboxEventInterface
	repoElastic     repository.ElasticRepositoryInterface
	cacheProduct    cache.ProductCacheInterface
	txManager       transaction.TransactionManager
	logger          *log.Logger
	cfg             *config.Config
}

func NewProductService(cfg *config.Config, repo repository.ProductRepositoryInterface, redisClient *redis.Client, cacheProduct cache.ProductCacheInterface, txManager transaction.TransactionManager, categoryService CategoryServiceInterface, repoOutbox repository.OutboxEventInterface, repoElastic repository.ElasticRepositoryInterface, logger *log.Logger) ProductServiceInterface {
	return &productService{
		cfg:             cfg,
		repo:            repo,
		redisClient:     redisClient,
		cacheProduct:    cacheProduct,
		txManager:       txManager,
		categoryService: categoryService,
		repoOutbox:      repoOutbox,
		repoElastic:     repoElastic,
		logger:          logger,
	}
}

// UpdateStockProduct implements [ProductServiceInterface].
func (p *productService) UpdateStockProduct(ctx context.Context, products []entity.ProductUpdateStockEntity) error {
	if err := p.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := p.repo.UpdateStockProduct(txCtx, products); err != nil {
			return err
		}

		return nil
	}); err != nil {
		p.logger.Errorf("[ProductService-1] UpdateStockProduct: %v", err)
		return err
	}

	return nil
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
		productId            int64
		publishProductCreate = p.cfg.PublisherName.ProductCreate
		outboxEventEntities  []entity.OutboxEventEntity
	)

	if err := p.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		categoryEntity, err := p.categoryService.GetCategoryBySlug(txCtx, req.CategorySlug)
		if err != nil {
			if err.Error() == utils.DATA_NOT_FOUND {
				err := errors.New(utils.RELATION_DATA_NOT_FOUND)
				return err
			}
			return err
		}

		req.CategoryID = categoryEntity.ID

		productIdCreated, err := p.repo.CreateProduct(txCtx, req)
		if err != nil {
			return err
		}

		if err := p.cacheProduct.DeleteProductCache(txCtx, productIdCreated); err != nil {
			return err
		}

		productEntity, err := p.cacheProduct.GetProductById(txCtx, productIdCreated)
		if err != nil {
			return err
		}

		jsonProduct, _ := json.Marshal(productEntity)

		outboxEventEntities = append(outboxEventEntities, entity.OutboxEventEntity{
			EventType:   publishProductCreate,
			Payload:     string(jsonProduct),
			AggregateID: fmt.Sprintf("%d", productIdCreated),
		})

		if err := p.repoOutbox.CreateBatchEvents(txCtx, outboxEventEntities); err != nil {
			return err
		}

		productId = productEntity.ID

		return nil
	}); err != nil {
		p.logger.Errorf("[ProductService-1] CreateProduct: %v", err)
		return 0, err
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
			"id": productId,
		}
		if err := p.repoOutbox.CreateEvent(txCtx, publishName, productDeletePayload, &productId); err != nil {
			return err
		}

		if err := p.cacheProduct.DeleteProductCache(txCtx, productId); err != nil {
			return err
		}

		return nil
	}); err != nil {
		p.logger.Errorf("[ProductService-1] DeleteProduct: %v", err)
		return err
	}

	return nil
}

// GetAllProducts implements ProductServiceInterface.
func (p *productService) GetAllProducts(ctx context.Context, query entity.QueryStringProduct) ([]entity.ProductEntity, int64, int64, error) {
	products, countData, totalPages, err := p.repoElastic.SearchProductElastic(ctx, query)
	if err == nil {
		if err := p.getAllProductsCategory(ctx, products); err != nil {
			if err.Error() == utils.DATA_NOT_FOUND {
				err := errors.New(utils.RELATION_DATA_NOT_FOUND)
				return nil, 0, 0, err
			}
		}

		return products, countData, totalPages, nil
	}

	if err := p.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		productEntities, count, pages, err := p.repo.GetAllProducts(txCtx, query)
		if err != nil {
			return err
		}

		if len(productEntities) == 0 {
			return nil
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
		product entity.ProductEntity
	)

	if err := p.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		productEntity, err := p.cacheProduct.GetProductById(txCtx, productId)
		if err != nil {
			return err
		}

		categoryEntity, err := p.categoryService.GetCategoryById(txCtx, productEntity.CategoryID)
		if err != nil {
			if err.Error() == utils.DATA_NOT_FOUND {
				err := errors.New(utils.RELATION_DATA_NOT_FOUND)
				return err
			}
			return err
		}

		productEntity.CategoryName = categoryEntity.Name
		productEntity.CategorySlug = categoryEntity.Slug

		product = *productEntity

		return nil
	}); err != nil {
		p.logger.Errorf("[ProductService-1] GetProductById: %v", err)
		return nil, err
	}

	return &product, nil
}

// UpdateProduct implements ProductServiceInterface.
func (p *productService) UpdateProduct(ctx context.Context, req entity.ProductEntity) error {
	publishName := p.cfg.PublisherName.ProductUpdate

	if err := p.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		categoryEntity, err := p.categoryService.GetCategoryBySlug(txCtx, req.CategorySlug)
		if err != nil {
			if err.Error() == utils.DATA_NOT_FOUND {
				err := errors.New(utils.RELATION_DATA_NOT_FOUND)
				return err
			}
			return err
		}

		req.CategoryID = categoryEntity.ID

		if err := p.repo.UpdateProduct(txCtx, req); err != nil {
			return err
		}

		if err := p.cacheProduct.DeleteProductCache(txCtx, req.ID); err != nil {
			return err
		}

		productEntity, err := p.cacheProduct.GetProductById(txCtx, req.ID)
		if err != nil {
			return err
		}

		if err := p.repoOutbox.CreateEvent(txCtx, publishName, productEntity, &productEntity.ID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		p.logger.Errorf("[ProductService-1] UpdateProduct: %v", err)
		return err
	}

	return nil
}

// getAllProductsCategory implements [ProductServiceInterface].
func (p *productService) getAllProductsCategory(ctx context.Context, products []entity.ProductEntity) error {
	categoryIds := map[int64]struct{}{}
	for _, product := range products {
		categoryIds[product.CategoryID] = struct{}{}
	}

	categoryList := make([]int64, 0, len(categoryIds))
	for id := range categoryIds {
		categoryList = append(categoryList, id)
	}

	resultCategories, err := p.categoryService.GetBatchCategories(ctx, categoryList)
	if err != nil {
		return err
	}

	for pIdx, product := range products {
		if c, ok := resultCategories[product.CategoryID]; ok {
			products[pIdx].CategoryName = c.Name
			products[pIdx].CategorySlug = c.Slug
		}
	}

	return nil
}
