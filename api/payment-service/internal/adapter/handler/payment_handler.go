package handler

import (
	"encoding/json"
	"net/http"
	"payment-service/config"
	"payment-service/internal/adapter"
	"payment-service/internal/adapter/handler/request"
	"payment-service/internal/adapter/handler/response"
	"payment-service/internal/core/domain/entity"
	"payment-service/internal/core/service"
	"payment-service/utils"
	"payment-service/utils/conv"
	"payment-service/utils/logger"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type PaymentHandlerInterface interface {
	CreatePayment(c echo.Context) error
	GetAllPayments(c echo.Context) error
	GetPaymentById(c echo.Context) error

	MidtransWebhook(c echo.Context) error

	GetAllPaymentsAdmin(c echo.Context) error
	GetPaymentByIdAdmin(c echo.Context) error
}

type paymentHandler struct {
	paymentService service.PaymentServiceInterface
}

// GetPaymentById implements [PaymentHandlerInterface].
func (p *paymentHandler) GetPaymentById(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		respPayment = response.PaymentDetailResponse{}
		jwtUserData = entity.JwtUserData{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[PaymentHandler-1] GetPaymentById: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[OrderHandler-2] GetPaymentById: %v", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed("id required"))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-3] GetPaymentById: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed("id invalid"))
	}

	result, err := p.paymentService.GetPaymentById(ctx, uint(id), jwtUserData, user)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-4] GetPaymentById: %v", err.Error())
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	respPayment = response.PaymentDetailResponse{
		ID:              int64(result.ID),
		OrderCode:       result.Order.OrderCode,
		PaymentMethod:   result.PaymentMethod,
		PaymentStatus:   result.PaymentStatus,
		GrossAmount:     result.GrossAmount,
		ShippingType:    result.Order.OrderShippingType,
		PaymentAt:       result.PaymentAt,
		OrderAt:         result.Order.OrderAt,
		OrderRemarks:    result.Order.OrderRemarks,
		CustomerName:    result.Customer.CustomerName,
		CustomerAddress: result.Customer.CustomerAddress,
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respPayment))
}

// GetPaymentByIdAdmin implements [PaymentHandlerInterface].
func (p *paymentHandler) GetPaymentByIdAdmin(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		respPayment = response.PaymentDetailResponse{}
		jwtUserData = entity.JwtUserData{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[PaymentHandler-1] GetPaymentByIdAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[OrderHandler-2] GetPaymentByIdAdmin: %v", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed("id required"))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-3] GetPaymentByIdAdmin: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed("id invalid"))
	}

	result, err := p.paymentService.GetPaymentById(ctx, uint(id), jwtUserData, user)
	if err != nil {
		c.Logger().Errorf("[OrderHandler-4] GetPaymentByIdAdmin: %v", err.Error())
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	respPayment = response.PaymentDetailResponse{
		ID:              int64(result.ID),
		OrderCode:       result.Order.OrderCode,
		PaymentMethod:   result.PaymentMethod,
		PaymentStatus:   result.PaymentStatus,
		GrossAmount:     result.GrossAmount,
		ShippingType:    result.Order.OrderShippingType,
		PaymentAt:       result.PaymentAt,
		OrderAt:         result.Order.OrderAt,
		OrderRemarks:    result.Order.OrderRemarks,
		CustomerName:    result.Customer.CustomerName,
		CustomerAddress: result.Customer.CustomerAddress,
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respPayment))
}

// GetAllPayments implements [PaymentHandlerInterface].
func (p *paymentHandler) GetAllPayments(c echo.Context) error {
	var (
		ctx          = c.Request().Context()
		respPayments = []response.PaymentListResponse{}
		jwtUserData  = entity.JwtUserData{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[PaymentHandler-1] GetAllPaymentsAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[PaymentHandler-2] GetAllPaymentsAdmin: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	userId := jwtUserData.UserID

	search := c.QueryParam("search")

	page, err := conv.ParseInt64QueryParam(c, "page", 1)
	if err != nil {
		c.Logger().Errorf("[PaymentHandler-3] GetAllPaymentsAdmin: %v", err.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	limit, err := conv.ParseInt64QueryParam(c, "limit", 10)
	if err != nil {
		c.Logger().Errorf("[PaymentHandler-4] GetAllPaymentsAdmin: %v", err.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	status := c.QueryParam("status")

	orderBy := c.QueryParam("order_by")
	if orderBy == "" {
		orderBy = "created_at"
	}

	orderType := c.QueryParam("order_type")
	if orderType != "asc" && orderType != "desc" {
		orderType = "desc"
	}

	reqEntity := entity.QueryStringPayment{
		Search:    search,
		Page:      page,
		Limit:     limit,
		Status:    status,
		OrderType: orderType,
		OrderBy:   orderBy,
		UserID:    userId,
	}

	results, countData, totalPages, err := p.paymentService.GetAllPayments(ctx, reqEntity, user)
	if err != nil {
		c.Logger().Errorf("[PaymentHandler-5] GetAllPaymentsAdmin: %v", err.Error())
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	for _, result := range results {
		respPayments = append(respPayments, response.PaymentListResponse{
			ID:            uint64(result.ID),
			OrderCode:     result.Order.OrderCode,
			PaymentStatus: result.PaymentStatus,
			PaymentMethod: result.PaymentMethod,
			GrossAmount:   result.GrossAmount,
			ShippingType:  result.Order.OrderShippingType,
		})
	}

	pagination := response.Pagination{
		Page:       page,
		TotalCount: countData,
		PerPage:    limit,
		TotalPage:  totalPages,
	}

	return c.JSON(http.StatusOK, response.ResponseWithPaginationsSuccess(respPayments, pagination))
}

// GetAllPaymentsAdmin implements [PaymentHandlerInterface].
func (p *paymentHandler) GetAllPaymentsAdmin(c echo.Context) error {
	var (
		ctx          = c.Request().Context()
		respPayments = []response.PaymentListResponse{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[PaymentHandler-1] GetAllPaymentsAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed("data token invalid"))
	}

	search := c.QueryParam("search")

	page, err := conv.ParseInt64QueryParam(c, "page", 1)
	if err != nil {
		c.Logger().Errorf("[PaymentHandler-2] GetAllPaymentsAdmin: %v", err.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	limit, err := conv.ParseInt64QueryParam(c, "limit", 10)
	if err != nil {
		c.Logger().Errorf("[PaymentHandler-3] GetAllPaymentsAdmin: %v", err.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	status := c.QueryParam("status")

	orderBy := c.QueryParam("order_by")
	if orderBy == "" {
		orderBy = "created_at"
	}

	orderType := c.QueryParam("order_type")
	if orderType != "asc" && orderType != "desc" {
		orderType = "desc"
	}

	reqEntity := entity.QueryStringPayment{
		Search:    search,
		Page:      page,
		Limit:     limit,
		Status:    status,
		OrderType: orderType,
		OrderBy:   orderBy,
	}

	results, countData, totalPages, err := p.paymentService.GetAllPayments(ctx, reqEntity, user)
	if err != nil {
		c.Logger().Errorf("[PaymentHandler-4] GetAllPaymentsAdmin: %v", err.Error())
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	for _, result := range results {
		respPayments = append(respPayments, response.PaymentListResponse{
			ID:            uint64(result.ID),
			OrderCode:     result.Order.OrderCode,
			PaymentStatus: result.PaymentStatus,
			PaymentMethod: result.PaymentMethod,
			GrossAmount:   result.GrossAmount,
			ShippingType:  result.Order.OrderShippingType,
		})
	}

	pagination := response.Pagination{
		Page:       page,
		TotalCount: countData,
		PerPage:    limit,
		TotalPage:  totalPages,
	}

	return c.JSON(http.StatusOK, response.ResponseWithPaginationsSuccess(respPayments, pagination))
}

// MidtransWebhook implements [PaymentHandlerInterface].
func (p *paymentHandler) MidtransWebhook(c echo.Context) error {
	var (
		ctx                 = c.Request().Context()
		notificationPayload = map[string]any{}
	)

	if err := c.Bind(&notificationPayload); err != nil {
		c.Logger().Errorf("[PaymentHandler-2] MidtransWebhook: %v", err.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	transactionStatus := notificationPayload["transaction_status"].(string)
	orderCode := notificationPayload["order_id"].(string)

	newStatus := ""
	switch transactionStatus {
	case "capture", "settlement":
		newStatus = "SUCCESS"
	case "deny", "cancel", "expire":
		newStatus = "FAILED"
	case "pending":
		newStatus = "PENDING"
	default:
		newStatus = "UNKNOWN"
	}

	if err := p.paymentService.UpdateStatusByOrderCode(ctx, orderCode, newStatus); err != nil {
		c.Logger().Errorf("[PaymentHandler-3] MidtransWebhook: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}

// CreatePayment implements [PaymentHandlerInterface].
func (p *paymentHandler) CreatePayment(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		req = request.PaymentRequest{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[PaymentHandler-1] CreatePayment: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[PaymentHandler-2] CreatePayment: %v", err.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[PaymentHandler-3] CreatePayment: %v", err.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.PaymentEntity{
		OrderID:       req.OrderID,
		PaymentMethod: req.PaymentMethod,
		GrossAmount:   req.GrassAmount,
		UserID:        req.UserID,
		Remarks:       req.Remarks,
	}

	result, err := p.paymentService.ProcessPayment(ctx, reqEntity, user)
	if err != nil {
		c.Logger().Errorf("[PaymentHandler-4] CreatePayment: %v", err.Error())
		if err.Error() == utils.INVALID_PAYMENT_METHOD {
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	respPayment := map[string]any{
		"payment_token": result.PaymentGatewayID,
	}

	return c.JSON(http.StatusCreated, response.ResponseSuccess(respPayment))
}

func NewPaymentHandler(paymentService service.PaymentServiceInterface, e *echo.Echo, cfg *config.Config) PaymentHandlerInterface {
	paymentHandler := &paymentHandler{
		paymentService: paymentService,
	}

	e.Use(middleware.Recover())
	e.Use(middleware.ContextTimeoutWithConfig(middleware.ContextTimeoutConfig{
		Timeout: 10 * time.Second,
	}))

	mid := adapter.NewMiddlewareAdapter(cfg, logger.NewLogger().Logger())

	e.POST("/payments/webhook", paymentHandler.MidtransWebhook)

	authGroup := e.Group("/auth", mid.CheckToken())
	authGroup.POST("/payments", paymentHandler.CreatePayment)
	authGroup.GET("/payments", paymentHandler.GetAllPayments)
	authGroup.GET("/payments/:id", paymentHandler.GetPaymentById)

	adminGroup := e.Group("/admin", mid.CheckToken())
	adminGroup.GET("/payments", paymentHandler.GetAllPaymentsAdmin)
	adminGroup.GET("/payments/:id", paymentHandler.GetPaymentByIdAdmin)

	return paymentHandler
}
