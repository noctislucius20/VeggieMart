package middleware

import (
	"net/http"
	"product-service/config"
	"product-service/internal/adapter/handler/response"
	"product-service/utils"
	"product-service/utils/helper"

	"github.com/labstack/echo/v4"
)

func InternalServiceMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			serviceName := c.Request().Header.Get("X-Service-Name")
			serviceSecret := c.Request().Header.Get("X-Service-Secret")

			if !helper.IsServiceAllowed(serviceName) {
				return c.JSON(http.StatusForbidden, response.ResponseFailed(utils.ErrServiceNotAllowed.Error()))
			}

			if serviceSecret != cfg.APIInternalService.SecretKey {
				return c.JSON(http.StatusForbidden, response.ResponseFailed(utils.ErrServiceSecretInvalid.Error()))
			}

			c.Logger().Infof("request from internal service: %s", serviceName)

			return next(c)
		}
	}
}
