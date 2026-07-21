package user

import (
	"api-gateway/utils/proxy"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

func RegisterPublicRoutes(g *echo.Group) {
	userGroup := g.Group("/users")

	userGroup.POST("/signin", proxyHandler)
	userGroup.POST("/signup", proxyHandler)
	userGroup.POST("/forgot-password", proxyHandler)
	userGroup.GET("/activate-account", proxyHandler)
	userGroup.PUT("/reset-password", proxyHandler)
}

func RegisterProtectedRoutes(g *echo.Group) {
	userGroup := g.Group("/users")

	userGroup.GET("/customers", proxyHandler)
	userGroup.POST("/customers/batch", proxyHandler)
	userGroup.GET("/customers/:id", proxyHandler)
	userGroup.POST("/customers", proxyHandler)
	userGroup.PUT("/customers/:id", proxyHandler)
	userGroup.DELETE("/customers/:id", proxyHandler)

	userGroup.GET("/profile", proxyHandler)
	userGroup.PUT("/profile", proxyHandler)
	userGroup.POST("/profile/image-upload", proxyHandler)

	userGroup.GET("/roles", proxyHandler)
	userGroup.GET("/roles/:id", proxyHandler)
	userGroup.POST("/roles", proxyHandler)
	userGroup.PUT("/roles/:id", proxyHandler)
	userGroup.DELETE("/roles/:id", proxyHandler)
}

func proxyHandler(c echo.Context) error {
	userServiceUrl := os.Getenv("USER_SERVICE_URL")
	if userServiceUrl == "" {
		userServiceUrl = "http://localhost:8080"
	}

	path := strings.TrimPrefix(c.Request().URL.Path, "/api")
	if path == "" {
		path = "/"
	}

	return proxy.ForwardRequest(c, userServiceUrl+path)
}
