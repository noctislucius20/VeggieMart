package handler

import (
	"encoding/json"
	"net/http"
	"notification-service/config"
	"notification-service/internal/adapter"
	"notification-service/internal/adapter/handler/response"
	"notification-service/internal/core/domain/entity"
	"notification-service/internal/core/service"
	middlewareGateway "notification-service/internal/middleware"
	"notification-service/utils"
	"notification-service/utils/ws"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type WebSocketHandlerInterface interface {
	WebSocket(c echo.Context) error
	AdminWebSocket(c echo.Context) error
}

type webSocketHandler struct {
}

func NewWebSocketHandler(e *echo.Echo, cfg *config.Config, jwtService service.JwtServiceInterface, redisClient *redis.Client, logger *log.Logger) WebSocketHandlerInterface {
	webSocketHandler := &webSocketHandler{}

	e.Use(middleware.Recover())

	mid := adapter.NewMiddlewareAdapter(cfg, logger, jwtService, redisClient)

	adminWsPermission := []string{"notifications:ws:all"}

	notificationGroup := e.Group("/notifications")
	notificationGroup.Use(middlewareGateway.GatewayValidationMiddleware(cfg))

	notificationGroup.GET("/ws", webSocketHandler.WebSocket, mid.CheckToken())
	notificationGroup.GET("/admin/ws", webSocketHandler.AdminWebSocket, mid.CheckToken(), mid.RequiredPermission(adminWsPermission...))

	return webSocketHandler
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WebSocket implements [WebSocketHandlerInterface].
func (w *webSocketHandler) WebSocket(c echo.Context) error {
	var jwtUserData entity.JwtUserData

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[WebsocketHandler] Websocket: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[WebsocketHandler] Websocket: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	userId := jwtUserData.UserID

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Errorf("[WebsocketHandler] Websocket: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	ws.AddWebSocketConn(userId, conn)
	defer ws.RemoveWebSocketConn(userId)
	defer conn.Close()

	for {
		if _, _, err := conn.NextReader(); err != nil {
			break
		}
	}

	return nil
}

// AdminWebSocket implements [WebSocketHandlerInterface].
func (w *webSocketHandler) AdminWebSocket(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Errorf("[WebsocketHandler] AdminWebSocket: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	ws.AddAdminWebSocketConn(conn)
	defer ws.RemoveAdminWebSocketConn(conn)
	defer conn.Close()

	for {
		if _, _, err := conn.NextReader(); err != nil {
			break
		}
	}

	return nil
}
