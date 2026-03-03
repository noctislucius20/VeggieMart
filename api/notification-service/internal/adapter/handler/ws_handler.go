package handler

import (
	"net/http"
	"notification-service/config"
	"notification-service/utils"
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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WebSocket implements [WebSocketHandlerInterface].
func (w *webSocketHandler) WebSocket(c echo.Context) error {
	userIdStr := c.QueryParam("user_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, utils.INVALID_ID)
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	utils.AddWebSocketConn(userId, conn)
	defer utils.RemoveWebSocketConn(userId)
	defer conn.Close()

	for {
		if _, _, err := conn.NextReader(); err != nil {
			break
		}
	}

	return nil
}

func NewWebSocketHandler(e *echo.Echo, cfg *config.Config) WebSocketHandlerInterface {
	webSocketHandler := &webSocketHandler{}

	e.Use(middleware.Recover())
	e.GET("/ws", webSocketHandler.WebSocket)

	return webSocketHandler
}
