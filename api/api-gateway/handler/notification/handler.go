package notification

import (
	"api-gateway/utils/proxy"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

func RegisterPublicRoutes(g *echo.Group) {
}

func RegisterProtectedRoutes(g *echo.Group) {
	notificationGroup := g.Group("/notifications")

	notificationGroup.GET("", proxyHandler)
	notificationGroup.GET("/ws", wsProxyHandler)
	notificationGroup.GET("/push", proxyHandler)
	notificationGroup.GET("/:id", proxyHandler)
	notificationGroup.PUT("/:id/read", proxyHandler)
	notificationGroup.PUT("/:id/sent", proxyHandler)
	notificationGroup.GET("/admin/ws", wsProxyHandler)
}

func getNotificationServiceUrl() string {
	notificationServiceUrl := os.Getenv("NOTIFICATION_SERVICE_URL")
	if notificationServiceUrl == "" {
		notificationServiceUrl = "http://localhost:8084"
	}
	return notificationServiceUrl
}

func proxyHandler(c echo.Context) error {
	path := strings.TrimPrefix(c.Request().URL.Path, "/api")
	if path == "" {
		path = "/"
	}

	return proxy.ForwardRequest(c, getNotificationServiceUrl()+path)
}

func wsProxyHandler(c echo.Context) error {
	path := strings.TrimPrefix(c.Request().URL.Path, "/api")
	if path == "" {
		path = "/"
	}

	return proxy.ForwardWebSocket(c, getNotificationServiceUrl()+path)
}
