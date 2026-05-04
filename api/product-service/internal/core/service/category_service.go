package service

import (
	"context"
	"product-service/internal/adapter/repository"
	"product-service/internal/adapter/repository/cache"
	"product-service/internal/core/domain/entity"
	"product-service/internal/core/service/transaction"
	"product-service/utils/conv"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
)

type CategoryServiceInterface interface {
	GetAllCategories(ctx context.Context, query entity.QueryStringEntity) ([]entity.CategoryEntity, int64, int64, error)
	GetBatchCategories(ctx context.Context, categoryIds []int64) (map[int64]entity.CategoryEntity, error)
	GetCategoryById(ctx context.Context, id int64) (*entity.CategoryEntity, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*entity.CategoryEntity, error)
	CreateCategory(ctx context.Context, req entity.CategoryEntity) (string, int64, error)
	UpdateCategory(ctx context.Context, req entity.CategoryEntity) error
	DeleteCategory(ctx context.Context, categoryId int64) error

	GetAllCategoriesPublished(ctx context.Context) ([]entity.CategoryEntity, error)
}

type categoryService struct {
	repo          repository.CategoryRepositoryInterface
	redisClient   *redis.Client
	cacheCategory cache.CategoryCacheInterface
	txManager     transaction.TransactionManager
	logger        *log.Logger
}

func NewCategoryService(repo repository.CategoryRepositoryInterface, redisClient *redis.Client, cacheCategory cache.CategoryCacheInterface, txManager transaction.TransactionManager, logger *log.Logger) CategoryServiceInterface {
	return &categoryService{
		repo:          repo,
		redisClient:   redisClient,
		cacheCategory: cacheCategory,
		txManager:     txManager,
		logger:        logger,
	}
}

// GetBatchCategories implements [CategoryServiceInterface].
func (c *categoryService) GetBatchCategories(ctx context.Context, categoryIds []int64) (map[int64]entity.CategoryEntity, error) {
	var (
		categoriesMap = make(map[int64]entity.CategoryEntity)
	)

	if err := c.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		categoryEntities, err := c.repo.GetBatchCategories(txCtx, categoryIds)
		if err != nil {
			return err
		}

		for _, c := range categoryEntities {
			categoriesMap[c.ID] = c
		}

		return nil
	}); err != nil {
		c.logger.Errorf("[CategoryService-1] GetBatchCategories: %v", err)
		return nil, err
	}

	return categoriesMap, nil
}

// GetCategoryById implements [CategoryServiceInterface].
func (c *categoryService) GetCategoryById(ctx context.Context, id int64) (*entity.CategoryEntity, error) {
	var (
		category *entity.CategoryEntity
	)

	if err := c.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		categoryEntity, err := c.cacheCategory.GetCategoryById(txCtx, id)
		if err != nil {
			return err
		}

		category = categoryEntity

		return nil
	}); err != nil {
		c.logger.Errorf("[CategoryService-1] GetCategoryById: %v", err)
		return nil, err
	}

	return category, nil
}

// GetCategoryBySlug implements [CategoryServiceInterface].
func (c *categoryService) GetCategoryBySlug(ctx context.Context, slug string) (*entity.CategoryEntity, error) {
	var (
		category *entity.CategoryEntity
	)

	if err := c.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		categoryEntity, err := c.cacheCategory.GetCategoryBySlug(txCtx, slug)
		if err != nil {
			return err
		}

		category = categoryEntity

		return nil
	}); err != nil {
		c.logger.Errorf("[CategoryService-1] GetCategoryBySlug: %v", err)
		return nil, err
	}

	return category, nil
}

// GetAllCategoriesPublished implements [CategoryServiceInterface].
func (c *categoryService) GetAllCategoriesPublished(ctx context.Context) ([]entity.CategoryEntity, error) {
	var categories []entity.CategoryEntity

	if err := c.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		categoryEntities, err := c.repo.GetAllPublishedCategories(txCtx)
		if err != nil {
			return err
		}

		categories = categoryEntities

		return nil
	}); err != nil {
		c.logger.Errorf("[CategoryService-1] GetAllCategoriesPublished: %v", err)
		return nil, err
	}

	return categories, nil
}

// CreateCategory implements CategoryServiceInterface.
func (c *categoryService) CreateCategory(ctx context.Context, req entity.CategoryEntity) (string, int64, error) {
	var categoryId int64

	if err := c.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		req.Slug = conv.GenerateSlug(req.Name)

		categoryIdCreated, err := c.repo.CreateCategory(txCtx, req)
		if err != nil {
			return err
		}

		if err := c.cacheCategory.DeleteCategoryCache(ctx, categoryIdCreated); err != nil {
			return err
		}

		categoryId = categoryIdCreated

		return nil
	}); err != nil {
		c.logger.Errorf("[CategoryService-1] CreateCategory: %v", err)
		return "", 0, err
	}

	return req.Slug, categoryId, nil
}

// DeleteCategory implements CategoryServiceInterface.
func (c *categoryService) DeleteCategory(ctx context.Context, categoryId int64) error {
	if err := c.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := c.repo.DeleteCategory(txCtx, categoryId); err != nil {
			return err
		}

		if err := c.cacheCategory.DeleteCategoryCache(ctx, categoryId); err != nil {
			return err
		}

		return nil
	}); err != nil {
		c.logger.Errorf("[CategoryService-1] DeleteCategory: %v", err)
		return err
	}

	return nil
}

// GetAllCategories implements CategoryServiceInterface.
func (c *categoryService) GetAllCategories(ctx context.Context, query entity.QueryStringEntity) ([]entity.CategoryEntity, int64, int64, error) {
	var (
		categories []entity.CategoryEntity
		countData  int64
		totalPages int64
	)

	if err := c.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		categoryEntities, count, pages, err := c.repo.GetAllCategories(txCtx, query)
		if err != nil {
			return err
		}

		categories, countData, totalPages = categoryEntities, count, pages

		return nil
	}); err != nil {
		c.logger.Errorf("[CategoryService-1] GetAllCategories: %v", err)
		return nil, 0, 0, err
	}

	return categories, countData, totalPages, nil
}

// UpdateCategory implements CategoryServiceInterface.
func (c *categoryService) UpdateCategory(ctx context.Context, req entity.CategoryEntity) error {
	if err := c.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		req.Slug = conv.GenerateSlug(req.Name)

		if err := c.repo.UpdateCategory(txCtx, req); err != nil {
			return err
		}

		if err := c.cacheCategory.DeleteCategoryCache(ctx, req.ID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		c.logger.Errorf("[CategoryService-1] UpdateCategory: %v", err)
		return err
	}

	return nil
}
