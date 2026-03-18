package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"product-service/internal/core/domain/entity"
	"product-service/internal/core/domain/model"
	"product-service/utils"
	"strings"

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

	SearchProducts(ctx context.Context, query entity.QueryStringProduct) ([]entity.ProductEntity, int64, int64, error)

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
			Select("id", "image", "name", "sale_price", "weight", "unit").
			Where("id IN ?", productIds[i:end]).
			Find(&batchProducts).Error; err != nil {
			p.logger.Errorf("[ProductRepository-1] GetBatchProducts: %v", err)
			return nil, err
		}

		modelProducts = append(modelProducts, batchProducts...)
	}

	if len(modelProducts) == 0 {
		err := errors.New(utils.DATA_NOT_FOUND)
		p.logger.Errorf("[ProductRepository-2] GetBatchProducts: %v", err)
		return nil, err
	}

	entities := []entity.ProductEntity{}
	for _, val := range modelProducts {
		entities = append(entities, entity.ProductEntity{
			ID:        val.ID,
			Name:      val.Name,
			Image:     val.Image,
			SalePrice: val.SalePrice,
			Weight:    val.Weight,
			Unit:      val.Unit,
		})
	}

	return entities, nil
}

// SearchProducts implements ProductRepositoryInterface.
func (p *productRepository) SearchProducts(ctx context.Context, query entity.QueryStringProduct) ([]entity.ProductEntity, int64, int64, error) {
	offset := (query.Page - 1) * query.Limit

	categoryFilter := ""
	if query.CategoryID != 0 {
		categoryFilter = fmt.Sprintf(`{ "match": { "category_id": "%d" } },`, query.CategoryID)
	}

	priceFilter := ""
	if query.StartPrice > 0 && query.EndPrice > 0 {
		priceFilter = fmt.Sprintf(`{ "range": { "regular_price": { "gte": %d, "lte": %d } } }`, query.StartPrice, query.EndPrice)
	}

	searchFilter := `{ "match_all": {} }`
	if query.Search != "" {
		searchFilter = fmt.Sprintf(`{ "multi_match": { "query": "%s", "fields": ["name", "description"] } }`, query.Search)
	}

	mainQuery := fmt.Sprintf(`{
		"from": %d,
		"size": %d,
		"query": {
			"bool": {
				"must": [
					%s
					%s
				],
				"filter": [
					%s
				]
			}
		},
		"sort": [
			{ "id": "asc" }
		]
	}`, offset, query.Limit, categoryFilter, searchFilter, priceFilter)

	res, err := p.esClient.Search(
		p.esClient.Search.WithContext(ctx),
		p.esClient.Search.WithIndex("products"),
		p.esClient.Search.WithBody(strings.NewReader(mainQuery)),
		p.esClient.Search.WithPretty(),
	)
	if err != nil {
		p.logger.Errorf("[ProductRepository-1] SearchProducts: %v", err)
		return nil, 0, 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err := errors.New(res.Status())
		p.logger.Errorf("[ProductRepository-2] SearchProducts: %v", err)
		return nil, 0, 0, err
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		p.logger.Errorf("[ProductRepository-3] SearchProducts: %v", err)
		return nil, 0, 0, err
	}

	totalData := int64(0)
	if hitsTotal, found := result["hits"].(map[string]any)["total"].(map[string]any); found {
		totalData = int64(hitsTotal["value"].(float64))
	}

	totalPage := int64(0)
	if query.Limit > 0 {
		totalPage = int64(math.Ceil(float64(totalData) / float64(query.Limit)))
	}

	products := []entity.ProductEntity{}
	hits, found := result["hits"].(map[string]any)["hits"].([]any)
	if found {
		for _, hit := range hits {
			product := entity.ProductEntity{}

			source := hit.(map[string]any)["_source"]
			data, _ := json.Marshal(source)
			json.Unmarshal(data, &product)

			products = append(products, product)
		}
	}

	return products, totalData, totalPage, nil
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
		p.logger.Errorf("[ProductRepository-1] DeleteProduct: %v", tx.Error)
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		err := errors.New(utils.DATA_NOT_FOUND)
		p.logger.Errorf("[ProductRepository-2] DeleteProduct: %v", err)
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
		p.logger.Errorf("[ProductRepository-1] UpdateProduct: %v", tx.Error)
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		err := errors.New(utils.DATA_NOT_FOUND)
		p.logger.Errorf("[ProductRepository-2] UpdateProduct: %v", err)
		return err
	}

	modelProductChild := []model.Product{}
	if len(req.Childs) > 0 {
		if err := db.WithContext(ctx).
			Where("parent_id = ?", req.ID).
			Delete(&model.Product{}).Error; err != nil {
			p.logger.Errorf("[ProductRepository-3] UpdateProduct: %v", err)
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
			p.logger.Errorf("[ProductRepository-4] UpdateProduct: %v", err)
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
		p.logger.Errorf("[ProductRepository-1] CreateProduct: %v", err)
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
			p.logger.Errorf("[ProductRepository-2] CreateProduct: %v", err)
			return 0, err
		}
	}

	childEntities := []entity.ProductEntity{}
	for _, val := range modelProductChild {
		childEntities = append(childEntities, entity.ProductEntity{
			ID:           val.ID,
			CategorySlug: val.Categories.Slug,
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
			CategoryName: val.Categories.Name,
			CreatedAt:    val.CreatedAt,
			Childs:       nil,
		})
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

	if err := db.WithContext(ctx).Debug().
		Order("id ASC").
		Omit("updated_at", "deleted_at").
		Find(&modelProducts, "id = ? OR parent_id = ?", productId, productId).Error; err != nil {
		p.logger.Errorf("[ProductRepository-1] GetProductById: %v", err)
		return nil, err
	}

	if len(modelProducts) == 0 || (len(modelProducts) > 0 && modelProducts[0].ParentID != nil) {
		err := errors.New(utils.DATA_NOT_FOUND)
		p.logger.Errorf("[ProductRepository-2] GetProductById: %v", err)
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
		p.logger.Errorf("[ProductRepository-1] GetAllProducts: %v", err)
		return nil, 0, 0, err
	}

	totalPage := int(math.Ceil(float64(countData) / float64(query.Limit)))
	if err := sqlMain.Order(orderSort).Limit(int(query.Limit)).Offset(int(offset)).Find(&productsDto).Error; err != nil {
		p.logger.Errorf("[ProductRepository-2] GetAllProducts: %v", err)
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
