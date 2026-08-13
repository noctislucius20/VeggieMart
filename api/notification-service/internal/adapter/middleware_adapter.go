package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"notification-service/config"
	"notification-service/internal/adapter/handler/response"
	"notification-service/internal/core/domain/entity"
	"notification-service/internal/core/service"
	"notification-service/utils"
	"notification-service/utils/helper"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type MiddlewareAdapterInterface interface {
	CheckToken() echo.MiddlewareFunc

	RequiredPermission(requiredPermissions ...string) echo.MiddlewareFunc
}

type middlewareAdapter struct {
	cfg         *config.Config
	jwtService  service.JwtServiceInterface
	redisClient *redis.Client
	logger      *log.Logger
}

func NewMiddlewareAdapter(cfg *config.Config, logger *log.Logger, jwtService service.JwtServiceInterface, redisClient *redis.Client) MiddlewareAdapterInterface {
	return &middlewareAdapter{
		cfg:         cfg,
		jwtService:  jwtService,
		redisClient: redisClient,
		logger:      logger,
	}
}

// CheckToken implements MiddlewareAdapterInterface.
func (m *middlewareAdapter) CheckToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var tokenString string

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			} else if c.IsWebSocket() {
				tokenString = c.QueryParam("token")
			}

			if tokenString == "" {
				err := utils.ErrTokenInvalid
				m.logger.Errorf("[MiddlewareAdapter] CheckToken: %v", err)
				return c.JSON(http.StatusUnauthorized, response.ResponseFailed(err.Error()))
			}

			_, err := m.jwtService.ValidateToken(tokenString)
			if err != nil {
				err := utils.ErrSessionExpired
				m.logger.Errorf("[MiddlewareAdapter] CheckToken: %v", err)
				return c.JSON(http.StatusUnauthorized, response.ResponseFailed(err.Error()))
			}

			keyIdxSession := fmt.Sprintf("user:session:%s", tokenString)
			getIdxSession, err := m.redisClient.Get(c.Request().Context(), keyIdxSession).Result()
			if err != nil {
				m.logger.Errorf("[MiddlewareAdapter] CheckToken: %v", err)
				if errors.Is(err, redis.Nil) {
					err := utils.ErrTokenInvalid
					return c.JSON(http.StatusUnauthorized, response.ResponseFailed(err.Error()))
				}
				return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
			}

			keySession := fmt.Sprintf("user:id:%s:session", getIdxSession)
			getSession, err := m.redisClient.Get(c.Request().Context(), keySession).Result()
			if err != nil {
				m.logger.Errorf("[MiddlewareAdapter] CheckToken: %v", err)
				return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
			}

			c.Set("user", getSession)

			return next(c)
		}
	}
}

// RequiredPermission implements [MiddlewareAdapterInterface].
func (m *middlewareAdapter) RequiredPermission(requiredPermissions ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var (
				jwtUserData entity.JwtUserData
				roleEntity  entity.RoleEntity
				permissions []string
			)

			user, ok := c.Get("user").(string)
			if !ok || user == "" {
				err := utils.ErrTokenInvalid
				m.logger.Errorf("[MiddlewareAdapter] RequiredPermission: %v", err)
				return c.JSON(http.StatusUnauthorized, response.ResponseFailed(err.Error()))
			}

			if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
				m.logger.Errorf("[MiddlewareAdapter] RequiredPermission: %v", err)
				return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
			}

			keyRolePermission := fmt.Sprintf("role:id:%d", jwtUserData.RoleID)
			rolePermission, err := m.redisClient.Get(c.Request().Context(), keyRolePermission).Result()
			if err != nil {
				m.logger.Errorf("[MiddlewareAdapter] RequiredPermission: %v", err)
				if errors.Is(err, redis.Nil) {
					return c.JSON(http.StatusForbidden, response.ResponseFailed(utils.ErrAccessForbidden.Error()))
				}
				return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
			}

			if err := json.Unmarshal([]byte(rolePermission), &roleEntity); err != nil {
				m.logger.Errorf("[MiddlewareAdapter] RequiredPermission: %v", err)
				return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
			}

			for _, p := range roleEntity.Permissions {
				permissions = append(permissions, fmt.Sprintf("%s:%s:%s", p.Resource, p.Action, p.Scope))
			}

			if allowed := helper.HasRequiredPermissions(permissions, requiredPermissions); !allowed {
				err := utils.ErrAccessForbidden
				m.logger.Errorf("[MiddlewareAdapter] RequiredPermission: %v", err)
				return c.JSON(http.StatusForbidden, response.ResponseFailed(err.Error()))
			}

			return next(c)
		}
	}
}
