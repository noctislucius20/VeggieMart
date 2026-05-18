package handler

import (
	"net/http"
	"product-service/config"
	"product-service/internal/adapter"
	"product-service/internal/adapter/handler/request"
	"product-service/internal/adapter/handler/response"
	"product-service/internal/core/domain/entity"
	"product-service/internal/core/service"
	"product-service/utils"
	"product-service/utils/conv"
	"product-service/utils/logger"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ProductHandlerInterface interface {
	GetAllProductsAdmin(c echo.Context) error
	GetBatchProducts(c echo.Context) error
	GetProductByIdAdmin(c echo.Context) error
	CreateProductAdmin(c echo.Context) error
	UpdateProductAdmin(c echo.Context) error
	DeleteProductAdmin(c echo.Context) error

	GetAllProductsHome(c echo.Context) error
	GetAllProductsShop(c echo.Context) error
	GetDetailProductHome(c echo.Context) error
}

type productHandler struct {
	service service.ProductServiceInterface
}

func NewProductHandler(e *echo.Echo, cfg *config.Config, service service.ProductServiceInterface, jwtService service.JwtServiceInterface, redisClient *redis.Client) ProductHandlerInterface {
	productHandler := &productHandler{service: service}

	e.Use(middleware.Recover())
	e.Use(middleware.ContextTimeoutWithConfig(middleware.ContextTimeoutConfig{
		Timeout: 10 * time.Second,
	}))

	mid := adapter.NewMiddlewareAdapter(cfg, logger.NewLogger().Logger(), jwtService, redisClient)

	// homeProduct := e.Group("/products")
	e.GET("/products/home", productHandler.GetAllProductsHome)
	e.GET("/products/shop", productHandler.GetAllProductsShop)
	e.GET("/products/home/:id", productHandler.GetDetailProductHome)

	adminPermission := []string{
		"products:read:all",
		"products:write:all",
		"products:update:all",
		"products:delete:all",
	}

	authPermission := []string{
		"products:read:own",
		"products:write:own",
		"products:update:own",
		"products:delete:own",
	}

	// adminGroup := e.Group("/admin", mid.CheckToken())
	e.GET("/products", productHandler.GetAllProductsAdmin, mid.CheckToken(), mid.RequiredPermission(adminPermission...))
	e.GET("/products/:id", productHandler.GetProductByIdAdmin, mid.CheckToken(), mid.RequiredPermission(adminPermission...))
	e.POST("/products", productHandler.CreateProductAdmin, mid.CheckToken(), mid.RequiredPermission(adminPermission...))
	e.DELETE("/products/:id", productHandler.DeleteProductAdmin, mid.CheckToken(), mid.RequiredPermission(adminPermission...))
	e.PUT("/products/:id", productHandler.UpdateProductAdmin, mid.CheckToken(), mid.RequiredPermission(adminPermission...))

	// authGroup := e.Group("auth", mid.CheckToken())
	e.POST("/products/batch", productHandler.GetBatchProducts, mid.CheckToken(), mid.RequiredPermission(authPermission...))

	return productHandler
}

// GetDetailProductHome implements ProductHandlerInterface.
func (p *productHandler) GetDetailProductHome(c echo.Context) error {
	var (
		ctx            = c.Request().Context()
		respHomeDetail = response.ProductHomeDetailResponse{}
	)

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[ProductHandler-1] GetDetailProductHome: %v", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ID_REQUIRED))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[ProductHandler-2] GetDetailProductHome: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ID_INVALID))
	}

	result, err := p.service.GetProductById(ctx, id)
	if err != nil {
		c.Logger().Errorf("[ProductHandler-3] GetDetailProductHome: %v", err)
		switch err.Error() {
		case utils.DATA_NOT_FOUND:
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		case utils.RELATION_DATA_NOT_FOUND:
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
		}
	}

	respHomeDetail = response.ProductHomeDetailResponse{
		ID:           result.ID,
		ProductName:  result.Name,
		CategoryName: result.CategoryName,
		Description:  result.Description,
		Unit:         result.Unit,
	}

	for _, child := range result.Childs {
		respHomeDetail.Childs = append(respHomeDetail.Childs, response.ProductHomeChildResponse{
			ID:           child.ID,
			Weight:       child.Weight,
			Stock:        child.Stock,
			RegularPrice: int64(child.RegularPrice),
			SalePrice:    int64(child.SalePrice),
			Image:        child.Image,
		})
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respHomeDetail))
}

// GetAllProductsShop implements ProductHandlerInterface.
func (p *productHandler) GetAllProductsShop(c echo.Context) error {
	var (
		ctx          = c.Request().Context()
		respHomeList = []response.ProductHomeListResponse{}
	)

	search := c.QueryParam("search")

	orderBy := "created_at"
	orderType := "desc"
	if c.QueryParam("order_by") != "" {
		switch strings.ToLower(c.QueryParam("order_by")) {
		case "price_asc":
			orderBy = "regular_price"
			orderType = "asc"
		case "price_desc":
			orderBy = "regular_price"
			orderType = "desc"
		case "newest":
			orderBy = "id"
			orderType = "desc"
		}
	}

	page, err := conv.ParseInt64QueryParam(c, "page", 1)
	if err != nil {
		c.Logger().Errorf("[ProductHandler-1] GetAllProductsShop: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	limit, err := conv.ParseInt64QueryParam(c, "limit", 10)
	if err != nil {
		c.Logger().Errorf("[ProductHandler-2] GetAllProductsShop: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	startPrice, err := conv.ParseInt64QueryParam(c, "start_price", 0)
	if err != nil {
		c.Logger().Errorf("[ProductHandler-3] GetAllProductsShop: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	endPrice, err := conv.ParseInt64QueryParam(c, "end_price", 0)
	if err != nil {
		c.Logger().Errorf("[ProductHandler-4] GetAllProductsShop: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	// TODO: check if price range is valid

	reqEntity := entity.QueryStringProduct{
		Search:     search,
		Page:       page,
		Limit:      limit,
		OrderBy:    orderBy,
		OrderType:  orderType,
		StartPrice: startPrice,
		EndPrice:   endPrice,
	}

	results, totalPage, countData, err := p.service.GetAllProducts(ctx, reqEntity)
	if err != nil {
		c.Logger().Errorf("[ProductHandler-5] GetAllProductsShop: %v", err)
		switch err.Error() {
		case utils.RELATION_DATA_NOT_FOUND:
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
		}
	}

	for _, result := range results {
		respHomeList = append(respHomeList, response.ProductHomeListResponse{
			ID:           result.ID,
			ProductName:  result.Name,
			ProductImage: result.Image,
			CategoryName: result.CategoryName,
			SalePrice:    int64(result.SalePrice),
			RegularPrice: int64(result.RegularPrice),
		})
	}

	pagination := response.Pagination{
		Page:       page,
		TotalCount: countData,
		PerPage:    limit,
		TotalPage:  totalPage,
	}

	return c.JSON(http.StatusOK, response.ResponseWithPaginationsSuccess(respHomeList, pagination))
}

// GetAllProductsHome implements ProductHandlerInterface.
func (p *productHandler) GetAllProductsHome(c echo.Context) error {
	var (
		ctx          = c.Request().Context()
		respHomeList = []response.ProductHomeListResponse{}
	)

	orderBy := "created_at"
	orderType := "desc"
	page := int64(1)
	limit := int64(5)

	reqEntity := entity.QueryStringProduct{
		Page:      page,
		Limit:     limit,
		OrderBy:   orderBy,
		OrderType: orderType,
	}

	results, _, _, err := p.service.GetAllProducts(ctx, reqEntity)
	if err != nil {
		c.Logger().Errorf("[ProductHandler-1] GetAllProductsHome: %v", err)
		switch err.Error() {
		case utils.RELATION_DATA_NOT_FOUND:
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
		}
	}

	for _, result := range results {
		respHomeList = append(respHomeList, response.ProductHomeListResponse{
			ID:           result.ID,
			ProductName:  result.Name,
			ProductImage: result.Image,
			CategoryName: result.CategoryName,
			SalePrice:    int64(result.SalePrice),
			RegularPrice: int64(result.RegularPrice),
		})
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respHomeList))
}

// DeleteProductAdmin implements ProductHandlerInterface.
func (p *productHandler) DeleteProductAdmin(c echo.Context) error {
	var (
		ctx = c.Request().Context()
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[ProductHandler-1] DeleteProductAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[ProductHandler-2] DeleteProductAdmin: %v", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ID_REQUIRED))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[ProductHandler-3] DeleteProductAdmin: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ID_INVALID))
	}

	if err := p.service.DeleteProduct(ctx, id); err != nil {
		c.Logger().Errorf("[ProductHandler-4] DeleteProductAdmin: %v", err)
		switch err.Error() {
		case utils.DATA_NOT_FOUND:
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
		}
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}

// UpdateProductAdmin implements ProductHandlerInterface.
func (p *productHandler) UpdateProductAdmin(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		req = request.ProductRequest{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[ProductHandler-1] UpdateProductAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[ProductHandler-2] UpdateProductAdmin: %v", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ID_REQUIRED))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[ProductHandler-3] UpdateProductAdmin: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ID_INVALID))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[ProductHandler-4] UpdateProductAdmin: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[ProductHandler-5] UpdateProductAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.ProductEntity{
		ID:           id,
		CategorySlug: req.CategorySlug,
		ParentID:     nil,
		Name:         req.ProductName,
		Image:        req.VariantDetail[0].ProductImage,
		Description:  req.ProductDescription,
		RegularPrice: float64(req.VariantDetail[0].RegularPrice),
		SalePrice:    float64(req.VariantDetail[0].SalePrice),
		Unit:         req.Unit,
		Weight:       req.VariantDetail[0].Weight,
		Stock:        req.VariantDetail[0].Stock,
		Variant:      req.Variant,
		Status:       req.Status,
	}

	productChilds := []entity.ProductEntity{}
	if len(req.VariantDetail) > 1 {
		for i := 1; i < len(req.VariantDetail); i++ {
			productChilds = append(productChilds, entity.ProductEntity{
				Image:        req.VariantDetail[i].ProductImage,
				RegularPrice: float64(req.VariantDetail[i].RegularPrice),
				SalePrice:    float64(req.VariantDetail[i].SalePrice),
				Weight:       req.VariantDetail[i].Weight,
				Stock:        req.VariantDetail[i].Stock,
			})
		}

		reqEntity.Childs = productChilds
	}

	if err := p.service.UpdateProduct(ctx, reqEntity); err != nil {
		c.Logger().Errorf("[ProductHandler-6] UpdateProductAdmin: %v", err)
		switch err.Error() {
		case utils.DATA_NOT_FOUND:
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		case utils.RELATION_DATA_NOT_FOUND:
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
		}
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}

// CreateProductAdmin implements ProductHandlerInterface.
func (p *productHandler) CreateProductAdmin(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		req = request.ProductRequest{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[ProductHandler-1] CreateProductAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[ProductHandler-2] CreateProductAdmin: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[ProductHandler-3] CreateProductAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.ProductEntity{
		CategorySlug: req.CategorySlug,
		ParentID:     nil,
		Name:         req.ProductName,
		Image:        req.VariantDetail[0].ProductImage,
		Description:  req.ProductDescription,
		RegularPrice: float64(req.VariantDetail[0].RegularPrice),
		SalePrice:    float64(req.VariantDetail[0].SalePrice),
		Unit:         req.Unit,
		Weight:       req.VariantDetail[0].Weight,
		Stock:        req.VariantDetail[0].Stock,
		Variant:      req.Variant,
		Status:       req.Status,
	}

	productChilds := []entity.ProductEntity{}
	if len(req.VariantDetail) > 1 {
		for i := 1; i < len(req.VariantDetail); i++ {
			productChilds = append(productChilds, entity.ProductEntity{
				Image:        req.VariantDetail[i].ProductImage,
				RegularPrice: float64(req.VariantDetail[i].RegularPrice),
				SalePrice:    float64(req.VariantDetail[i].SalePrice),
				Weight:       req.VariantDetail[i].Weight,
				Stock:        req.VariantDetail[i].Stock,
			})
		}

		reqEntity.Childs = productChilds
	}

	productId, err := p.service.CreateProduct(ctx, reqEntity)
	if err != nil {
		switch err.Error() {
		case utils.RELATION_DATA_NOT_FOUND:
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
		}
	}

	respProductId := map[string]int64{
		"product_id": productId,
	}

	return c.JSON(http.StatusCreated, response.ResponseSuccess(respProductId))
}

// GetProductByIdAdmin implements ProductHandlerInterface.
func (p *productHandler) GetProductByIdAdmin(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		respProduct = response.ProductDetailResponse{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[ProductHandler-1] GetProductByIdAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[ProductHandler-2] GetProductByIdAdmin: %v", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ID_REQUIRED))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[ProductHandler-3] GetProductByIdAdmin: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ID_INVALID))
	}

	result, err := p.service.GetProductById(ctx, id)
	if err != nil {
		c.Logger().Errorf("[ProductHandler-4] GetProductByIdAdmin: %v", err)
		switch err.Error() {
		case utils.DATA_NOT_FOUND:
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		case utils.RELATION_DATA_NOT_FOUND:
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
		}
	}

	responseChilds := []response.ProductChildResponse{}
	if len(result.Childs) > 0 {
		for _, child := range result.Childs {
			responseChilds = append(responseChilds, response.ProductChildResponse{
				ID:           child.ID,
				Weight:       child.Weight,
				Stock:        child.Stock,
				RegularPrice: int64(child.RegularPrice),
				SalePrice:    int64(child.SalePrice),
			})
		}
	}

	respProduct = response.ProductDetailResponse{
		ID:            result.ID,
		ProductName:   result.Name,
		ParentID:      result.ParentID,
		ProductImage:  result.Image,
		CategoryName:  result.CategoryName,
		ProductStatus: result.Status,
		SalePrice:     int64(result.SalePrice),
		RegularPrice:  int64(result.RegularPrice),
		CreatedAt:     result.CreatedAt,
		Unit:          result.Unit,
		Weight:        result.Weight,
		Stock:         result.Stock,
		Childs:        responseChilds,
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respProduct))
}

// GetBatchProducts implements [ProductHandlerInterface].
func (p *productHandler) GetBatchProducts(c echo.Context) error {
	var (
		ctx       = c.Request().Context()
		respBatch = []response.ProductBatchResponse{}
		req       = request.ProductBatchRequest{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[ProductHandler-1] GetBatchProducts: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[ProductHandler-2] GetBatchProducts: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[ProductHandler-3] GetBatchProducts: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	results, err := p.service.GetBatchProducts(ctx, req.IDProducts)
	if err != nil {
		switch err.Error() {
		case utils.DATA_NOT_FOUND:
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
		}
	}

	for _, result := range results {
		respBatch = append(respBatch, response.ProductBatchResponse{
			ID:           result.ID,
			ProductImage: result.Image,
			ProductName:  result.Name,
			SalePrice:    int64(result.SalePrice),
			Weight:       result.Weight,
			Unit:         result.Unit,
		})
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respBatch))
}

// GetAllProductsAdmin implements ProductHandlerInterface.
func (p *productHandler) GetAllProductsAdmin(c echo.Context) error {
	var (
		ctx          = c.Request().Context()
		respProducts = []response.ProductListResponse{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[ProductHandler-1] GetAllProductsAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	search := c.QueryParam("search")

	orderBy := c.QueryParam("order_by")
	if orderBy == "" {
		orderBy = "created_at"
	}

	orderType := c.QueryParam("order_type")
	if orderType != "asc" && orderType != "desc" {
		orderType = "desc"
	}

	categoryId, err := conv.ParseInt64QueryParam(c, "category", 0)
	if err != nil {
		c.Logger().Errorf("[ProductHandler-2] GetAllProductsAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	page, err := conv.ParseInt64QueryParam(c, "page", 1)
	if err != nil {
		c.Logger().Errorf("[ProductHandler-3] GetAllProductsAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	limit, err := conv.ParseInt64QueryParam(c, "limit", 10)
	if err != nil {
		c.Logger().Errorf("[ProductHandler-4] GetAllProductsAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	startPrice, err := conv.ParseInt64QueryParam(c, "start_price", 0)
	if err != nil {
		c.Logger().Errorf("[ProductHandler-5] GetAllProductsAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	endPrice, err := conv.ParseInt64QueryParam(c, "end_price", 0)
	if err != nil {
		c.Logger().Errorf("[ProductHandler-6] GetAllProductsAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.QueryStringProduct{
		Search:     search,
		Page:       page,
		Limit:      limit,
		OrderBy:    orderBy,
		OrderType:  orderType,
		CategoryID: categoryId,
		StartPrice: startPrice,
		EndPrice:   endPrice,
	}

	results, totalPages, countData, err := p.service.GetAllProducts(ctx, reqEntity)
	if err != nil {
		c.Logger().Errorf("[ProductHandler-7] GetAllProductsAdmin: %v", err)
		switch err.Error() {
		case utils.RELATION_DATA_NOT_FOUND:
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
		}
	}

	for _, result := range results {
		respProducts = append(respProducts, response.ProductListResponse{
			ID:            result.ID,
			ProductName:   result.Name,
			ParentID:      result.ParentID,
			ProductImage:  result.Image,
			CategoryName:  result.CategoryName,
			ProductStatus: result.Status,
			SalePrice:     int64(result.SalePrice),
			CreatedAt:     result.CreatedAt,
		})
	}

	pagination := response.Pagination{
		Page:       page,
		TotalCount: countData,
		PerPage:    limit,
		TotalPage:  totalPages,
	}

	return c.JSON(http.StatusOK, response.ResponseWithPaginationsSuccess(respProducts, pagination))
}
