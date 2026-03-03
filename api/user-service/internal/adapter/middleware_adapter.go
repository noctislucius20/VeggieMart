package adapter

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"user-service/config"
	"user-service/internal/adapter/handler/response"
	"user-service/internal/core/domain/entity"
	"user-service/internal/core/service"
	"user-service/utils"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type MiddlewareAdapterInterface interface {
	CheckToken() echo.MiddlewareFunc
}

type middlewareAdapter struct {
	cfg         *config.Config
	jwtService  service.JwtServiceInterface
	logger      *log.Logger
	redisClient *redis.Client
}

// CheckToken implements MiddlewareAdapterInterface.
func (m *middlewareAdapter) CheckToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				err := errors.New(utils.TOKEN_INVALID)
				m.logger.Errorf("[MiddlewareAdapter-1] CheckToken: %v", err.Error())
				return c.JSON(http.StatusUnauthorized, response.ResponseFailed(err.Error()))
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			_, err := m.jwtService.ValidateToken(tokenString)
			if err != nil {
				err := errors.New(utils.SESSION_EXPIRED)
				m.logger.Errorf("[MiddlewareAdapter-2] CheckToken: %v", err.Error())
				return c.JSON(http.StatusUnauthorized, response.ResponseFailed(err.Error()))
			}

			getSession, err := m.redisClient.Get(c.Request().Context(), tokenString).Result()
			if err != nil {
				m.logger.Errorf("[MiddlewareAdapter-3] CheckToken: %v", err.Error())
				if errors.Is(err, redis.Nil) {
					err := errors.New(utils.TOKEN_INVALID)
					return c.JSON(http.StatusUnauthorized, response.ResponseFailed(err.Error()))
				}
				return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
			}

			c.Set("user", getSession)

			jwtUserData := entity.JwtUserData{}
			err = json.Unmarshal([]byte(getSession), &jwtUserData)
			if err != nil {
				m.logger.Errorf("[MiddlewareAdapter-4] CheckToken: %v", err.Error())
				return c.JSON(http.StatusInternalServerError, response.ResponseFailed(err.Error()))
			}

			path := c.Request().URL.Path
			segments := strings.Split(strings.Trim(path, "/"), "/")

			if strings.ToLower(jwtUserData.RoleName) == "Customer" && segments[0] == "admin" {
				err := errors.New(utils.ACCESS_FORBIDDEN)
				m.logger.Errorf("[MiddlewareAdapter-5] CheckToken: %v", err.Error())
				return c.JSON(http.StatusForbidden, response.ResponseFailed(err.Error()))
			}

			return next(c)
		}
	}
}

func NewMiddlewareAdapter(cfg *config.Config, logger *log.Logger, jwtService service.JwtServiceInterface, redisClient *redis.Client) MiddlewareAdapterInterface {
	return &middlewareAdapter{
		cfg:         cfg,
		jwtService:  jwtService,
		redisClient: redisClient,
		logger:      logger,
	}
}
