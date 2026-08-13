package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"payment-service/config"
	"payment-service/internal/adapter/handler/response"
	"payment-service/internal/core/domain/entity"
	"payment-service/internal/core/service"
	"payment-service/utils"
	"payment-service/utils/helper"
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

func NewMiddlewareAdapter(cfg *config.Config, jwtService service.JwtServiceInterface, redisClient *redis.Client, logger *log.Logger) MiddlewareAdapterInterface {
	return &middlewareAdapter{
		cfg:         cfg,
		jwtService:  jwtService,
		redisClient: redisClient,
		logger:      logger,
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

// CheckToken implements MiddlewareAdapterInterface.
func (m *middlewareAdapter) CheckToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				err := utils.ErrTokenInvalid
				m.logger.Errorf("[MiddlewareAdapter] CheckToken: %v", err)
				return c.JSON(http.StatusUnauthorized, response.ResponseFailed(err.Error()))
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

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

			// jwtUserData := entity.JwtUserData{}
			// err = json.Unmarshal([]byte(getSession), &jwtUserData)
			// if err != nil {
			// 	m.logger.Errorf("[MiddlewareAdapter] CheckToken: %v", err)
			// 	return c.JSON(http.StatusInternalServerError, response.ResponseFailed(err.Error()))
			// }

			// path := c.Request().URL.Path
			// segments := strings.Split(strings.Trim(path, "/"), "/")

			// if strings.ToLower(jwtUserData.RoleName) == "customer" && segments[0] == "admin" {
			// 	err := utils.ErrAccessForbidden
			// 	m.logger.Errorf("[MiddlewareAdapter] CheckToken: %v", err)
			// 	return c.JSON(http.StatusForbidden, response.ResponseFailed(err.Error()))
			// }

			return next(c)
		}
	}
}
