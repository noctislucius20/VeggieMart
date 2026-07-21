package order

import (
	"api-gateway/utils/proxy"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

func RegisterProtectedRoutes(g *echo.Group) {
	orderGroup := g.Group("/orders")

	orderGroup.POST("", proxyHandler)
	orderGroup.GET("", proxyHandler)
	orderGroup.GET("/:id", proxyHandler)
	orderGroup.GET("/:orderCode/code", proxyHandler)
	orderGroup.GET("/admin", proxyHandler)
	orderGroup.GET("/:id/admin", proxyHandler)
	orderGroup.PUT("/:id/status", proxyHandler)
}

func proxyHandler(c echo.Context) error {
	orderServiceUrl := os.Getenv("ORDER_SERVICE_URL")
	if orderServiceUrl == "" {
		orderServiceUrl = "http://localhost:8082"
	}

	path := strings.TrimPrefix(c.Request().URL.Path, "/api")
	if path == "" {
		path = "/"
	}

	return proxy.ForwardRequest(c, orderServiceUrl+path)
}
