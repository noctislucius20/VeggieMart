package product

import (
	"api-gateway/utils/proxy"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

func RegisterPublicRoutes(g *echo.Group) {
	productGroup := g.Group("/products")

	productGroup.GET("/home", proxyHandler)
	productGroup.GET("/home/:id", proxyHandler)
	productGroup.GET("/shop", proxyHandler)

	productGroup.GET("/categories/shop", proxyHandler)
	productGroup.GET("/categories/home", proxyHandler)
}

func RegisterProtectedRoutes(g *echo.Group) {
	productGroup := g.Group("/products")

	productGroup.GET("", proxyHandler)
	productGroup.GET("/:id", proxyHandler)
	productGroup.POST("", proxyHandler)
	productGroup.DELETE("/:id", proxyHandler)
	productGroup.PUT("/:id", proxyHandler)
	productGroup.POST("/batch", proxyHandler)
	productGroup.POST("/stock", proxyHandler)
	productGroup.POST("/image-upload", proxyHandler)

	productGroup.POST("/categories", proxyHandler)
	productGroup.GET("/categories", proxyHandler)
	productGroup.GET("/categories/:id", proxyHandler)
	productGroup.GET("/categories/:slug/slug", proxyHandler)
	productGroup.PUT("/categories/:id", proxyHandler)
	productGroup.DELETE("/categories/:id", proxyHandler)

	productGroup.POST("/carts", proxyHandler)
	productGroup.GET("/carts", proxyHandler)
	productGroup.DELETE("/carts", proxyHandler)
	productGroup.DELETE("/carts/all", proxyHandler)
}

func proxyHandler(c echo.Context) error {
	productServiceUrl := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceUrl == "" {
		productServiceUrl = "http://localhost:8081"
	}

	path := strings.TrimPrefix(c.Request().URL.Path, "/api")
	if path == "" {
		path = "/"
	}

	return proxy.ForwardRequest(c, productServiceUrl+path)
}
