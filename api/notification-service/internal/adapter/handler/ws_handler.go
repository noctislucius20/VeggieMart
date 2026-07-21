package handler

import (
	"net/http"
	"notification-service/config"
	"notification-service/internal/adapter/handler/response"
	middlewareGateway "notification-service/internal/middleware"
	"notification-service/utils"
	"notification-service/utils/ws"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type WebSocketHandlerInterface interface {
	WebSocket(c echo.Context) error
}

type webSocketHandler struct {
}

func NewWebSocketHandler(e *echo.Echo, cfg *config.Config) WebSocketHandlerInterface {
	webSocketHandler := &webSocketHandler{}

	e.Use(middleware.Recover())

	notificationGroup := e.Group("/notifications")
	notificationGroup.Use(middlewareGateway.GatewayValidationMiddleware(cfg))

	notificationGroup.GET("/ws", webSocketHandler.WebSocket)

	return webSocketHandler
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WebSocket implements [WebSocketHandlerInterface].
func (w *webSocketHandler) WebSocket(c echo.Context) error {
	userIdStr := c.QueryParam("user_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.Logger().Errorf("[WebsocketHandler-1] Websocket: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.INVALID_ID))
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Errorf("[WebsocketHandler-1] Websocket: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
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
