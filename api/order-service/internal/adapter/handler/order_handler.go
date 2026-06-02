package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/config"
	"order-service/internal/adapter"
	"order-service/internal/adapter/handler/request"
	"order-service/internal/adapter/handler/response"
	"order-service/internal/core/domain/entity"
	"order-service/internal/core/service"
	"order-service/utils"
	"order-service/utils/conv"
	"order-service/utils/logger"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type OrderHandlerInterface interface {
	GetOrderById(c echo.Context) error
	GetAllOrders(c echo.Context) error
	CreateOrder(c echo.Context) error
	GetOrderByOrderCode(c echo.Context) error

	GetAllOrdersAdmin(c echo.Context) error
	GetOrderByIdAdmin(c echo.Context) error
	UpdateOrderStatusByAdmin(c echo.Context) error
}

type orderHandler struct {
	orderService service.OrderServiceInterface
}

func NewOrderHandler(e *echo.Echo, cfg *config.Config, orderService service.OrderServiceInterface, jwtService service.JwtServiceInterface, redisClient *redis.Client) OrderHandlerInterface {
	orderHandler := &orderHandler{
		orderService: orderService,
	}

	e.Use(middleware.Recover())
	e.Use(middleware.ContextTimeoutWithConfig(middleware.ContextTimeoutConfig{
		Timeout: 10 * time.Second,
	}))

	mid := adapter.NewMiddlewareAdapter(cfg, logger.NewLogger().Logger(), jwtService, redisClient)

	authPermission := []string{
		"orders:read:own",
		"orders:write:own",
		"orders:update:own",
		"orders:delete:own",
	}

	adminPermission := []string{
		"orders:read:all",
		"orders:write:all",
		"orders:update:all",
		"orders:delete:all",
	}

	// authGroup := e.Group("/auth", mid.CheckToken())
	e.POST("/orders", orderHandler.CreateOrder, mid.CheckToken(), mid.RequiredPermission(authPermission...), mid.DistanceCheck())
	e.GET("/orders", orderHandler.GetAllOrders, mid.CheckToken(), mid.RequiredPermission(authPermission...))
	e.GET("/orders/:id", orderHandler.GetOrderById, mid.CheckToken(), mid.RequiredPermission(authPermission...))
	e.GET("/orders/:orderCode/code", orderHandler.GetOrderByOrderCode, mid.CheckToken(), mid.RequiredPermission(authPermission...))

	// adminGroup := e.Group("/admin", mid.CheckToken())
	e.GET("/orders/admin", orderHandler.GetAllOrdersAdmin, mid.CheckToken(), mid.RequiredPermission(adminPermission...))
	e.GET("/orders/:id/admin", orderHandler.GetOrderByIdAdmin, mid.CheckToken(), mid.RequiredPermission(adminPermission...))
	e.PUT("/orders/:id/status", orderHandler.UpdateOrderStatusByAdmin, mid.CheckToken(), mid.RequiredPermission(adminPermission...))

	return orderHandler
}

// GetOrderById implements [OrderHandlerInterface].
func (o *orderHandler) GetOrderById(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		respOrder   = response.OrderDetailResponse{}
		jwtUserData = entity.JwtUserData{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[OrderHandler-1] GetOrderById: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[OrderHandler-2] GetOrderById: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[OrderHandler-3] GetOrderById: %v", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ID_REQUIRED))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-4] GetOrderById: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ID_INVALID))
	}

	result, err := o.orderService.GetOrderById(ctx, id, jwtUserData)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-5] GetOrderById: %v", err)
		switch err.Error() {
		case utils.DATA_NOT_FOUND:
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
		}
	}

	respOrder = response.OrderDetailResponse{
		ID:            result.ID,
		OrderCode:     result.OrderCode,
		Status:        result.Status,
		TotalAmount:   result.TotalAmount,
		OrderDatetime: fmt.Sprintf("%s %s", result.OrderDate, result.OrderTime),
		ShippingFee:   result.ShippingFee,
		ShippingType:  result.ShippingType,
		Remarks:       result.Remarks,
		Customer: response.OrderCustomer{
			CustomerID:      result.BuyerID,
			CustomerName:    result.BuyerName,
			CustomerPhone:   result.BuyerPhone,
			CustomerAddress: result.BuyerAddress,
			CustomerEmail:   result.BuyerEmail,
		},
	}

	for _, item := range result.OrderItems {
		respOrder.OrderItems = append(respOrder.OrderItems, response.OrderItemsDetail{
			ProductName:  item.ProductName,
			ProductImage: item.ProductImage,
			ProductPrice: item.Price,
			Quantity:     item.Quantity,
		})
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respOrder))
}

// GetOrderByOrderCode implements [OrderHandlerInterface].
func (o *orderHandler) GetOrderByOrderCode(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		respOrder   = response.OrderDetailResponse{}
		jwtUserData = entity.JwtUserData{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[OrderHandler-1] GetOrderByOrderCode: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[OrderHandler-2] GetOrderByOrderCode: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	orderCode := c.Param("orderCode")
	if orderCode == "" {
		c.Logger().Errorf("[OrderHandler-3] GetOrderByOrderCode: %v", "order code required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ORDER_CODE_REQUIRED))
	}

	result, err := o.orderService.GetOrderByOrderCode(ctx, orderCode, jwtUserData)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-4] GetOrderByOrderCode: %v", err)
		switch err.Error() {
		case utils.DATA_NOT_FOUND:
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
		}
	}

	respOrder = response.OrderDetailResponse{
		ID:            result.ID,
		OrderCode:     result.OrderCode,
		Status:        result.Status,
		TotalAmount:   result.TotalAmount,
		OrderDatetime: fmt.Sprintf("%s %s", result.OrderDate, result.OrderTime),
		ShippingFee:   result.ShippingFee,
		ShippingType:  result.ShippingType,
		Remarks:       result.Remarks,
		Customer: response.OrderCustomer{
			CustomerID:      result.BuyerID,
			CustomerName:    result.BuyerName,
			CustomerPhone:   result.BuyerPhone,
			CustomerAddress: result.BuyerAddress,
			CustomerEmail:   result.BuyerEmail,
		},
	}

	for _, item := range result.OrderItems {
		respOrder.OrderItems = append(respOrder.OrderItems, response.OrderItemsDetail{
			ProductName:  item.ProductName,
			ProductImage: item.ProductImage,
			ProductPrice: item.Price,
			Quantity:     item.Quantity,
		})
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respOrder))
}

// GetAllOrders implements [OrderHandlerInterface].
func (o *orderHandler) GetAllOrders(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		respOrders  = []response.OrderCustomerList{}
		jwtUserData = entity.JwtUserData{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[OrderHandler-1] GetAllOrders: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[OrderHandler-2] GetAllOrders: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	userId := jwtUserData.UserID

	search := c.QueryParam("search")

	page, err := conv.ParseInt64QueryParam(c, "page", 1)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-3] GetAllOrders: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	limit, err := conv.ParseInt64QueryParam(c, "limit", 10)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-4] GetAllOrders: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	status := c.QueryParam("status")

	reqEntity := entity.OrderQueryString{
		Search:  search,
		Page:    page,
		Limit:   limit,
		Status:  status,
		BuyerID: userId,
	}

	results, countData, totalPages, err := o.orderService.GetAllOrders(ctx, reqEntity, user)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-5] GetAllOrders: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	for _, result := range results {
		// productImage := ""
		// for _, item := range result.OrderItems {
		// 	productImage = item.ProductImage
		// }
		respOrders = append(respOrders, response.OrderCustomerList{
			ID:            result.ID,
			OrderCode:     result.OrderCode,
			Status:        result.Status,
			ProductName:   result.OrderItems[0].ProductName,
			TotalAmount:   result.TotalAmount,
			ProductImage:  result.OrderItems[0].ProductImage,
			OrderDatetime: fmt.Sprintf("%s %s", result.OrderDate, result.OrderTime),
		})
	}

	pagination := response.Pagination{
		Page:       page,
		TotalCount: countData,
		PerPage:    limit,
		TotalPage:  totalPages,
	}

	return c.JSON(http.StatusOK, response.ResponseWithPaginationsSuccess(respOrders, pagination))
}

// UpdateOrderStatusByAdmin implements [OrderHandlerInterface].
func (o *orderHandler) UpdateOrderStatusByAdmin(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		req = request.OrderUpdateStatusRequest{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[OrderHandler-1] UpdateOrderStatusByAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[OrderHandler-2] UpdateOrderStatusByAdmin: %v", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ID_REQUIRED))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-3] UpdateOrderStatusByAdmin: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ID_INVALID))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[OrderHandler-4] UpdateOrderStatusByAdmin: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[OrderHandler-5] UpdateOrderStatusByAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.OrderEntity{
		ID:      id,
		Status:  req.Status,
		Remarks: req.Remarks,
	}

	if err := o.orderService.UpdateOrderStatus(ctx, reqEntity, user); err != nil {
		c.Logger().Errorf("[OrderHandler-6] UpdateOrderStatusByAdmin: %v", err)
		switch err.Error() {
		case utils.DATA_NOT_FOUND:
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		case utils.INVALID_STATUS_TRANSITION:
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
		}
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}

// CreateOrder implements [OrderHandlerInterface].
func (o *orderHandler) CreateOrder(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		req         = request.CreateOrderRequest{}
		jwtUserData = entity.JwtUserData{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[OrderHandler-1] CreateOrder: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[OrderHandler-2] CreateOrder: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[OrderHandler-3] CreateOrder: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[OrderHandler-4] CreateOrder: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.OrderEntity{
		BuyerID:      jwtUserData.UserID,
		OrderDate:    req.OrderDate,
		TotalAmount:  req.TotalAmount,
		ShippingType: req.ShippingType,
		Remarks:      req.Remarks,
		OrderTime:    req.OrderTime,
	}

	orderDetails := []entity.OrderItemEntity{}
	for _, item := range req.OrderDetails {
		orderDetails = append(orderDetails, entity.OrderItemEntity{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}
	reqEntity.OrderItems = orderDetails

	orderId, orderCode, err := o.orderService.CreateOrder(ctx, reqEntity, user)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-5] CreateOrder: %v", err)
		switch err.Error() {
		case utils.RELATION_DATA_NOT_FOUND:
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
		case utils.STOCK_UNAVAILABLE:
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
		}
	}

	respOrderId := map[string]any{
		"order_id":   orderId,
		"order_code": orderCode,
	}

	return c.JSON(http.StatusCreated, response.ResponseSuccess(respOrderId))
}

// GetOrderByIdAdmin implements [OrderHandlerInterface].
func (o *orderHandler) GetOrderByIdAdmin(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		respOrder   = response.OrderDetailResponse{}
		jwtUserData = entity.JwtUserData{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[OrderHandler-1] GetOrderByIdAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[OrderHandler-2] GetOrderByIdAdmin: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[OrderHandler-3] GetOrderByIdAdmin: %v", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ID_REQUIRED))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-4] GetOrderByIdAdmin: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ID_INVALID))
	}

	result, err := o.orderService.GetOrderById(ctx, id, jwtUserData)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-5] GetOrderByIdAdmin: %v", err)
		switch err.Error() {
		case utils.DATA_NOT_FOUND:
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
		}
	}

	respOrder = response.OrderDetailResponse{
		ID:            result.ID,
		OrderCode:     result.OrderCode,
		Status:        result.Status,
		TotalAmount:   result.TotalAmount,
		OrderDatetime: fmt.Sprintf("%s %s", result.OrderDate, result.OrderTime),
		ShippingFee:   result.ShippingFee,
		ShippingType:  result.ShippingType,
		Remarks:       result.Remarks,
		Customer: response.OrderCustomer{
			CustomerID:      result.BuyerID,
			CustomerName:    result.BuyerName,
			CustomerPhone:   result.BuyerPhone,
			CustomerAddress: result.BuyerAddress,
			CustomerEmail:   result.BuyerEmail,
		},
	}

	for _, item := range result.OrderItems {
		respOrder.OrderItems = append(respOrder.OrderItems, response.OrderItemsDetail{
			ProductName:  item.ProductName,
			ProductImage: item.ProductImage,
			ProductPrice: item.Price,
			Quantity:     item.Quantity,
		})
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respOrder))
}

// GetAllOrdersAdmin implements [OrderHandlerInterface].
func (o *orderHandler) GetAllOrdersAdmin(c echo.Context) error {
	var (
		ctx        = c.Request().Context()
		respOrders = []response.OrderListResponse{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[OrderHandler-1] GetAllOrdersAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	search := c.QueryParam("search")

	page, err := conv.ParseInt64QueryParam(c, "page", 1)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-2] GetAllOrdersAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	limit, err := conv.ParseInt64QueryParam(c, "limit", 10)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-3] GetAllOrdersAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	status := c.QueryParam("status")

	reqEntity := entity.OrderQueryString{
		Search: search,
		Page:   page,
		Limit:  limit,
		Status: status,
	}

	results, countData, totalPages, err := o.orderService.GetAllOrdersAdmin(ctx, reqEntity, user)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-4] GetAllOrdersAdmin: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	for _, result := range results {
		// productImage := ""
		// for _, item := range result.OrderItems {
		// 	productImage = item.ProductImage
		// }
		respOrders = append(respOrders, response.OrderListResponse{
			ID:           result.ID,
			OrderCode:    result.OrderCode,
			Status:       result.Status,
			TotalAmount:  result.TotalAmount,
			CustomerName: result.BuyerName,
			ProductImage: result.OrderItems[0].ProductImage,
		})
	}

	pagination := response.Pagination{
		Page:       page,
		TotalCount: countData,
		PerPage:    limit,
		TotalPage:  totalPages,
	}

	return c.JSON(http.StatusOK, response.ResponseWithPaginationsSuccess(respOrders, pagination))
}
