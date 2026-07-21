package handler

import (
	"api-gateway/handler/notification"
	"api-gateway/handler/order"
	"api-gateway/handler/payment"
	"api-gateway/handler/product"
	"api-gateway/handler/user"

	"github.com/labstack/echo/v4"
)

func RegisterPublicRoutes(g *echo.Group) {
	product.RegisterPublicRoutes(g)
	user.RegisterPublicRoutes(g)
	payment.RegisterPublicRoutes(g)
	notification.RegisterPublicRoutes(g)
}

func RegisterProtectedRoutes(g *echo.Group) {
	user.RegisterProtectedRoutes(g)
	product.RegisterProtectedRoutes(g)
	order.RegisterProtectedRoutes(g)
	payment.RegisterProtectedRoutes(g)
	notification.RegisterProtectedRoutes(g)
}

func RegisterAllRoutes(g *echo.Group, middlewareJwt echo.MiddlewareFunc) {
	RegisterPublicRoutes(g)

	protected := g.Group("")
	protected.Use(middlewareJwt)

	RegisterProtectedRoutes(protected)
}
