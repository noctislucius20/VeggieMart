package handler

import (
	"encoding/json"
	"net/http"
	"product-service/config"
	"product-service/internal/adapter"
	"product-service/internal/adapter/handler/request"
	"product-service/internal/adapter/handler/response"
	"product-service/internal/core/domain/entity"
	"product-service/internal/core/service"
	middlewareGateway "product-service/internal/middleware"
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
	RemoveAllFromCart(c echo.Context) error
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
		"carts:read:own",
		"carts:write:own",
		"carts:update:own",
		"carts:delete:own",
	}

	// adminGroup := e.Group("/admin", mid.CheckToken())
	productGroup := e.Group("/products")
	productGroup.Use(middlewareGateway.GatewayValidationMiddleware(cfg))

	authCartGroup := productGroup.Group("/carts", mid.CheckToken(), mid.RequiredPermission(authPermission...))
	authCartGroup.POST("", cartHandler.AddToCart)
	authCartGroup.GET("", cartHandler.GetCart)
	authCartGroup.DELETE("", cartHandler.RemoveFromCart)
	authCartGroup.DELETE("/all", cartHandler.RemoveAllFromCart)

	return cartHandler
}

// AddToCart implements [CartHandlerInterface].
func (ch *cartHandler) AddToCart(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		req         request.CartRequest
		jwtUserData entity.JwtUserData
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[CartHandler] AddToCart: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[CartHandler] AddToCart: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[CartHandler] AddToCart: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[CartHandler] AddToCart: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	if req.Quantity <= 0 {
		err := utils.ErrQuantityInvalid
		c.Logger().Errorf("[CartHandler] AddToCart: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.CartItem{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := ch.cartService.AddToCart(ctx, jwtUserData.UserID, reqEntity); err != nil {
		c.Logger().Errorf("[CartHandler] AddToCart: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
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

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[CartHandler] GetCart: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[CartHandler] GetCart: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	results, err := ch.cartService.GetCart(ctx, jwtUserData.UserID)
	if err != nil {
		c.Logger().Errorf("[CartHandler] GetCart: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
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

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[CartHandler] RemoveFromCart: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[CartHandler] RemoveFromCart: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	productIdParam := c.QueryParam("product_id")
	if productIdParam == "" {
		c.Logger().Errorf("[CartHandler] RemoveFromCart: %v", utils.ErrProductIDRequired.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ErrProductIDRequired.Error()))
	}

	productId, err := strconv.ParseInt(productIdParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[CartHandler] RemoveFromCart: %v", utils.ErrProductIDInvalid.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ErrProductIDInvalid.Error()))
	}

	if err := ch.cartService.RemoveFromCart(ctx, jwtUserData.UserID, productId); err != nil {
		c.Logger().Errorf("[CartHandler] RemoveFromCart: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	return c.JSON(http.StatusNoContent, response.ResponseSuccess(nil))
}

// RemoveAllFromCart implements [CartHandlerInterface].
func (ch *cartHandler) RemoveAllFromCart(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		jwtUserData entity.JwtUserData
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[CartHandler] RemoveAllFromCart: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[CartHandler] RemoveAllFromCart: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	if err := ch.cartService.RemoveAllFromCart(ctx, jwtUserData.UserID); err != nil {
		c.Logger().Errorf("[CartHandler] RemoveAllFromCart: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	return c.JSON(http.StatusNoContent, response.ResponseSuccess(nil))
}
