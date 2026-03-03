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

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type OrderHandlerInterface interface {
	GetOrderById(c echo.Context) error
	GetAllOrders(c echo.Context) error
	CreateOrder(c echo.Context) error
	GetBatchOrders(c echo.Context) error
	GetOrderByOrderCode(c echo.Context) error

	GetAllOrdersAdmin(c echo.Context) error
	GetOrderByIdAdmin(c echo.Context) error
	UpdateOrderStatusByAdmin(c echo.Context) error

	GetOrderIdByOrderCodePublic(c echo.Context) error
}

type orderHandler struct {
	orderService service.OrderServiceInterface
}

// GetOrderIdByOrderCodePublic implements [OrderHandlerInterface].
func (o *orderHandler) GetOrderIdByOrderCodePublic(c echo.Context) error {
	var (
		ctx = c.Request().Context()
	)

	orderCodeParam := c.Param("orderCode")
	if orderCodeParam == "" {
		c.Logger().Errorf("[OrderHandler-1] GetOrderIdByOrderCodePublic: %v", "order code required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed("order code required"))
	}

	result, err := o.orderService.GetOrderIdByOrderCodePublic(ctx, orderCodeParam)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-2] GetOrderIdByOrderCodePublic: %v", err.Error())
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	respOrder := map[string]int64{
		"order_id": result,
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respOrder))
}

// GetBatchOrders implements [OrderHandlerInterface].
func (o *orderHandler) GetBatchOrders(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		respBatch   = []response.OrderBatchResponse{}
		req         = request.OrderBatchRequest{}
		jwtUserData = entity.JwtUserData{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[OrderHandler-1] GetBatchOrders: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[OrderHandler-2] GetBatchOrders: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[OrderHandler-3] GetBatchOrders: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[OrderHandler-4] GetBatchOrders: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	results, err := o.orderService.GetBatchOrders(ctx, req.IDOrders, jwtUserData, user)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-5] GetBatchOrders: %v", err)
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(utils.DATA_NOT_FOUND))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(err.Error()))
	}

	for _, result := range results {
		respBatch = append(respBatch, response.OrderBatchResponse{
			ID:           result.ID,
			OrderCode:    result.OrderCode,
			ShippingType: result.ShippingType,
			Customer: response.OrderCustomer{
				CustomerName:  result.BuyerName,
				CustomerEmail: result.BuyerEmail,
			},
		})
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respBatch))
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

	userId := jwtUserData.UserID

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[OrderHandler-3] GetOrderById: %v", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed("id required"))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-4] GetOrderById: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed("id invalid"))
	}

	result, err := o.orderService.GetOrderById(ctx, id, userId, user)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-5] GetOrderById: %v", err.Error())
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
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
			CustomerID:      result.ID,
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
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed("data token invalid"))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[OrderHandler-2] GetOrderByOrderCode: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	orderCode := c.Param("orderCode")
	if orderCode == "" {
		c.Logger().Errorf("[OrderHandler-3] GetOrderByOrderCode: %v", "order code required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed("order code required"))
	}

	result, err := o.orderService.GetOrderByOrderCode(ctx, orderCode, jwtUserData, user)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-4] GetOrderByOrderCode: %v", err.Error())
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	respOrder = response.OrderDetailResponse{
		ID:            result.ID,
		OrderCode:     result.OrderCode,
		Status:        result.Status,
		TotalAmount:   result.TotalAmount,
		OrderDatetime: fmt.Sprintf("%s %s", result.OrderDate, result.OrderTime),
		ShippingFee:   result.ShippingFee,
		Remarks:       result.Remarks,
		Customer: response.OrderCustomer{
			CustomerID:      result.ID,
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
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed("data token invalid"))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[OrderHandler-2] GetAllOrders: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	userId := jwtUserData.UserID

	search := c.QueryParam("search")

	page, err := conv.ParseInt64QueryParam(c, "page", 1)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-3] GetAllOrders: %v", err.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	limit, err := conv.ParseInt64QueryParam(c, "limit", 10)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-4] GetAllOrders: %v", err.Error())
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
		c.Logger().Errorf("[OrderHandler-5] GetAllOrders: %v", err.Error())
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
			Weight:        result.OrderItems[0].ProductWeight,
			Unit:          result.OrderItems[0].ProductUnit,
			Quantity:      result.OrderItems[0].Quantity,
			OrderDateTime: result.OrderDate,
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
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed("data token invalid"))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[OrderHandler-2] UpdateOrderStatusByAdmin: %v", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed("id required"))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-3] UpdateOrderStatusByAdmin: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed("id invalid"))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[OrderHandler-4] UpdateOrderStatusByAdmin: %v", err.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[OrderHandler-5] UpdateOrderStatusByAdmin: %v", err.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.OrderEntity{
		ID:      id,
		Status:  req.Status,
		Remarks: req.Remarks,
	}

	if err := o.orderService.UpdateOrderStatus(ctx, reqEntity, user); err != nil {
		c.Logger().Errorf("[OrderHandler-6] UpdateOrderStatusByAdmin: %v", err)
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(utils.DATA_NOT_FOUND))
		}
		if err.Error() == utils.INVALID_STATUS_TRANSITION {
			return c.JSON(http.StatusConflict, response.ResponseFailed(utils.INVALID_STATUS_TRANSITION))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}

// CreateOrder implements [OrderHandlerInterface].
func (o *orderHandler) CreateOrder(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		req = request.CreateOrderRequest{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[OrderHandler-1] CreateOrder: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[OrderHandler-2] CreateOrder: %v", err.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[OrderHandler-3] CreateOrder: %v", err.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.OrderEntity{
		BuyerID:      req.BuyerID,
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
		c.Logger().Errorf("[OrderHandler-4] CreateOrder: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
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
		ctx       = c.Request().Context()
		respOrder = response.OrderDetailResponse{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[OrderHandler-1] GetOrderByIdAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed("data token invalid"))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[OrderHandler-2] GetOrderByIdAdmin: %v", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed("id required"))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-3] GetOrderByIdAdmin: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed("id invalid"))
	}

	result, err := o.orderService.GetOrderByIdAdmin(ctx, id, user)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-4] GetOrderByIdAdmin: %v", err.Error())
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	respOrder = response.OrderDetailResponse{
		ID:            result.ID,
		OrderCode:     result.OrderCode,
		Status:        result.Status,
		TotalAmount:   result.TotalAmount,
		OrderDatetime: fmt.Sprintf("%s %s", result.OrderDate, result.OrderTime),
		ShippingFee:   result.ShippingFee,
		Remarks:       result.Remarks,
		Customer: response.OrderCustomer{
			CustomerID:      result.ID,
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
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed("data token invalid"))
	}

	search := c.QueryParam("search")

	page, err := conv.ParseInt64QueryParam(c, "page", 1)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-2] GetAllOrdersAdmin: %v", err.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	limit, err := conv.ParseInt64QueryParam(c, "limit", 10)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-3] GetAllOrdersAdmin: %v", err.Error())
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
		c.Logger().Errorf("[OrderHandler-4] GetAllOrdersAdmin: %v", err.Error())
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

func NewOrderHandler(orderService service.OrderServiceInterface, e *echo.Echo, cfg *config.Config) OrderHandlerInterface {
	orderHandler := &orderHandler{
		orderService: orderService,
	}

	e.Use(middleware.Recover())
	e.Use(middleware.ContextTimeoutWithConfig(middleware.ContextTimeoutConfig{
		Timeout: 10 * time.Second,
	}))

	mid := adapter.NewMiddlewareAdapter(cfg, logger.NewLogger().Logger())

	e.GET("/public/orders/:orderCode/code", orderHandler.GetOrderIdByOrderCodePublic)

	authGroup := e.Group("/auth", mid.CheckToken())
	authGroup.POST("/orders", orderHandler.CreateOrder, mid.DistanceCheck())
	authGroup.GET("/orders", orderHandler.GetAllOrders)
	authGroup.POST("/orders/batch", orderHandler.GetBatchOrders)
	authGroup.GET("/orders/:id", orderHandler.GetOrderById)
	authGroup.GET("/orders/:orderCode/code", orderHandler.GetOrderByOrderCode)

	adminGroup := e.Group("/admin", mid.CheckToken())
	adminGroup.GET("/orders", orderHandler.GetAllOrdersAdmin)
	adminGroup.GET("/orders/:id", orderHandler.GetOrderByIdAdmin)
	adminGroup.PUT("/orders/:id/status", orderHandler.UpdateOrderStatusByAdmin)

	return orderHandler
}
