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
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
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

func NewPaymentHandler(paymentService service.PaymentServiceInterface, e *echo.Echo, cfg *config.Config, jwtService service.JwtServiceInterface, redisClient *redis.Client) PaymentHandlerInterface {
	paymentHandler := &paymentHandler{
		paymentService: paymentService,
	}

	e.Use(middleware.Recover())
	e.Use(middleware.ContextTimeoutWithConfig(middleware.ContextTimeoutConfig{
		Timeout: 10 * time.Second,
	}))

	mid := adapter.NewMiddlewareAdapter(cfg, jwtService, redisClient, logger.NewLogger().Logger())

	e.POST("/payments/webhook", paymentHandler.MidtransWebhook)

	authPermission := []string{
		"payments:read:own",
		"payments:write:own",
		"payments:update:own",
		"payments:delete:own",
	}

	adminPermission := []string{
		"payments:read:all",
		"payments:write:all",
		"payments:update:all",
		"payments:delete:all",
	}

	// authGroup := e.Group("/auth", mid.CheckToken())
	e.POST("/payments", paymentHandler.CreatePayment, mid.CheckToken(), mid.RequiredPermission(authPermission...))
	e.GET("/payments", paymentHandler.GetAllPayments, mid.CheckToken(), mid.RequiredPermission(authPermission...))
	e.GET("/payments/:id", paymentHandler.GetPaymentById, mid.CheckToken(), mid.RequiredPermission(authPermission...))

	// adminGroup := e.Group("/admin", mid.CheckToken())
	e.GET("/payments/admin", paymentHandler.GetAllPaymentsAdmin, mid.CheckToken(), mid.RequiredPermission(adminPermission...))
	e.GET("/payments/:id/admin", paymentHandler.GetPaymentByIdAdmin, mid.CheckToken(), mid.RequiredPermission(adminPermission...))

	return paymentHandler
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

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[PaymentHandler-2] GetPaymentById: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[PaymentHandler-3] GetPaymentById: %v", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ID_REQUIRED))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[PaymentHandler-4] GetPaymentById: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ID_INVALID))
	}

	result, err := p.paymentService.GetPaymentById(ctx, id, jwtUserData, user)
	if err != nil {
		c.Logger().Errorf("[PaymentHandler-5] GetPaymentById: %v", err)
		switch err.Error() {
		case utils.DATA_NOT_FOUND:
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
		}
	}

	respPayment = response.PaymentDetailResponse{
		ID:              int64(result.ID),
		OrderCode:       result.Order.OrderCode,
		PaymentMethod:   result.PaymentMethod,
		PaymentStatus:   result.PaymentStatus,
		GrossAmount:     result.GrossAmount,
		ShippingType:    result.Order.ShippingType,
		PaymentAt:       result.PaymentAt,
		OrderAt:         result.Order.OrderDatetime,
		OrderRemarks:    result.Order.Remarks,
		CustomerName:    result.Customer.CustomerName,
		CustomerEmail:   result.Customer.CustomerEmail,
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

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[PaymentHandler-2] GetPaymentById: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[PaymentHandler-3] GetPaymentByIdAdmin: %v", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ID_REQUIRED))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[PaymentHandler-4] GetPaymentByIdAdmin: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ID_INVALID))
	}

	result, err := p.paymentService.GetPaymentById(ctx, id, jwtUserData, user)
	if err != nil {
		c.Logger().Errorf("[PaymentHandler-5] GetPaymentByIdAdmin: %v", err)
		switch err.Error() {
		case utils.DATA_NOT_FOUND:
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
		}
	}

	respPayment = response.PaymentDetailResponse{
		ID:              int64(result.ID),
		PaymentMethod:   result.PaymentMethod,
		PaymentStatus:   result.PaymentStatus,
		GrossAmount:     result.GrossAmount,
		PaymentAt:       result.PaymentAt,
		OrderCode:       result.Order.OrderCode,
		ShippingType:    result.Order.ShippingType,
		OrderAt:         result.Order.OrderDatetime,
		OrderRemarks:    result.Order.Remarks,
		CustomerName:    result.Customer.CustomerName,
		CustomerEmail:   result.Customer.CustomerEmail,
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
		c.Logger().Errorf("[PaymentHandler-1] GetAllPayments: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[PaymentHandler-2] GetAllPayments: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	userId := jwtUserData.UserID

	search := c.QueryParam("search")

	page, err := conv.ParseInt64QueryParam(c, "page", 1)
	if err != nil {
		c.Logger().Errorf("[PaymentHandler-3] GetAllPayments: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	limit, err := conv.ParseInt64QueryParam(c, "limit", 10)
	if err != nil {
		c.Logger().Errorf("[PaymentHandler-4] GetAllPayments: %v", err)
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
		c.Logger().Errorf("[PaymentHandler-5] GetAllPayments: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	for _, result := range results {
		respPayments = append(respPayments, response.PaymentListResponse{
			ID:            int64(result.ID),
			OrderCode:     result.Order.OrderCode,
			PaymentStatus: result.PaymentStatus,
			PaymentMethod: result.PaymentMethod,
			GrossAmount:   result.GrossAmount,
			ShippingType:  result.Order.ShippingType,
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
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	search := c.QueryParam("search")

	page, err := conv.ParseInt64QueryParam(c, "page", 1)
	if err != nil {
		c.Logger().Errorf("[PaymentHandler-2] GetAllPaymentsAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	limit, err := conv.ParseInt64QueryParam(c, "limit", 10)
	if err != nil {
		c.Logger().Errorf("[PaymentHandler-3] GetAllPaymentsAdmin: %v", err)
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
		c.Logger().Errorf("[PaymentHandler-4] GetAllPaymentsAdmin: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	for _, result := range results {
		respPayments = append(respPayments, response.PaymentListResponse{
			ID:            int64(result.ID),
			OrderCode:     result.Order.OrderCode,
			PaymentStatus: result.PaymentStatus,
			PaymentMethod: result.PaymentMethod,
			GrossAmount:   result.GrossAmount,
			ShippingType:  result.Order.ShippingType,
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
		c.Logger().Errorf("[PaymentHandler-2] MidtransWebhook: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	transactionStatus := notificationPayload["transaction_status"].(string)
	orderCode := notificationPayload["order_id"].(string)

	newStatus := ""
	switch strings.ToLower(transactionStatus) {
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
		c.Logger().Errorf("[PaymentHandler-3] MidtransWebhook: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}

// CreatePayment implements [PaymentHandlerInterface].
func (p *paymentHandler) CreatePayment(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		req         = request.PaymentRequest{}
		jwtUserData = entity.JwtUserData{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[PaymentHandler-1] CreatePayment: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[PaymentHandler-2] CreatePayment: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[PaymentHandler-3] CreatePayment: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[PaymentHandler-4] CreatePayment: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.PaymentEntity{
		PaymentMethod: req.PaymentMethod,
		GrossAmount:   req.GrassAmount,
		UserID:        jwtUserData.UserID,
		Remarks:       req.Remarks,
		Order: entity.OrderEntity{
			ID: int64(req.OrderID),
		},
	}

	result, err := p.paymentService.ProcessPayment(ctx, reqEntity, jwtUserData)
	if err != nil {
		c.Logger().Errorf("[PaymentHandler-5] CreatePayment: %v", err)
		switch err.Error() {
		case utils.RELATION_DATA_NOT_FOUND:
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
		case utils.INVALID_PAYMENT_METHOD:
			return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
		}
	}

	respPayment := map[string]any{
		"payment_id":    result.ID,
		"payment_token": result.PaymentGatewayID,
	}

	return c.JSON(http.StatusCreated, response.ResponseSuccess(respPayment))
}
