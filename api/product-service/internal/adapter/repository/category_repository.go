package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"product-service/internal/core/domain/entity"
	"product-service/internal/core/domain/model"
	"product-service/utils"
	"strings"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type CategoryRepositoryInterface interface {
	GetAllCategories(ctx context.Context, query entity.QueryStringEntity) ([]entity.CategoryEntity, int64, int64, error)
	GetCategoryByIdOrSlug(ctx context.Context, id int64, slug string) (*entity.CategoryEntity, error)
	CreateCategory(ctx context.Context, req entity.CategoryEntity) (int64, error)
	UpdateCategory(ctx context.Context, req entity.CategoryEntity) (string, error)
	DeleteCategory(ctx context.Context, categoryId int64) (string, error)

	GetAllPublishedCategories(ctx context.Context) ([]entity.CategoryEntity, error)

	getDB(ctx context.Context) *gorm.DB
}

type categoryRepository struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewCategoryRepository(db *gorm.DB, logger *log.Logger) CategoryRepositoryInterface {
	return &categoryRepository{db: db, logger: logger}
}

// getDB implements [UserRepositoryInterface].
func (c *categoryRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}

	return c.db
}

// GetAllPublishedCategories implements [CategoryRepositoryInterface].
func (c *categoryRepository) GetAllPublishedCategories(ctx context.Context) ([]entity.CategoryEntity, error) {
	var (
		db              = c.getDB(ctx)
		modelCategories []model.Category
	)

	if err := db.WithContext(ctx).
		Select("id", "parent_id", "name", "icon", "slug").
		Where("status = ?", true).
		Order("COALESCE(parent_id, id), parent_id IS NOT NULL, id").
		Find(&modelCategories).Error; err != nil {
		c.logger.Errorf("[CategoryRepository-1] GetAllPublishedCategories: %v", err)
		return nil, err
	}

	entities := []entity.CategoryEntity{}
	for _, val := range modelCategories {
		entities = append(entities, entity.CategoryEntity{
			ID:       val.ID,
			ParentID: val.ParentID,
			Name:     val.Name,
			Icon:     val.Icon,
			Slug:     val.Slug,
		})
	}

	return entities, nil
}

// DeleteCategory implements CategoryRepositoryInterface.
func (c *categoryRepository) DeleteCategory(ctx context.Context, categoryId int64) (string, error) {
	var (
		db            = c.getDB(ctx)
		modelCategory = model.Category{
			ID: categoryId,
		}
		categoryDeleteDTO model.CategoryDTO
	)

	if err := db.WithContext(ctx).
		Select("categories.id AS category_id",
			"products.id AS product_id",
			"categories.slug AS slug",
		).
		Joins("JOIN products ON categories.id = products.category_id").
		Where("categories.id = ?", categoryId).
		First(&categoryDeleteDTO).Error; err != nil {
		c.logger.Errorf("[CategoryRepository-1] DeleteCategory: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := errors.New(utils.DATA_NOT_FOUND)
			return "", err
		}
		return "", err
	}

	if categoryDeleteDTO.ProductID > 0 {
		err := errors.New(utils.DATA_STILL_IN_USED)
		c.logger.Errorf("[CategoryRepository-2] DeleteCategory: %v", err)
		return "", err
	}

	if err := db.WithContext(ctx).Delete(&modelCategory).Error; err != nil {
		c.logger.Errorf("[CategoryRepository-3] DeleteCategory: %v", err)
		return "", err
	}

	return categoryDeleteDTO.Slug, nil
}

// UpdateCategory implements CategoryRepositoryInterface.
func (c *categoryRepository) UpdateCategory(ctx context.Context, req entity.CategoryEntity) (string, error) {
	var (
		db                = c.getDB(ctx)
		modelCategory     model.Category
		categoryUpdateDTO model.CategoryDTO
	)

	if err := db.WithContext(ctx).
		Select("id AS category_id", "slug").
		Where("id = ?", req.ID).
		First(&categoryUpdateDTO).Error; err != nil {
		c.logger.Errorf("[CategoryRepository-1] UpdateCategory: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := errors.New(utils.DATA_NOT_FOUND)
			return "", err
		}
		return "", err
	}

	modelCategory = model.Category{
		ID:          req.ID,
		ParentID:    req.ParentID,
		Name:        req.Name,
		Icon:        req.Icon,
		Slug:        req.Slug,
		Description: req.Description,
		Status:      strings.ToLower(req.Status) == "published",
	}

	tx := db.WithContext(ctx).Updates(&modelCategory)
	if tx.Error != nil {
		c.logger.Errorf("[CategoryRepository-2] UpdateCategory: %v", tx.Error)
		if strings.Contains(tx.Error.Error(), "duplicate key") {
			err := errors.New(utils.DATA_ALREADY_EXISTS)
			return "", err
		}
		return "", tx.Error
	}

	return categoryUpdateDTO.Slug, nil
}

// CreateCategory implements CategoryRepositoryInterface.
func (c *categoryRepository) CreateCategory(ctx context.Context, req entity.CategoryEntity) (int64, error) {
	var (
		db            = c.getDB(ctx)
		modelCategory = model.Category{
			ParentID:    req.ParentID,
			Name:        req.Name,
			Icon:        req.Icon,
			Slug:        req.Slug,
			Description: req.Description,
			Status:      strings.ToLower(req.Status) == "published",
		}
	)

	if err := db.WithContext(ctx).Create(&modelCategory).Error; err != nil {
		c.logger.Errorf("[CategoryRepository-1] CreateCategory: %v", err)
		if strings.Contains(err.Error(), "duplicate key") {
			err := errors.New(utils.DATA_ALREADY_EXISTS)
			return 0, err
		}
		return 0, err
	}

	return modelCategory.ID, nil
}

// GetCategoryByIdOrSlug implements CategoryRepositoryInterface.
func (c *categoryRepository) GetCategoryByIdOrSlug(ctx context.Context, id int64, slug string) (*entity.CategoryEntity, error) {
	var (
		db            = c.getDB(ctx)
		modelCategory model.Category
		status        string
	)

	if err := db.WithContext(ctx).
		Omit("created_at", "updated_at", "deleted_at").
		First(&modelCategory, "id = ? OR slug = ?", id, slug).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New(utils.DATA_NOT_FOUND)
			c.logger.Errorf("[CategoryRepository-1] GetCategoryByIdOrSlug: %v", err)
			return nil, err
		}
		c.logger.Errorf("[CategoryRepository-2] GetCategoryByIdOrSlug: %v", err)
		return nil, err
	}

	if modelCategory.Status {
		status = "PUBLISHED"
	} else {
		status = "UNPUBLISHED"
	}

	return &entity.CategoryEntity{
		ID:          modelCategory.ID,
		ParentID:    modelCategory.ParentID,
		Name:        modelCategory.Name,
		Icon:        modelCategory.Icon,
		Status:      status,
		Slug:        modelCategory.Slug,
		Description: modelCategory.Description,
	}, nil
}

// GetAllCategories implements CategoryRepositoryInterface.
func (c *categoryRepository) GetAllCategories(ctx context.Context, query entity.QueryStringEntity) ([]entity.CategoryEntity, int64, int64, error) {
	var (
		db              = c.getDB(ctx)
		modelCategories []model.Category
		countData       int64
	)

	orderSort := fmt.Sprintf("%s %s", query.OrderBy, query.OrderType)
	offset := (query.Page - 1) * query.Limit

	sqlMain := db.WithContext(ctx).Preload("Products").
		Where("name ILIKE ? OR slug ILIKE ?", "%"+query.Search+"%", "%"+query.Search+"%")

	if err := sqlMain.Model(&modelCategories).
		Count(&countData).Error; err != nil {
		c.logger.Errorf("[CategoryRepository-1] GetAllCategories: %v", err)
		return nil, 0, 0, err
	}

	totalPage := int(math.Ceil(float64(countData) / float64(query.Limit)))
	if err := sqlMain.Order(orderSort).
		Limit(int(query.Limit)).
		Offset(int(offset)).
		Find(&modelCategories).Error; err != nil {
		c.logger.Errorf("[CategoryRepository-2] GetAllCategories: %v", err)
		return nil, 0, 0, err
	}

	entities := []entity.CategoryEntity{}
	for _, val := range modelCategories {
		productEntities := []entity.ProductEntity{}
		for _, prd := range val.Products {
			productEntities = append(productEntities, entity.ProductEntity{
				ID:       prd.ID,
				ParentID: prd.ParentID,
				Name:     prd.Name,
				Image:    prd.Image,
			})
		}

		var status string
		if val.Status {
			status = "PUBLISHED"
		} else {
			status = "UNPUBLISHED"
		}

		entities = append(entities, entity.CategoryEntity{
			ID:          val.ID,
			ParentID:    val.ParentID,
			Name:        val.Name,
			Icon:        val.Icon,
			Status:      status,
			Slug:        val.Slug,
			Description: val.Description,
			Products:    productEntities,
		})
	}

	return entities, countData, int64(totalPage), nil
}
