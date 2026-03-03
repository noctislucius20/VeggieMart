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
		Timeout: 10 * time.Second,
	}))

	mid := adapter.NewMiddlewareAdapter(cfg, logger.NewLogger().Logger(), jwtService, redisClient)

	homeCategory := e.Group("/categories")
	homeCategory.GET("/shop", categoryHandler.GetAllCategoriesShop)
	homeCategory.GET("/home", categoryHandler.GetAllCategoriesHome)

	adminGroup := e.Group("/admin", mid.CheckToken())
	adminGroup.POST("/categories", categoryHandler.CreateCategoryAdmin)
	adminGroup.GET("/categories", categoryHandler.GetAllCategoriesAdmin)
	adminGroup.GET("/categories/:id", categoryHandler.GetCategoryByIdAdmin)
	adminGroup.GET("/categories/:slug/slug", categoryHandler.GetCategoryBySlugAdmin)
	adminGroup.PUT("/categories/:id", categoryHandler.UpdateCategoryAdmin)
	adminGroup.DELETE("/categories/:id", categoryHandler.DeleteCategoryAdmin)

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
		c.Logger().Warnf("[CategoryHandler-1] GetAllCategoriesShop: %v", err)
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(utils.DATA_NOT_FOUND))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(err.Error()))
	}

	for _, result := range results {
		// Check if the data is as a parent and populate the response list.
		// Otherwise, it will categorized as a child.
		if result.ParentID == nil {
			respCategoriesShopList = append(respCategoriesShopList, response.CategoryShopListResponse{
				Name: result.Name,
				Slug: result.Slug,
			})
			continue
		}

		// Add childs category to their parent.
		respCategoriesShopList[len(respCategoriesShopList)-1].Childs = append(respCategoriesShopList[len(respCategoriesShopList)-1].Childs, response.CategoryShopListResponse{
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
		c.Logger().Warnf("[CategoryHandler-1] GetAllCategoriesHome: %v", err)
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(utils.DATA_NOT_FOUND))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(err.Error()))
	}

	for _, result := range results {
		if result.ParentID == nil {
			respCategoriesHomeList = append(respCategoriesHomeList, response.CategoryHomeListResponse{
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

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[CategoryHandler-1] UpdateCategoryAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[CategoryHandler-2] UpdateCategoryAdmin: id is required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed("id is required"))
	}

	categoryId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[CategoryHandler-3] UpdateCategoryAdmin: id is invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed("id is invalid"))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[CategoryHandler-4] UpdateCategoryAdmin: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[CategoryHandler-5] UpdateCategoryAdmin: %v", err)
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
		c.Logger().Warnf("[CategoryHandler-6] UpdateCategoryAdmin: %v", err)
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(utils.DATA_NOT_FOUND))
		}
		if err.Error() == "409" {
			return c.JSON(http.StatusConflict, response.ResponseFailed("data already exists"))
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

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[CategoryHandler-1] DeleteCategoryAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[CategoryHandler-2] DeleteCategoryAdmin: id is required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed("id is required"))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[CategoryHandler-3] DeleteCategoryAdmin: id is invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed("id is invalid"))
	}

	if err := ch.categoryService.DeleteCategory(ctx, id); err != nil {
		c.Logger().Warnf("[CategoryHandler-4] DeleteCategoryAdmin: %v", err)
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(utils.DATA_NOT_FOUND))
		}
		if err.Error() == "422" {
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed("data is still being used"))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(err.Error()))
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}

// CreateCategoryAdmin implements CategoryHandlerInterface.
func (ch *categoryHandler) CreateCategoryAdmin(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		req = request.CategoryRequest{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[CategoryHandler-1] CreateCategoryAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[CategoryHandler-2] CreateCategoryAdmin: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[CategoryHandler-3] CreateCategoryAdmin: %v", err)
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
		c.Logger().Errorf("[CategoryHandler-4] CreateCategoryAdmin: %v", err)
		if err.Error() == "409" {
			return c.JSON(http.StatusConflict, response.ResponseFailed("data already exists"))
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

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[CategoryHandler-1] GetCategoryByIdAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[CategoryHandler-2] GetCategoryByIdAdmin: id is required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed("id is required"))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[CategoryHandler-3] GetCategoryByIdAdmin: id is invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed("id is invalid"))
	}

	result, err := ch.categoryService.GetCategoryByIdOrSlug(ctx, id, "")
	if err != nil {
		c.Logger().Warnf("[CategoryHandler-4] GetCategoryByIdAdmin: %v", err)
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(utils.DATA_NOT_FOUND))
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

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[CategoryHandler-1] GetCategoryBySlugAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	slug := c.Param("slug")
	if slug == "" {
		c.Logger().Errorf("[CategoryHandler-2] GetCategoryBySlugAdmin: slug is required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed("slug is required"))
	}

	result, err := ch.categoryService.GetCategoryByIdOrSlug(ctx, 0, slug)
	if err != nil {
		c.Logger().Warnf("[CategoryHandler-3] GetCategoryBySlugAdmin: %v", err)
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(utils.DATA_NOT_FOUND))
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

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[CategoryHandler-1] GetAllCategoriesAdmin: %v", "data token not found")
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

	page, err := conv.ParseInt64QueryParam(c, "page", 1)
	if err != nil {
		c.Logger().Errorf("[CategoryHandler-2] GetAllCategoriesAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	limit, err := conv.ParseInt64QueryParam(c, "limit", 10)
	if err != nil {
		c.Logger().Errorf("[CategoryHandler-3] GetAllCategoriesAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.QueryStringEntity{
		Search:    search,
		Page:      page,
		Limit:     limit,
		OrderBy:   orderBy,
		OrderType: orderType,
	}

	results, countData, totalPages, err := ch.categoryService.GetAllCategories(ctx, reqEntity)
	if err != nil {
		c.Logger().Warnf("[CategoryHandler-4] GetAllCategoriesAdmin: %v", err)
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(utils.DATA_NOT_FOUND))
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
