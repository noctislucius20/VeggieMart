package payment

import (
	"api-gateway/utils/proxy"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

func RegisterPublicRoutes(g *echo.Group) {
	paymentGroup := g.Group("/payments")

	paymentGroup.POST("/webhook", proxyHandler)
}

func RegisterProtectedRoutes(g *echo.Group) {
	paymentGroup := g.Group("/payments")

	paymentGroup.POST("", proxyHandler)
	paymentGroup.GET("", proxyHandler)
	paymentGroup.GET("/order/:order_id", proxyHandler)
	paymentGroup.GET("/:id", proxyHandler)
	paymentGroup.GET("/admin", proxyHandler)
	paymentGroup.GET("/order/:order_id/admin", proxyHandler)
	paymentGroup.GET("/:id/admin", proxyHandler)
}

func proxyHandler(c echo.Context) error {
	paymentServiceUrl := os.Getenv("PAYMENT_SERVICE_URL")
	if paymentServiceUrl == "" {
		paymentServiceUrl = "http://localhost:8083"
	}

	path := strings.TrimPrefix(c.Request().URL.Path, "/api")
	if path == "" {
		path = "/"
	}

	return proxy.ForwardRequest(c, paymentServiceUrl+path)
}
