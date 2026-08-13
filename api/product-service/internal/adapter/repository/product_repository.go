package repository

import (
	"context"
	"fmt"
	"math"
	"product-service/internal/core/domain/entity"
	"product-service/internal/core/domain/model"
	"product-service/utils"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type ProductRepositoryInterface interface {
	GetAllProducts(ctx context.Context, query entity.QueryStringProduct) ([]entity.ProductEntity, int64, int64, error)
	GetBatchProducts(ctx context.Context, productIds []int64) ([]entity.ProductEntity, error)
	GetProductById(ctx context.Context, productId int64) (*entity.ProductEntity, error)
	CreateProduct(ctx context.Context, req entity.ProductEntity) (int64, error)
	UpdateProduct(ctx context.Context, req entity.ProductEntity) error
	DeleteProduct(ctx context.Context, productId int64) error
	UpdateStockProduct(ctx context.Context, products []entity.ProductUpdateStockEntity) error

	getDB(ctx context.Context) *gorm.DB
}

type productRepository struct {
	db       *gorm.DB
	esClient *elasticsearch.Client
	logger   *log.Logger
}

func NewProductRepository(db *gorm.DB, esClient *elasticsearch.Client, logger *log.Logger) ProductRepositoryInterface {
	return &productRepository{db: db, esClient: esClient, logger: logger}
}

// getDB implements [ProductRepositoryInterface].
func (p *productRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}

	return p.db
}

// UpdateStockProduct implements [ProductRepositoryInterface].
func (p *productRepository) UpdateStockProduct(ctx context.Context, products []entity.ProductUpdateStockEntity) error {
	var (
		db = p.getDB(ctx)
	)

	for _, product := range products {
		tx := db.WithContext(ctx).
			Model(&model.Product{}).
			Where("id = ? AND stock >= ?", product.ProductID, product.Quantity).
			Update("stock", gorm.Expr("stock - ?", product.Quantity))

		if tx.Error != nil {
			p.logger.Errorf("[ProductRepository] UpdateProduct: %v", tx.Error)
			return tx.Error
		}

		if tx.RowsAffected == 0 {
			err := utils.ErrStockUnavailable
			p.logger.Errorf("[ProductRepository] UpdateProduct: %v", err)
			return err
		}
	}

	return nil
}

// GetBatchProducts implements [ProductRepositoryInterface].
func (p *productRepository) GetBatchProducts(ctx context.Context, productIds []int64) ([]entity.ProductEntity, error) {
	var (
		db            = p.getDB(ctx)
		modelProducts []model.Product
	)

	chunkSize := 150

	for i := 0; i < len(productIds); i += chunkSize {
		end := min(i+chunkSize, len(productIds))

		batchProducts := []model.Product{}
		if err := db.WithContext(ctx).
			Select("id", "image", "name", "sale_price", "weight", "unit", "regular_price").
			Where("id IN ?", productIds[i:end]).
			Find(&batchProducts).Error; err != nil {
			p.logger.Errorf("[ProductRepository] GetBatchProducts: %v", err)
			return nil, err
		}

		modelProducts = append(modelProducts, batchProducts...)
	}

	if len(modelProducts) == 0 {
		err := utils.ErrDataNotFound
		p.logger.Errorf("[ProductRepository] GetBatchProducts: %v", err)
		return nil, err
	}

	entities := []entity.ProductEntity{}
	for _, val := range modelProducts {
		entities = append(entities, entity.ProductEntity{
			ID:           val.ID,
			Name:         val.Name,
			RegularPrice: val.RegularPrice,
			Image:        val.Image,
			SalePrice:    val.SalePrice,
			Weight:       val.Weight,
			Unit:         val.Unit,
		})
	}

	return entities, nil
}

// DeleteProduct implements ProductRepositoryInterface.
func (p *productRepository) DeleteProduct(ctx context.Context, productId int64) error {
	var (
		db           = p.getDB(ctx)
		modelProduct model.Product
	)

	tx := db.WithContext(ctx).
		Where("id = ? OR parent_id = ?", productId, productId).
		Delete(&modelProduct)
	if tx.Error != nil {
		p.logger.Errorf("[ProductRepository] DeleteProduct: %v", tx.Error)
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		err := utils.ErrDataNotFound
		p.logger.Errorf("[ProductRepository] DeleteProduct: %v", err)
		return err
	}

	return nil
}

// UpdateProduct implements ProductRepositoryInterface.
func (p *productRepository) UpdateProduct(ctx context.Context, req entity.ProductEntity) error {
	var (
		db           = p.getDB(ctx)
		modelProduct = &model.Product{
			CategoryID:   req.CategoryID,
			Name:         req.Name,
			Image:        req.Image,
			Description:  req.Description,
			RegularPrice: req.RegularPrice,
			SalePrice:    req.SalePrice,
			Unit:         req.Unit,
			Weight:       req.Weight,
			Stock:        req.Stock,
			Variant:      req.Variant,
			Status:       req.Status,
		}
	)

	tx := db.WithContext(ctx).
		Where("id = ? AND parent_id IS NULL", req.ID).
		Updates(&modelProduct)
	if tx.Error != nil {
		p.logger.Errorf("[ProductRepository] UpdateProduct: %v", tx.Error)
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		err := utils.ErrDataNotFound
		p.logger.Errorf("[ProductRepository] UpdateProduct: %v", err)
		return err
	}

	modelProductChild := []model.Product{}
	if len(req.Childs) > 0 {
		if err := db.WithContext(ctx).
			Where("parent_id = ?", req.ID).
			Delete(&model.Product{}).Error; err != nil {
			p.logger.Errorf("[ProductRepository] UpdateProduct: %v", err)
			return err
		}

		for _, val := range req.Childs {
			modelProductChild = append(modelProductChild, model.Product{
				CategoryID:   req.CategoryID,
				ParentID:     &req.ID,
				Name:         req.Name,
				Image:        val.Image,
				Description:  req.Description,
				RegularPrice: val.RegularPrice,
				SalePrice:    val.SalePrice,
				Unit:         req.Unit,
				Weight:       val.Weight,
				Stock:        val.Stock,
				Variant:      req.Variant,
				Status:       req.Status,
			})
		}

		if err := db.WithContext(ctx).
			CreateInBatches(&modelProductChild, 50).Error; err != nil {
			p.logger.Errorf("[ProductRepository] UpdateProduct: %v", err)
			return err
		}
	}

	return nil

}

// CreateProduct implements ProductRepositoryInterface.
func (p *productRepository) CreateProduct(ctx context.Context, req entity.ProductEntity) (int64, error) {
	var (
		db           = p.getDB(ctx)
		modelProduct = &model.Product{
			CategoryID:   req.CategoryID,
			Name:         req.Name,
			Image:        req.Image,
			Description:  req.Description,
			RegularPrice: req.RegularPrice,
			SalePrice:    req.SalePrice,
			Unit:         req.Unit,
			Weight:       req.Weight,
			Stock:        req.Stock,
			Variant:      req.Variant,
			Status:       req.Status,
		}
	)

	if err := db.WithContext(ctx).Create(&modelProduct).Error; err != nil {
		p.logger.Errorf("[ProductRepository] CreateProduct: %v", err)
		return 0, err
	}

	modelProductChild := []model.Product{}
	if len(req.Childs) > 0 {
		for _, val := range req.Childs {
			modelProductChild = append(modelProductChild, model.Product{
				CategoryID:   req.CategoryID,
				ParentID:     &modelProduct.ID,
				Name:         req.Name,
				Image:        val.Image,
				Description:  req.Description,
				RegularPrice: val.RegularPrice,
				SalePrice:    val.SalePrice,
				Unit:         req.Unit,
				Weight:       val.Weight,
				Stock:        val.Stock,
				Variant:      req.Variant,
				Status:       req.Status,
			})

		}

		if err := db.WithContext(ctx).CreateInBatches(&modelProductChild, 50).Error; err != nil {
			p.logger.Errorf("[ProductRepository] CreateProduct: %v", err)
			return 0, err
		}
	}

	return modelProduct.ID, nil
}

// GetProductById implements ProductRepositoryInterface.
func (p *productRepository) GetProductById(ctx context.Context, productId int64) (*entity.ProductEntity, error) {
	var (
		db            = p.getDB(ctx)
		modelProducts []model.Product
		productEntity *entity.ProductEntity
	)

	if err := db.WithContext(ctx).
		Order("id ASC").
		Omit("updated_at", "deleted_at").
		Find(&modelProducts, "id = ? OR parent_id = ?", productId, productId).Error; err != nil {
		p.logger.Errorf("[ProductRepository] GetProductById: %v", err)
		return nil, err
	}

	if len(modelProducts) == 0 || (len(modelProducts) > 0 && modelProducts[0].ParentID != nil) {
		err := utils.ErrDataNotFound
		p.logger.Errorf("[ProductRepository] GetProductById: %v", err)
		return nil, err
	}

	for _, val := range modelProducts {
		if val.ParentID == nil {
			productEntity = &entity.ProductEntity{
				ID:           val.ID,
				Name:         val.Name,
				Image:        val.Image,
				Description:  val.Description,
				RegularPrice: val.RegularPrice,
				SalePrice:    val.SalePrice,
				Unit:         val.Unit,
				Weight:       val.Weight,
				Stock:        val.Stock,
				Variant:      val.Variant,
				Status:       val.Status,
				CreatedAt:    val.CreatedAt,
				CategoryID:   val.CategoryID,
			}
			continue
		}

		productEntity.Childs = append(productEntity.Childs, entity.ProductEntity{
			ID:           val.ID,
			ParentID:     val.ParentID,
			Name:         val.Name,
			Image:        val.Image,
			Description:  val.Description,
			RegularPrice: val.RegularPrice,
			SalePrice:    val.SalePrice,
			Unit:         val.Unit,
			Weight:       val.Weight,
			Stock:        val.Stock,
			Variant:      val.Variant,
			Status:       val.Status,
			CreatedAt:    val.CreatedAt,
		})
	}

	return productEntity, nil
}

// GetAllProducts implements ProductRepositoryInterface.
func (p *productRepository) GetAllProducts(ctx context.Context, query entity.QueryStringProduct) ([]entity.ProductEntity, int64, int64, error) {
	var (
		db          = p.getDB(ctx)
		productsDto []model.ProductDTO
		countData   int64
	)

	orderSort := fmt.Sprintf("products.%s %s", query.OrderBy, query.OrderType)
	offset := (query.Page - 1) * query.Limit

	productSelectedField := `products.id AS product_id,
							products.parent_id AS product_parent_id,
							products.name AS product_name,
							products.image AS product_image,
							products.description AS product_description,
							products.regular_price AS product_regular_price,
							products.sale_price AS product_sale_price,
							products.unit AS product_unit,
							products.weight AS product_weight,
							products.stock AS product_stock,
							products.variant AS product_variant,
							products.status AS product_status,
							products.created_at AS product_created_at,
							categories.name AS category_name`

	sqlMain := db.WithContext(ctx).
		Model(&model.Product{}).
		Select(productSelectedField).
		Joins("LEFT JOIN categories ON categories.id = products.category_id").
		Where("products.parent_id IS NULL").
		Where("products.status = ?", "active")

	if query.Search != "" {
		sqlMain = sqlMain.Where(`products.name ILIKE ? OR products.description ILIKE ?`, "%"+query.Search+"%", "%"+query.Search+"%")
	}

	if query.CategoryID != 0 {
		sqlMain = sqlMain.Where("categories.id = ?", query.CategoryID)
	}

	if query.StartPrice > 0 {
		sqlMain = sqlMain.Where("products.sale_price >= ?", query.StartPrice)
	}

	if query.EndPrice > 0 {
		sqlMain = sqlMain.Where("products.sale_price <= ?", query.EndPrice)
	}

	if err := sqlMain.Count(&countData).Error; err != nil {
		p.logger.Errorf("[ProductRepository] GetAllProducts: %v", err)
		return nil, 0, 0, err
	}

	totalPage := int(math.Ceil(float64(countData) / float64(query.Limit)))
	if err := sqlMain.Order(orderSort).
		Limit(int(query.Limit)).
		Offset(int(offset)).
		Find(&productsDto).Error; err != nil {
		p.logger.Errorf("[ProductRepository] GetAllProducts: %v", err)
		return nil, 0, 0, err
	}

	entities := []entity.ProductEntity{}
	for _, val := range productsDto {
		entities = append(entities, entity.ProductEntity{
			ID:           val.ProductID,
			ParentID:     &val.ProductParentID,
			Name:         val.ProductName,
			Image:        val.ProductImage,
			Description:  val.ProductDescription,
			RegularPrice: val.ProductRegularPrice,
			SalePrice:    val.ProductSalePrice,
			Unit:         val.ProductUnit,
			Weight:       val.ProductWeight,
			Stock:        val.ProductStock,
			Variant:      val.ProductVariant,
			Status:       val.ProductStatus,
			CategoryName: val.CategoryName,
			CreatedAt:    val.ProductCreatedAt,
		})
	}

	return entities, countData, int64(totalPage), nil
}
