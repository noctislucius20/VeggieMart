package handler

import (
	"encoding/json"
	"net/http"
	"notification-service/config"
	"notification-service/internal/adapter"
	"notification-service/internal/adapter/handler/response"
	"notification-service/internal/core/domain/entity"
	"notification-service/internal/core/service"
	"notification-service/utils"
	"notification-service/utils/conv"
	"notification-service/utils/logger"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type NotificationHandlerInterface interface {
	GetAllNotifications(c echo.Context) error
	GetNotificationById(c echo.Context) error
	MarkAsReadNotification(c echo.Context) error
}

type notificationHandler struct {
	notificationService service.NotificationServiceInterface
}

// MarkAsReadNotification implements [NotificationHandlerInterface].
func (n *notificationHandler) MarkAsReadNotification(c echo.Context) error {
	var (
		ctx = c.Request().Context()
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[NotificationHandler-1] MarkAsReadNotification: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[NotificationHandler-2] MarkAsReadNotification: %v", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed("id required"))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[NotificationHandler-3] MarkAsReadNotification: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.INVALID_ID))
	}

	if err := n.notificationService.MarkAsReadNotification(ctx, uint(id)); err != nil {
		c.Logger().Errorf("[NotificationHandler-4] MarkAsReadNotification: %v", err)
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(utils.DATA_NOT_FOUND))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}

// GetNotificationById implements [NotificationHandlerInterface].
func (n *notificationHandler) GetNotificationById(c echo.Context) error {
	var (
		ctx              = c.Request().Context()
		respNotification = response.NotificationDetailResponse{}
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[NotificationHandler-1] GetNotificationById: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[NotificationHandler-2] GetNotificationById: %v", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed("id required"))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[NotificationHandler-3] GetNotificationById: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed("id invalid"))
	}

	result, err := n.notificationService.GetNotificationById(ctx, uint(id))
	if err != nil {
		c.Logger().Errorf("[NotificationHandler-4] GetNotificationById: %v", err.Error())
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	respNotification = response.NotificationDetailResponse{
		ID:               result.ID,
		Subject:          *result.Subject,
		Message:          result.Message,
		Status:           result.Status,
		SentAt:           result.SentAt.Format("2006-01-02 15:04:05"),
		ReadAt:           result.ReadAt.Format("2006-01-02 15:04:05"),
		NotificationType: result.NotificationType,
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

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[NotificationHandler-1] GetAllNotifications: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed("data token invalid"))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[NotificationHandler-2] GetAllNotifications: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
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
		c.Logger().Errorf("[NotificationHandler-3] GetAllNotifications: %v", err.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	limit, err := conv.ParseInt64QueryParam(c, "limit", 10)
	if err != nil {
		c.Logger().Errorf("[NotificationHandler-4] GetAllNotifications: %v", err.Error())
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
		c.Logger().Errorf("[NotificationHandler-5] GetAllNotifications: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
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

func NewNotificationHandler(notificationService service.NotificationServiceInterface, e *echo.Echo, cfg *config.Config) NotificationHandlerInterface {
	notificationHandler := notificationHandler{
		notificationService: notificationService,
	}

	e.Use(middleware.Recover())
	e.Use(middleware.ContextTimeoutWithConfig(middleware.ContextTimeoutConfig{
		Timeout: 10 * time.Second,
	}))

	mid := adapter.NewMiddlewareAdapter(cfg, logger.NewLogger().Logger())

	authGroup := e.Group("/auth", mid.CheckToken())
	authGroup.GET("/notifications", notificationHandler.GetAllNotifications)
	authGroup.GET("/notifications/:id", notificationHandler.GetNotificationById)
	authGroup.PUT("/notifications/:id", notificationHandler.MarkAsReadNotification)

	return &notificationHandler
}
