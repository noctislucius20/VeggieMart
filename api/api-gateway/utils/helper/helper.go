package helper

import (
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetEnvAsFloat(key string, defaultVal float64) float64 {
	if value := os.Getenv(key); value != "" {
		floatVal, _ := strconv.ParseFloat(value, 64)
		return floatVal
	}

	return defaultVal
}

func GetEnvAsInt(key string, defaultVal int64) int64 {
	if value := os.Getenv(key); value != "" {
		intVal, _ := strconv.ParseInt(value, 10, 64)
		return intVal
	}

	return defaultVal
}

func SetGatewayHeaders(c echo.Context, reqHeader *http.Header) {
	reqHeader.Set("X-API-Gateway", "true")
	reqHeader.Set("X-API-Gateway-Version", "1.0")
	reqHeader.Set("X-Request-ID", c.Response().Header().Get("X-Request-ID"))
	reqHeader.Set("X-Forwarded-For", c.RealIP())
	reqHeader.Set("X-Real-IP", c.RealIP())

	gatewaySecret := os.Getenv("GATEWAY_SECRET_KEY")
	if gatewaySecret != "" {
		reqHeader.Set("X-Gateway-Secret", gatewaySecret)
	}
}
