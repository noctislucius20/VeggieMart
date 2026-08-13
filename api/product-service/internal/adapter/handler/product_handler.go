package handler

import (
	"net/http"
	"product-service/config"
	"product-service/internal/adapter"
	"product-service/internal/adapter/handler/request"
	"product-service/internal/adapter/handler/response"
	"product-service/internal/core/domain/entity"
	"product-service/internal/core/service"
	middlewareGateway "product-service/internal/middleware"
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
	UpdateStockProduct(c echo.Context) error

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

	productGroup := e.Group("/products")
	internalProductGroup := e.Group("/internal/products")

	productGroup.Use(middlewareGateway.GatewayValidationMiddleware(cfg))
	internalProductGroup.Use(middlewareGateway.InternalServiceMiddleware(cfg))

	productGroup.GET("/home", productHandler.GetAllProductsHome)
	productGroup.GET("/shop", productHandler.GetAllProductsShop)
	productGroup.GET("/home/:id", productHandler.GetDetailProductHome)

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
	productGroup.GET("", productHandler.GetAllProductsAdmin, mid.CheckToken(), mid.RequiredPermission(adminPermission...))
	productGroup.GET("/:id", productHandler.GetProductByIdAdmin, mid.CheckToken(), mid.RequiredPermission(adminPermission...))
	productGroup.POST("", productHandler.CreateProductAdmin, mid.CheckToken(), mid.RequiredPermission(adminPermission...))
	productGroup.DELETE("/:id", productHandler.DeleteProductAdmin, mid.CheckToken(), mid.RequiredPermission(adminPermission...))
	productGroup.PUT("/:id", productHandler.UpdateProductAdmin, mid.CheckToken(), mid.RequiredPermission(adminPermission...))

	// authGroup := e.Group("auth", mid.CheckToken())
	internalProductGroup.POST("/batch", productHandler.GetBatchProducts, mid.CheckToken(), mid.RequiredPermission(authPermission...))
	internalProductGroup.POST("/stock", productHandler.UpdateStockProduct, mid.CheckToken(), mid.RequiredPermission(authPermission...))

	return productHandler
}

// UpdateStockProduct implements [ProductHandlerInterface].
func (p *productHandler) UpdateStockProduct(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		req = []request.ProductUpdateStockRequest{}
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[ProductHandler] UpdateStockProduct: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[ProductHandler] UpdateStockProduct: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[ProductHandler] UpdateStockProduct: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := []entity.ProductUpdateStockEntity{}
	for _, rp := range req {
		reqEntity = append(reqEntity, entity.ProductUpdateStockEntity{
			ProductID: rp.ProductID,
			Quantity:  rp.Quantity,
		})
	}

	if err := p.service.UpdateStockProduct(ctx, reqEntity); err != nil {
		c.Logger().Errorf("[ProductHandler] UpdateStockProduct: %v", err)
		switch err.Error() {
		case utils.ErrStockUnavailable.Error():
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}

// GetDetailProductHome implements ProductHandlerInterface.
func (p *productHandler) GetDetailProductHome(c echo.Context) error {
	var (
		ctx            = c.Request().Context()
		respHomeDetail = response.ProductHomeDetailResponse{}
	)

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[ProductHandler] GetDetailProductHome: %v", utils.ErrIDRequired.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ErrIDRequired.Error()))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[ProductHandler] GetDetailProductHome: %v", utils.ErrIDInvalid.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ErrIDInvalid.Error()))
	}

	result, err := p.service.GetProductById(ctx, id)
	if err != nil {
		c.Logger().Errorf("[ProductHandler] GetDetailProductHome: %v", err)
		switch err.Error() {
		case utils.ErrDataNotFound.Error():
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		case utils.ErrRelationDataNotFound.Error():
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
		}
	}

	respHomeDetail = response.ProductHomeDetailResponse{
		ID:           result.ID,
		ProductName:  result.Name,
		Weight:       result.Weight,
		Image:        result.Image,
		Stock:        result.Stock,
		RegularPrice: int64(result.RegularPrice),
		SalePrice:    int64(result.SalePrice),
		CategoryName: result.CategoryName,
		CategorySlug: result.CategorySlug,
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
	status := "ACTIVE"
	price := c.QueryParam("price")

	orderBy := "created_at"
	orderType := "desc"
	if c.QueryParam("order_by") != "" {
		switch strings.ToLower(c.QueryParam("order_by")) {
		case "price_asc":
			orderBy = "sale_price"
			orderType = "asc"
		case "price_desc":
			orderBy = "sale_price"
			orderType = "desc"
		case "newest":
			orderBy = "id"
			orderType = "desc"
		}
	}

	page, err := conv.ParseInt64QueryParam(c, "page", 1)
	if err != nil {
		c.Logger().Errorf("[ProductHandler] GetAllProductsShop: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	limit, err := conv.ParseInt64QueryParam(c, "limit", 10)
	if err != nil {
		c.Logger().Errorf("[ProductHandler] GetAllProductsShop: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	categoryId, err := conv.ParseInt64QueryParam(c, "category", 0)
	if err != nil {
		c.Logger().Errorf("[ProductHandler] GetAllProductsShop: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	startPrice, endPrice, err := conv.RangePriceFormat(price)
	if err != nil {
		c.Logger().Errorf("[ProductHandler] GetAllProductsShop: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ErrPriceRangeInvalid.Error()))
	}

	if startPrice > 0 && endPrice > 0 && startPrice > endPrice {
		c.Logger().Errorf("[ProductHandler] GetAllProductsShop: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ErrPriceRangeInvalid.Error()))
	}

	reqEntity := entity.QueryStringProduct{
		CategoryID: categoryId,
		Status:     status,
		Search:     search,
		Page:       page,
		Limit:      limit,
		OrderBy:    orderBy,
		OrderType:  orderType,
		StartPrice: startPrice,
		EndPrice:   endPrice,
	}

	results, countData, totalPage, err := p.service.GetAllProducts(ctx, reqEntity)
	if err != nil {
		c.Logger().Errorf("[ProductHandler] GetAllProductsShop: %v", err)
		switch err.Error() {
		case utils.ErrRelationDataNotFound.Error():
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
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
	status := "ACTIVE"
	page := int64(1)
	limit := int64(5)

	reqEntity := entity.QueryStringProduct{
		Status:    status,
		Page:      page,
		Limit:     limit,
		OrderBy:   orderBy,
		OrderType: orderType,
	}

	results, _, _, err := p.service.GetAllProducts(ctx, reqEntity)
	if err != nil {
		c.Logger().Errorf("[ProductHandler] GetAllProductsHome: %v", err)
		switch err.Error() {
		case utils.ErrRelationDataNotFound.Error():
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
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

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[ProductHandler] DeleteProductAdmin: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[ProductHandler] DeleteProductAdmin: %v", utils.ErrIDRequired.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ErrIDRequired.Error()))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[ProductHandler] DeleteProductAdmin: %v", utils.ErrIDInvalid.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ErrIDInvalid.Error()))
	}

	if err := p.service.DeleteProduct(ctx, id); err != nil {
		c.Logger().Errorf("[ProductHandler] DeleteProductAdmin: %v", err)
		switch err.Error() {
		case utils.ErrDataNotFound.Error():
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusNoContent, response.ResponseSuccess(nil))
}

// UpdateProductAdmin implements ProductHandlerInterface.
func (p *productHandler) UpdateProductAdmin(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		req = request.ProductRequest{}
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[ProductHandler] UpdateProductAdmin: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[ProductHandler] UpdateProductAdmin: %v", utils.ErrIDRequired.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ErrIDRequired.Error()))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[ProductHandler] UpdateProductAdmin: %v", utils.ErrIDInvalid.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ErrIDInvalid.Error()))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[ProductHandler] UpdateProductAdmin: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[ProductHandler] UpdateProductAdmin: %v", err)
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
		c.Logger().Errorf("[ProductHandler] UpdateProductAdmin: %v", err)
		switch err.Error() {
		case utils.ErrDataNotFound.Error():
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		case utils.ErrRelationDataNotFound.Error():
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
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

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[ProductHandler] CreateProductAdmin: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[ProductHandler] CreateProductAdmin: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[ProductHandler] CreateProductAdmin: %v", err)
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
		case utils.ErrRelationDataNotFound.Error():
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
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

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[ProductHandler] GetProductByIdAdmin: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[ProductHandler] GetProductByIdAdmin: %v", utils.ErrIDRequired.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ErrIDRequired.Error()))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[ProductHandler] GetProductByIdAdmin: %v", utils.ErrIDInvalid.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ErrIDInvalid.Error()))
	}

	result, err := p.service.GetProductById(ctx, id)
	if err != nil {
		c.Logger().Errorf("[ProductHandler] GetProductByIdAdmin: %v", err)
		switch err.Error() {
		case utils.ErrDataNotFound.Error():
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		case utils.ErrRelationDataNotFound.Error():
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
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
				ProductImage: child.Image,
			})
		}
	}

	respProduct = response.ProductDetailResponse{
		ID:                 result.ID,
		ProductName:        result.Name,
		ParentID:           result.ParentID,
		ProductImage:       result.Image,
		CategoryName:       result.CategoryName,
		ProductDescription: result.Description,
		ProductStatus:      result.Status,
		SalePrice:          int64(result.SalePrice),
		RegularPrice:       int64(result.RegularPrice),
		CreatedAt:          result.CreatedAt,
		Unit:               result.Unit,
		Weight:             result.Weight,
		Stock:              result.Stock,
		Childs:             responseChilds,
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

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[ProductHandler] GetBatchProducts: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[ProductHandler] GetBatchProducts: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[ProductHandler] GetBatchProducts: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	results, err := p.service.GetBatchProducts(ctx, req.IDProducts)
	if err != nil {
		switch err.Error() {
		case utils.ErrDataNotFound.Error():
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
		}
	}

	for _, result := range results {
		respBatch = append(respBatch, response.ProductBatchResponse{
			ID:           result.ID,
			ProductImage: result.Image,
			ProductName:  result.Name,
			RegularPrice: int64(result.RegularPrice),
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

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[ProductHandler] GetAllProductsAdmin: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	search := c.QueryParam("search")
	status := c.QueryParam("status")
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
		c.Logger().Errorf("[ProductHandler] GetAllProductsAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	page, err := conv.ParseInt64QueryParam(c, "page", 1)
	if err != nil {
		c.Logger().Errorf("[ProductHandler] GetAllProductsAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	limit, err := conv.ParseInt64QueryParam(c, "limit", 5)
	if err != nil {
		c.Logger().Errorf("[ProductHandler] GetAllProductsAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	startPrice, err := conv.ParseInt64QueryParam(c, "start_price", 0)
	if err != nil {
		c.Logger().Errorf("[ProductHandler] GetAllProductsAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	endPrice, err := conv.ParseInt64QueryParam(c, "end_price", 0)
	if err != nil {
		c.Logger().Errorf("[ProductHandler] GetAllProductsAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.QueryStringProduct{
		Status:     status,
		Search:     search,
		Page:       page,
		Limit:      limit,
		OrderBy:    orderBy,
		OrderType:  orderType,
		CategoryID: categoryId,
		StartPrice: startPrice,
		EndPrice:   endPrice,
	}

	results, countData, totalPages, err := p.service.GetAllProducts(ctx, reqEntity)
	if err != nil {
		c.Logger().Errorf("[ProductHandler] GetAllProductsAdmin: %v", err)
		switch err.Error() {
		case utils.ErrRelationDataNotFound.Error():
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
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
