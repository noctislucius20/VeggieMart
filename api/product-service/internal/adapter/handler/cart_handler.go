package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"product-service/config"
	"product-service/internal/adapter"
	"product-service/internal/adapter/handler/request"
	"product-service/internal/adapter/handler/response"
	"product-service/internal/core/domain/entity"
	"product-service/internal/core/service"
	"product-service/utils"
	"product-service/utils/logger"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CartHandlerInterface interface {
	AddToCart(c echo.Context) error
	GetCart(c echo.Context) error
	RemoveFromCart(c echo.Context) error
}

type cartHandler struct {
	cartService service.CartServiceInterface
}

func NewCartHandler(e *echo.Echo, cartService service.CartServiceInterface, cfg *config.Config, jwtService service.JwtServiceInterface, redisClient *redis.Client) CartHandlerInterface {
	cartHandler := &cartHandler{cartService: cartService}

	e.Use(middleware.Recover())
	e.Use(middleware.ContextTimeoutWithConfig(middleware.ContextTimeoutConfig{
		Timeout: 10 * time.Second,
	}))

	mid := adapter.NewMiddlewareAdapter(cfg, logger.NewLogger().Logger(), jwtService, redisClient)

	authPermission := []string{
		"carts:read:all",
		"carts:write:all",
		"carts:delete:all",
	}

	// adminGroup := e.Group("/admin", mid.CheckToken())
	e.POST("/carts", cartHandler.AddToCart, mid.CheckToken(), mid.RequiredPermission(authPermission...))
	e.GET("/carts", cartHandler.GetCart, mid.CheckToken(), mid.RequiredPermission(authPermission...))
	e.DELETE("/carts", cartHandler.RemoveFromCart, mid.CheckToken(), mid.RequiredPermission(authPermission...))

	return cartHandler
}

// AddToCart implements [CartHandlerInterface].
func (ch *cartHandler) AddToCart(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		req         request.CartRequest
		jwtUserData entity.JwtUserData
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[CartHandler-1] AddToCart: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[CartHandler-2] AddToCart: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[CartHandler-3] AddToCart: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[CartHandler-4] AddToCart: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	if req.Quantity <= 0 {
		err := errors.New(utils.QUANTITY_INVALID)
		c.Logger().Errorf("[CartHandler-5] AddToCart: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.CartItem{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := ch.cartService.AddToCart(ctx, jwtUserData.UserID, reqEntity); err != nil {
		c.Logger().Errorf("[CartHandler-6] AddToCart: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	return c.JSON(http.StatusCreated, response.ResponseSuccess(nil))
}

// GetCart implements [CartHandlerInterface].
func (ch *cartHandler) GetCart(c echo.Context) error {
	var (
		ctx          = c.Request().Context()
		respCartList = []response.CartResponse{}
		jwtUserData  entity.JwtUserData
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[CartHandler-1] GetCart: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[CartHandler-2] GetCart: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	results, err := ch.cartService.GetCart(ctx, jwtUserData.UserID)
	if err != nil {
		c.Logger().Errorf("[CartHandler-3] GetCart: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	for _, result := range results {
		respCartList = append(respCartList, response.CartResponse{
			ID:            result.ProductID,
			ProductName:   result.ProductDetail.Name,
			ProductImage:  result.ProductDetail.Image,
			ProductStatus: result.ProductDetail.Status,
			SalePrice:     int64(result.ProductDetail.SalePrice),
			Quantity:      result.Quantity,
			Unit:          result.ProductDetail.Unit,
			Weight:        result.ProductDetail.Weight,
		})
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respCartList))
}

// RemoveFromCart implements CartHandlerInterface.
func (ch *cartHandler) RemoveFromCart(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		jwtUserData entity.JwtUserData
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[CartHandler-1] RemoveFromCart: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[CartHandler-2] RemoveFromCart: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	productIdParam := c.QueryParam("product_id")
	if productIdParam == "" {
		c.Logger().Errorf("[CartHandler-3] RemoveFromCart: %v", "product_id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.PRODUCT_ID_REQUIRED))
	}

	productId, err := strconv.ParseInt(productIdParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[CartHandler-4] RemoveFromCart: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.PRODUCT_ID_INVALID))
	}

	err = ch.cartService.RemoveFromCart(ctx, jwtUserData.UserID, productId)
	if err != nil {
		c.Logger().Errorf("[CartHandler-5] RemoveFromCart: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}
