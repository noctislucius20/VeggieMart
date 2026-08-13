package jwt

import (
	"api-gateway/response"
	"api-gateway/utils"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type middlewareJwt struct {
	redisClient *redis.Client
}

type MiddlewareJWTInterface interface {
	CheckToken() echo.MiddlewareFunc

	validateToken(encodeToken string) (*jwt.Token, error)
}

func NewMiddlewareJWT(redisClient *redis.Client) MiddlewareJWTInterface {
	return &middlewareJwt{
		redisClient: redisClient,
	}
}

// MiddlewareJWT implements [MiddlewareJWTInterface].
func (m *middlewareJwt) CheckToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().URL.Path == "/health" {
				return next(c)
			}

			var tokenString string

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			} else if c.IsWebSocket() {
				tokenString = c.QueryParam("token")
			}

			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
			}

			_, err := m.validateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenExpired.Error()))
			}

			return next(c)
		}
	}
}

// validateToken implements [MiddlewareJWTInterface].
func (m *middlewareJwt) validateToken(encodeToken string) (*jwt.Token, error) {
	return jwt.Parse(encodeToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		secretKey := os.Getenv("JWT_SECRET_KEY")
		if secretKey == "" {
			secretKey = "your-secret-key"
		}

		return []byte(secretKey), nil
	})
}
