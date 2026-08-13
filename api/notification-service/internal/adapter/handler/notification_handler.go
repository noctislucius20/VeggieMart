package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"notification-service/config"
	"notification-service/internal/adapter"
	"notification-service/internal/adapter/handler/response"
	"notification-service/internal/core/domain/entity"
	"notification-service/internal/core/service"
	middlewareGateway "notification-service/internal/middleware"
	"notification-service/utils"
	"notification-service/utils/conv"
	"notification-service/utils/logger"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type NotificationHandlerInterface interface {
	GetAllNotifications(c echo.Context) error
	GetAllPushNotification(c echo.Context) error
	GetNotificationById(c echo.Context) error
	MarkAsReadNotification(c echo.Context) error
	MarkAsSentNotification(c echo.Context) error
}

type notificationHandler struct {
	notificationService service.NotificationServiceInterface
}

func NewNotificationHandler(e *echo.Echo, cfg *config.Config, notificationService service.NotificationServiceInterface, jwtService service.JwtServiceInterface, redisClient *redis.Client) NotificationHandlerInterface {
	notificationHandler := notificationHandler{
		notificationService: notificationService,
	}

	e.Use(middleware.Recover())
	e.Use(middleware.ContextTimeoutWithConfig(middleware.ContextTimeoutConfig{
		Timeout: 100 * time.Second,
	}))

	mid := adapter.NewMiddlewareAdapter(cfg, logger.NewLogger().Logger(), jwtService, redisClient)

	authPermission := []string{
		"notifications:read:own",
		"notifications:write:own",
		"notifications:update:own",
		"notifications:delete:own",
	}

	authNotificationGroup := e.Group("/notifications", mid.CheckToken(), mid.RequiredPermission(authPermission...))
	authNotificationGroup.Use(middlewareGateway.GatewayValidationMiddleware(cfg))

	authNotificationGroup.GET("", notificationHandler.GetAllNotifications)
	authNotificationGroup.GET("/push", notificationHandler.GetAllPushNotification)
	authNotificationGroup.GET("/:id", notificationHandler.GetNotificationById)
	authNotificationGroup.PUT("/:id/read", notificationHandler.MarkAsReadNotification)
	authNotificationGroup.PUT("/:id/sent", notificationHandler.MarkAsSentNotification)

	return &notificationHandler
}

// GetAllPushNotification implements [NotificationHandlerInterface].
func (n *notificationHandler) GetAllPushNotification(c echo.Context) error {
	var (
		ctx               = c.Request().Context()
		respNotifications = []response.NotificationResponseList{}
		jwtUserData       = entity.JwtUserData{}
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[NotificationHandler] GetAllPushNotification: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[NotificationHandler] GetAllPushNotification: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
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

	isRead := false
	if isReadStr := c.QueryParam("is_read"); isReadStr != "" {
		isRead, _ = conv.ParseStringToBool(isReadStr)
	}

	page, err := conv.ParseInt64QueryParam(c, "page", 1)
	if err != nil {
		c.Logger().Errorf("[NotificationHandler] GetAllPushNotification: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	limit, err := conv.ParseInt64QueryParam(c, "limit", 5)
	if err != nil {
		c.Logger().Errorf("[NotificationHandler] GetAllPushNotification: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	status := c.QueryParam("status")

	reqEntity := entity.NotificationQueryString{
		Search:    search,
		Status:    status,
		Page:      page,
		Limit:     limit,
		OrderBy:   orderBy,
		OrderType: orderType,
		IsRead:    isRead,
	}

	results, countData, totalPages, err := n.notificationService.GetAllPushNotification(ctx, reqEntity, jwtUserData)
	if err != nil {
		c.Logger().Errorf("[NotificationHandler] GetAllPushNotification: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	for _, result := range results {
		respNotifications = append(respNotifications, response.NotificationResponseList{
			ID:                 result.ID,
			NotificationType:   result.NotificationType,
			NotificationTypeID: result.NotificationTypeID,
			Subject:            *result.Subject,
			Message:            result.Message,
			Status:             result.Status,
			SentAt:             result.SentAt.Format("2006-01-02 15:04:05"),
		})
	}

	pagination := response.Pagination{
		Page:       page,
		TotalCount: countData,
		PerPage:    limit,
		TotalPage:  totalPages,
	}

	return c.JSON(http.StatusOK, response.ResponseWithPaginationsSuccess(respNotifications, pagination))
}

// MarkAsSentNotification implements [NotificationHandlerInterface].
func (n *notificationHandler) MarkAsSentNotification(c echo.Context) error {
	var (
		ctx = c.Request().Context()
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[NotificationHandler] MarkAsSentNotification: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[NotificationHandler] MarkAsSentNotification: %v", utils.ErrIDRequired.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ErrIDRequired.Error()))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[NotificationHandler] MarkAsSentNotification: %v", utils.ErrIDInvalid.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ErrIDInvalid.Error()))
	}

	if err := n.notificationService.MarkAsSentNotification(ctx, int64(id)); err != nil {
		c.Logger().Errorf("[NotificationHandler] MarkAsSentNotification: %v", err)
		if errors.Is(err, utils.ErrDataNotFound) {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(utils.ErrDataNotFound.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}

// MarkAsReadNotification implements [NotificationHandlerInterface].
func (n *notificationHandler) MarkAsReadNotification(c echo.Context) error {
	var (
		ctx = c.Request().Context()
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[NotificationHandler] MarkAsReadNotification: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[NotificationHandler] MarkAsReadNotification: %v", utils.ErrIDRequired.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ErrIDRequired.Error()))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[NotificationHandler] MarkAsReadNotification: %v", utils.ErrIDInvalid.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ErrIDInvalid.Error()))
	}

	if err := n.notificationService.MarkAsReadNotification(ctx, int64(id)); err != nil {
		c.Logger().Errorf("[NotificationHandler] MarkAsReadNotification: %v", err)
		if errors.Is(err, utils.ErrDataNotFound) {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(utils.ErrDataNotFound.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}

// GetNotificationById implements [NotificationHandlerInterface].
func (n *notificationHandler) GetNotificationById(c echo.Context) error {
	var (
		ctx              = c.Request().Context()
		respNotification = response.NotificationDetailResponse{}
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[NotificationHandler] GetNotificationById: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[NotificationHandler] GetNotificationById: %v", utils.ErrIDRequired.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ErrIDRequired.Error()))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[NotificationHandler] GetNotificationById: %v", utils.ErrIDInvalid.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ErrIDInvalid.Error()))
	}

	result, err := n.notificationService.GetNotificationById(ctx, int64(id))
	if err != nil {
		c.Logger().Errorf("[NotificationHandler] GetNotificationById: %v", err)
		if errors.Is(err, utils.ErrDataNotFound) {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	respNotification = response.NotificationDetailResponse{
		ID:                 result.ID,
		Subject:            *result.Subject,
		Message:            result.Message,
		Status:             result.Status,
		SentAt:             result.SentAt.Format("2006-01-02 15:04:05"),
		ReadAt:             result.ReadAt.Format("2006-01-02 15:04:05"),
		NotificationMethod: result.NotificationMethod,
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respNotification))
}

// NotificationDetailResponse implements [NotificationHandlerInterface].
func (n *notificationHandler) GetAllNotifications(c echo.Context) error {
	var (
		ctx               = c.Request().Context()
		respNotifications = []response.NotificationResponseList{}
		jwtUserData       = entity.JwtUserData{}
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[NotificationHandler] GetAllNotifications: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[NotificationHandler] GetAllNotifications: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	userId := jwtUserData.UserID

	search := c.QueryParam("search")

	orderBy := c.QueryParam("order_by")
	if orderBy == "" {
		orderBy = "created_at"
	}

	orderType := c.QueryParam("order_type")
	if orderType != "asc" && orderType != "desc" {
		orderType = "desc"
	}

	isRead := false
	if isReadStr := c.QueryParam("is_read"); isReadStr != "" {
		isRead, _ = conv.ParseStringToBool(isReadStr)
	}

	page, err := conv.ParseInt64QueryParam(c, "page", 1)
	if err != nil {
		c.Logger().Errorf("[NotificationHandler] GetAllNotifications: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	limit, err := conv.ParseInt64QueryParam(c, "limit", 5)
	if err != nil {
		c.Logger().Errorf("[NotificationHandler] GetAllNotifications: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	status := c.QueryParam("status")

	reqEntity := entity.NotificationQueryString{
		Search:    search,
		Status:    status,
		Page:      page,
		Limit:     limit,
		UserID:    userId,
		OrderBy:   orderBy,
		OrderType: orderType,
		IsRead:    isRead,
	}

	results, countData, totalPages, err := n.notificationService.GetAllNotifications(ctx, reqEntity)
	if err != nil {
		c.Logger().Errorf("[NotificationHandler] GetAllNotifications: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	for _, result := range results {
		respNotifications = append(respNotifications, response.NotificationResponseList{
			ID:      result.ID,
			Subject: *result.Subject,
			Status:  result.Status,
			SentAt:  result.SentAt.Format("2006-01-02 15:04:05"),
		})
	}

	pagination := response.Pagination{
		Page:       page,
		TotalCount: countData,
		PerPage:    limit,
		TotalPage:  totalPages,
	}

	return c.JSON(http.StatusOK, response.ResponseWithPaginationsSuccess(respNotifications, pagination))
}
