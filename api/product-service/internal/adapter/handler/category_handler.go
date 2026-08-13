package handler

import (
	"errors"
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
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CategoryHandlerInterface interface {
	GetAllCategoriesAdmin(c echo.Context) error
	GetCategoryBySlugAdmin(c echo.Context) error
	GetCategoryByIdAdmin(c echo.Context) error
	CreateCategoryAdmin(c echo.Context) error
	DeleteCategoryAdmin(c echo.Context) error
	UpdateCategoryAdmin(c echo.Context) error

	GetAllCategoriesHome(c echo.Context) error
	GetAllCategoriesShop(c echo.Context) error
}

type categoryHandler struct {
	categoryService service.CategoryServiceInterface
}

func NewCategoryHandler(e *echo.Echo, categoryService service.CategoryServiceInterface, cfg *config.Config, jwtService service.JwtServiceInterface, redisClient *redis.Client) CategoryHandlerInterface {
	categoryHandler := &categoryHandler{categoryService: categoryService}

	e.Use(middleware.Recover())
	e.Use(middleware.ContextTimeoutWithConfig(middleware.ContextTimeoutConfig{
		Timeout: 100 * time.Second,
	}))

	mid := adapter.NewMiddlewareAdapter(cfg, logger.NewLogger().Logger(), jwtService, redisClient)

	productGroup := e.Group("/products")
	productGroup.Use(middlewareGateway.GatewayValidationMiddleware(cfg))

	productGroup.GET("/categories/shop", categoryHandler.GetAllCategoriesShop)
	productGroup.GET("/categories/home", categoryHandler.GetAllCategoriesHome)

	adminPermission := []string{
		"categories:read:all",
		"categories:write:all",
		"categories:update:all",
		"categories:delete:all",
	}

	// adminGroup := e.Group("/admin", mid.CheckToken())
	adminCategoryGroup := productGroup.Group("/categories", mid.CheckToken(), mid.RequiredPermission(adminPermission...))
	adminCategoryGroup.POST("", categoryHandler.CreateCategoryAdmin)
	adminCategoryGroup.GET("", categoryHandler.GetAllCategoriesAdmin)
	adminCategoryGroup.GET("/:id", categoryHandler.GetCategoryByIdAdmin)
	adminCategoryGroup.GET("/:slug/slug", categoryHandler.GetCategoryBySlugAdmin)
	adminCategoryGroup.PUT("/:id", categoryHandler.UpdateCategoryAdmin)
	adminCategoryGroup.DELETE("/:id", categoryHandler.DeleteCategoryAdmin)

	return categoryHandler
}

// GetAllCategoriesShop implements [CategoryHandlerInterface].
func (ch *categoryHandler) GetAllCategoriesShop(c echo.Context) error {
	var (
		ctx                    = c.Request().Context()
		respCategoriesShopList = []response.CategoryShopListResponse{}
	)

	results, err := ch.categoryService.GetAllCategoriesPublished(ctx)
	if err != nil {
		c.Logger().Warnf("[CategoryHandler] GetAllCategoriesShop: %v", err)
		if errors.Is(err, utils.ErrDataNotFound) {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(utils.ErrDataNotFound.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(err.Error()))
	}

	for _, result := range results {
		// Check if the data is as a parent and populate the response list.
		// Otherwise, it will categorized as a child.
		if result.ParentID == nil {
			respCategoriesShopList = append(respCategoriesShopList, response.CategoryShopListResponse{
				ID:   result.ID,
				Name: result.Name,
				Slug: result.Slug,
			})
			continue
		}

		// Add childs category to their parent.
		respCategoriesShopList[len(respCategoriesShopList)-1].Childs = append(respCategoriesShopList[len(respCategoriesShopList)-1].Childs, response.CategoryShopListResponse{
			ID:   result.ID,
			Name: result.Name,
			Slug: result.Slug,
		})
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respCategoriesShopList))
}

// GetAllCategoriesHome implements [CategoryHandlerInterface].
func (ch *categoryHandler) GetAllCategoriesHome(c echo.Context) error {
	var (
		ctx                    = c.Request().Context()
		respCategoriesHomeList = []response.CategoryHomeListResponse{}
	)

	results, err := ch.categoryService.GetAllCategoriesPublished(ctx)
	if err != nil {
		c.Logger().Warnf("[CategoryHandler] GetAllCategoriesHome: %v", err)
		if errors.Is(err, utils.ErrDataNotFound) {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(utils.ErrDataNotFound.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(err.Error()))
	}

	for _, result := range results {
		if result.ParentID == nil {
			respCategoriesHomeList = append(respCategoriesHomeList, response.CategoryHomeListResponse{
				ID:   result.ID,
				Name: result.Name,
				Icon: result.Icon,
				Slug: result.Slug,
			})
		}
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respCategoriesHomeList))
}

// UpdateCategoryAdmin implements CategoryHandlerInterface.
func (ch *categoryHandler) UpdateCategoryAdmin(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		req = request.CategoryRequest{}
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[CategoryHandler] UpdateCategoryAdmin: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[CategoryHandler] UpdateCategoryAdmin: %v", utils.ErrIDRequired.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ErrIDRequired.Error()))
	}

	categoryId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[CategoryHandler] UpdateCategoryAdmin: %v", utils.ErrIDInvalid.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ErrIDInvalid.Error()))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[CategoryHandler] UpdateCategoryAdmin: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[CategoryHandler] UpdateCategoryAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.CategoryEntity{
		ID:          categoryId,
		ParentID:    req.ParentID,
		Name:        req.Name,
		Icon:        req.Icon,
		Status:      req.Status,
		Description: req.Description,
	}

	err = ch.categoryService.UpdateCategory(ctx, reqEntity)
	if err != nil {
		c.Logger().Errorf("[CategoryHandler] UpdateCategoryAdmin: %v", err)
		if errors.Is(err, utils.ErrDataNotFound) {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(utils.ErrDataNotFound.Error()))
		}
		if errors.Is(err, utils.ErrDataAlreadyExists) {
			return c.JSON(http.StatusConflict, response.ResponseFailed(utils.ErrDataAlreadyExists.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(err.Error()))
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}

// DeleteCategoryAdmin implements CategoryHandlerInterface.
func (ch *categoryHandler) DeleteCategoryAdmin(c echo.Context) error {
	var (
		ctx = c.Request().Context()
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[CategoryHandler] DeleteCategoryAdmin: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[CategoryHandler] DeleteCategoryAdmin: %v", utils.ErrIDRequired.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ErrIDRequired.Error()))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[CategoryHandler] DeleteCategoryAdmin: %v", utils.ErrIDInvalid.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ErrIDInvalid.Error()))
	}

	if err := ch.categoryService.DeleteCategory(ctx, id); err != nil {
		c.Logger().Warnf("[CategoryHandler] DeleteCategoryAdmin: %v", err)
		if errors.Is(err, utils.ErrDataNotFound) {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(utils.ErrDataNotFound.Error()))
		}
		if errors.Is(err, utils.ErrDataStillInUsed) {
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ErrDataStillInUsed.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(err.Error()))
	}

	return c.JSON(http.StatusNoContent, response.ResponseSuccess(nil))
}

// CreateCategoryAdmin implements CategoryHandlerInterface.
func (ch *categoryHandler) CreateCategoryAdmin(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		req = request.CategoryRequest{}
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[CategoryHandler] CreateCategoryAdmin: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[CategoryHandler] CreateCategoryAdmin: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[CategoryHandler] CreateCategoryAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.CategoryEntity{
		Name:        req.Name,
		ParentID:    req.ParentID,
		Icon:        req.Icon,
		Status:      req.Status,
		Description: req.Description,
	}

	slug, categoryId, err := ch.categoryService.CreateCategory(ctx, reqEntity)
	if err != nil {
		c.Logger().Errorf("[CategoryHandler] CreateCategoryAdmin: %v", err)
		if errors.Is(err, utils.ErrDataAlreadyExists) {
			return c.JSON(http.StatusConflict, response.ResponseFailed(utils.ErrDataAlreadyExists.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(err.Error()))
	}

	respCategoryId := map[string]any{
		"category_id": categoryId,
		"slug":        slug,
	}

	return c.JSON(http.StatusCreated, response.ResponseSuccess(respCategoryId))
}

// GetCategoryByIdAdmin implements CategoryHandlerInterface.
func (ch *categoryHandler) GetCategoryByIdAdmin(c echo.Context) error {
	var (
		ctx            = c.Request().Context()
		respCategories = response.CategoryResponse{}
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[CategoryHandler] GetCategoryByIdAdmin: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[CategoryHandler] GetCategoryByIdAdmin: %v", utils.ErrIDRequired.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ErrIDRequired.Error()))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[CategoryHandler] GetCategoryByIdAdmin: %v", utils.ErrIDInvalid.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ErrIDInvalid.Error()))
	}

	result, err := ch.categoryService.GetCategoryById(ctx, id)
	if err != nil {
		c.Logger().Warnf("[CategoryHandler] GetCategoryByIdAdmin: %v", err)
		if errors.Is(err, utils.ErrDataNotFound) {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(utils.ErrDataNotFound.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(err.Error()))
	}

	respCategories = response.CategoryResponse{
		ID:          result.ID,
		Name:        result.Name,
		Icon:        result.Icon,
		Slug:        result.Slug,
		Status:      result.Status,
		Description: result.Description,
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respCategories))
}

// GetCategoryBySlugAdmin implements CategoryHandlerInterface.
func (ch *categoryHandler) GetCategoryBySlugAdmin(c echo.Context) error {
	var (
		ctx            = c.Request().Context()
		respCategories = response.CategoryResponse{}
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[CategoryHandler] GetCategoryBySlugAdmin: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	slug := c.Param("slug")
	if slug == "" {
		c.Logger().Errorf("[CategoryHandler] GetCategoryBySlugAdmin: %v", utils.ErrSlugRequired.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ErrSlugRequired.Error()))
	}

	result, err := ch.categoryService.GetCategoryBySlug(ctx, slug)
	if err != nil {
		c.Logger().Warnf("[CategoryHandler] GetCategoryBySlugAdmin: %v", err)
		if errors.Is(err, utils.ErrDataNotFound) {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(utils.ErrDataNotFound.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(err.Error()))
	}

	respCategories = response.CategoryResponse{
		ID:          result.ID,
		Name:        result.Name,
		Icon:        result.Icon,
		Slug:        result.Slug,
		Status:      result.Status,
		Description: result.Description,
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respCategories))
}

// GetAllCategoriesAdmin implements CategoryHandlerInterface.
func (ch *categoryHandler) GetAllCategoriesAdmin(c echo.Context) error {
	var (
		ctx            = c.Request().Context()
		respCategories = []response.CategoryListResponse{}
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[CategoryHandler] GetAllCategoriesAdmin: %v", utils.ErrTokenInvalid.Error())
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

	page, err := conv.ParseInt64QueryParam(c, "page", 1)
	if err != nil {
		c.Logger().Errorf("[CategoryHandler] GetAllCategoriesAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	limit, err := conv.ParseInt64QueryParam(c, "limit", 5)
	if err != nil {
		c.Logger().Errorf("[CategoryHandler] GetAllCategoriesAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.QueryStringEntity{
		Search:    search,
		Page:      page,
		Limit:     limit,
		OrderBy:   orderBy,
		OrderType: orderType,
		Status:    status,
	}

	results, countData, totalPages, err := ch.categoryService.GetAllCategories(ctx, reqEntity)
	if err != nil {
		c.Logger().Warnf("[CategoryHandler] GetAllCategoriesAdmin: %v", err)
		if errors.Is(err, utils.ErrDataNotFound) {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(utils.ErrDataNotFound.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(err.Error()))
	}

	for _, result := range results {
		respCategories = append(respCategories, response.CategoryListResponse{
			ID:           result.ID,
			Name:         result.Name,
			Icon:         result.Icon,
			Slug:         result.Slug,
			Status:       result.Status,
			TotalProduct: len(result.Products),
		})
	}

	pagination := response.Pagination{
		Page:       page,
		TotalCount: countData,
		PerPage:    limit,
		TotalPage:  totalPages,
	}

	return c.JSON(http.StatusOK, response.ResponseWithPaginationsSuccess(respCategories, pagination))
}
