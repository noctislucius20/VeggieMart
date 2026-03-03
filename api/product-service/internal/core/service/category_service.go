package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"product-service/internal/adapter/repository"
	"product-service/internal/core/domain/entity"
	"product-service/internal/core/service/transaction"
	"product-service/utils"
	"product-service/utils/conv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
)

type CategoryServiceInterface interface {
	GetAllCategories(ctx context.Context, query entity.QueryStringEntity) ([]entity.CategoryEntity, int64, int64, error)
	GetCategoryByIdOrSlug(ctx context.Context, id int64, slug string) (*entity.CategoryEntity, error)
	CreateCategory(ctx context.Context, req entity.CategoryEntity) (string, int64, error)
	UpdateCategory(ctx context.Context, req entity.CategoryEntity) error
	DeleteCategory(ctx context.Context, categoryId int64) error

	GetAllCategoriesPublished(ctx context.Context) ([]entity.CategoryEntity, error)
}

type categoryService struct {
	repo        repository.CategoryRepositoryInterface
	redisClient *redis.Client
	txManager   transaction.TransactionManager
	logger      *log.Logger
}

func NewCategoryService(repo repository.CategoryRepositoryInterface, redisClient *redis.Client, txManager transaction.TransactionManager, logger *log.Logger) CategoryServiceInterface {
	return &categoryService{
		repo:        repo,
		redisClient: redisClient,
		txManager:   txManager,
		logger:      logger,
	}
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

		id, err := c.repo.CreateCategory(txCtx, req)
		if err != nil {
			return err
		}

		categoryId = id

		return nil
	}); err != nil {
		c.logger.Errorf("[CategoryService-1] CreateCategory: %v", err)
		return "", 0, err
	}

	// redis delete key
	key := fmt.Sprintf("category:%s", req.Slug)
	if err := c.redisClient.Del(ctx, key); err != nil {
		c.logger.Errorf("[CreateCategory-2] CreateCategory: %v", err)
	}

	return req.Slug, categoryId, nil
}

// DeleteCategory implements CategoryServiceInterface.
func (c *categoryService) DeleteCategory(ctx context.Context, categoryId int64) error {
	var slug string

	if err := c.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		slugDeleted, err := c.repo.DeleteCategory(txCtx, categoryId)
		if err != nil {
			return err
		}

		slug = slugDeleted

		return nil
	}); err != nil {
		c.logger.Errorf("[CategoryService-1] DeleteCategory: %v", err)
		return err
	}

	// redis delete key
	key := fmt.Sprintf("category:%s", slug)
	if err := c.redisClient.Del(ctx, key); err != nil {
		c.logger.Errorf("[CreateCategory-2] CreateCategory: %v", err)
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

// GetCategoryByIdOrSlug implements CategoryServiceInterface.
func (c *categoryService) GetCategoryByIdOrSlug(ctx context.Context, id int64, slug string) (*entity.CategoryEntity, error) {
	var (
		category *entity.CategoryEntity
		key      = fmt.Sprintf("category:%s", slug)
	)

	// Check redis if data exists.
	val, err := c.redisClient.Get(ctx, key).Result()
	if err == nil {
		// if key exists but value null, return data not found error
		if val == "null" {
			err := errors.New(utils.DATA_NOT_FOUND)
			return nil, err
		}

		json.Unmarshal([]byte(val), &category)
		return category, nil
	}

	// Query DB
	if err := c.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		categoryEntity, err := c.repo.GetCategoryByIdOrSlug(txCtx, id, slug)
		if err != nil {
			return err
		}

		category = categoryEntity

		return nil

	}); err != nil {
		// Save to redis (create key with null value if data not found)
		if err.Error() == utils.DATA_NOT_FOUND {
			if err := c.redisClient.Set(ctx, key, "null", 1*time.Minute); err != nil {
				c.logger.Errorf("[CategoryService-1] GetCategoryByIdOrSlug: %v", err)
			}
		}

		c.logger.Errorf("[CategoryService-2] GetCategoryByIdOrSlug: %v", err)
		return nil, err
	}

	// Save to redis
	jsonData, _ := json.Marshal(category)
	if err := c.redisClient.Set(ctx, key, jsonData, 1*time.Hour).Err(); err != nil {
		c.logger.Errorf("[CategoryService-3] GetCategoryByIdOrSlug: %v", err)
	}

	return category, nil
}

// UpdateCategory implements CategoryServiceInterface.
func (c *categoryService) UpdateCategory(ctx context.Context, req entity.CategoryEntity) error {
	var slug string

	if err := c.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		req.Slug = conv.GenerateSlug(req.Name)

		slugOld, err := c.repo.UpdateCategory(txCtx, req)
		if err != nil {
			return err
		}

		slug = slugOld

		return nil
	}); err != nil {
		c.logger.Errorf("[CategoryService-1] UpdateCategory: %v", err)
		return err
	}

	// redis delete key
	key := fmt.Sprintf("category:%s", slug)
	if err := c.redisClient.Del(ctx, key); err != nil {
		c.logger.Errorf("[CategoryService-2] UpdateCategory: %v", err)
	}

	return nil
}
