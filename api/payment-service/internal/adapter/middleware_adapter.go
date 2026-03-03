package adapter

import (
	"encoding/json"
	"errors"
	"net/http"
	"payment-service/config"
	"payment-service/internal/adapter/handler/response"
	"payment-service/internal/core/domain/entity"
	"payment-service/utils"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type MiddlewareAdapterInterface interface {
	CheckToken() echo.MiddlewareFunc
}

type middlewareAdapter struct {
	cfg    *config.Config
	logger *log.Logger
}

// CheckToken implements MiddlewareAdapterInterface.
func (m *middlewareAdapter) CheckToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			redisConn := m.cfg.NewRedisClient()

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				err := errors.New(utils.TOKEN_INVALID)
				m.logger.Errorf("[MiddlewareAdapter-1] CheckToken: %v", err.Error())
				return c.JSON(http.StatusUnauthorized, response.ResponseFailed(err.Error()))
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			_, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}

				return []byte(m.cfg.App.JwtSecretKey), nil
			})
			if err != nil {
				err := errors.New(utils.SESSION_EXPIRED)
				m.logger.Errorf("[MiddlewareAdapter-2] CheckToken: %v", err.Error())
				return c.JSON(http.StatusUnauthorized, response.ResponseFailed(err.Error()))
			}

			getSession, err := redisConn.Get(c.Request().Context(), tokenString).Result()
			if err != nil {
				if errors.Is(err, redis.Nil) {
					err := errors.New(utils.TOKEN_INVALID)
					m.logger.Errorf("[MiddlewareAdapter-3] CheckToken: %v", err.Error())
					return c.JSON(http.StatusUnauthorized, response.ResponseFailed(err.Error()))
				}
				m.logger.Errorf("[MiddlewareAdapter-4] CheckToken: %v", err.Error())
				return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
			}

			c.Set("user", getSession)

			jwtUserData := entity.JwtUserData{}
			err = json.Unmarshal([]byte(getSession), &jwtUserData)
			if err != nil {
				m.logger.Errorf("[MiddlewareAdapter-5] CheckToken: %v", err.Error())
				return c.JSON(http.StatusInternalServerError, response.ResponseFailed(err.Error()))
			}

			path := c.Request().URL.Path
			segments := strings.Split(strings.Trim(path, "/"), "/")

			if strings.ToLower(jwtUserData.RoleName) == "customer" && segments[0] == "admin" {
				err := errors.New(utils.ACCESS_FORBIDDEN)
				m.logger.Errorf("[MiddlewareAdapter-6] CheckToken: %v", err.Error())
				return c.JSON(http.StatusForbidden, response.ResponseFailed(err.Error()))
			}

			return next(c)
		}
	}
}

func NewMiddlewareAdapter(cfg *config.Config, logger *log.Logger) MiddlewareAdapterInterface {
	return &middlewareAdapter{
		cfg: cfg, logger: logger,
	}
}
