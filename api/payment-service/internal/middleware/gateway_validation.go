package middleware

import (
	"net/http"
	"payment-service/config"
	"payment-service/internal/adapter/handler/response"
	"payment-service/utils"

	"github.com/labstack/echo/v4"
)

func GatewayValidationMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requireGateway := cfg.APIGateway.RequestAPI
			if requireGateway == "" {
				requireGateway = "true"
			}

			expectedGateway := cfg.APIGateway.SecretKey

			if requireGateway == "false" {
				return next(c)
			}

			if c.Request().URL.Path == "/health" {
				return next(c)
			}

			gatewayHeader := c.Request().Header.Get("X-API-Gateway")
			if gatewayHeader != "true" {
				c.Logger().Errorf("[GatewayMiddleware-1] GatewayValidationMiddleware: %v", utils.GATEWAY_REQUIRED)
				return c.JSON(http.StatusForbidden, response.ResponseFailed(utils.GATEWAY_REQUIRED))
			}

			if expectedGateway != "" {
				receivedSecret := c.Request().Header.Get("X-Gateway-Secret")
				if receivedSecret != expectedGateway {
					c.Logger().Errorf("[GatewayMiddleware-2] GatewayValidationMiddleware: %v", utils.GATEWAY_SECRET_INVALID)
					return c.JSON(http.StatusForbidden, response.ResponseFailed(utils.GATEWAY_SECRET_INVALID))
				}
			}

			gatewayVersion := c.Request().Header.Get("X-API-Gateway-Version")
			if gatewayVersion == "" {
				c.Logger().Warn("missing X-API-Gateway-Version header")
			}

			secretStatus := "not configured"
			if expectedGateway != "" {
				secretStatus = "validated"
			}

			c.Logger().Infof("request from API Gateway - Version: %s, Request-ID: %s, Secret: %s", gatewayVersion, c.Request().Header.Get("X-Request-ID"), secretStatus)

			return next(c)
		}
	}
}
