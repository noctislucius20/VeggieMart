package adapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"order-service/config"
	"order-service/internal/adapter/handler/response"
	"order-service/internal/core/domain/entity"
	"order-service/internal/core/service"
	"order-service/utils"
	"order-service/utils/conv"
	"order-service/utils/helper"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type MiddlewareAdapterInterface interface {
	CheckToken() echo.MiddlewareFunc
	DistanceCheck() echo.MiddlewareFunc
	RequiredPermission(requiredPermissions ...string) echo.MiddlewareFunc
	IdempotencyCreateOrder() echo.MiddlewareFunc

	haversineDistance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64
}

type middlewareAdapter struct {
	cfg         *config.Config
	jwtService  service.JwtServiceInterface
	redisClient *redis.Client
	logger      *log.Logger
}

type responseRecorder struct {
	io.Writer
	http.ResponseWriter
	body *bytes.Buffer
}

func NewMiddlewareAdapter(cfg *config.Config, logger *log.Logger, jwtService service.JwtServiceInterface, redisClient *redis.Client) MiddlewareAdapterInterface {
	return &middlewareAdapter{
		cfg:         cfg,
		jwtService:  jwtService,
		redisClient: redisClient,
		logger:      logger,
	}
}

func (r *responseRecorder) Write(data []byte) (int, error) {
	r.body.Write(data)
	return r.ResponseWriter.Write(data)
}

// haversineDistance implements [MiddlewareAdapterInterface].
func (m *middlewareAdapter) haversineDistance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	const R = 6371

	lat1Rad := lat1 * math.Pi / 180
	lng1Rad := lng1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lng2Rad := lng2 * math.Pi / 180

	dLat := lat2Rad - lat1Rad
	dLng := lng2Rad - lng1Rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Asin(math.Sqrt(a))

	return R * c
}

// DistanceCheck implements [MiddlewareAdapterInterface].
func (m *middlewareAdapter) DistanceCheck() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			latParam := c.QueryParam("lat")
			lngParam := c.QueryParam("lng")
			if latParam == "" || lngParam == "" {
				err := utils.ErrLatOrLngRequired
				m.logger.Errorf("[MiddlewareAdapter] DistanceCheck: %v", err)
				return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
			}

			lat, lng, err := conv.ParseLatLngToFloat64(latParam, lngParam)
			if err != nil {
				err := utils.ErrLatOrLngInvalid
				m.logger.Errorf("[MiddlewareAdapter] DistanceCheck: %v", err)
				return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
			}

			latRef, lngRef, _ := conv.ParseLatLngToFloat64(m.cfg.App.LatitudeRef, m.cfg.App.LongitudeRef)
			distance := m.haversineDistance(latRef, lngRef, lat, lng)
			if distance > float64(m.cfg.App.MaxDistance) {
				err := utils.ErrDistanceTooFar
				m.logger.Errorf("[MiddlewareAdapter] DistanceCheck: %v", err)
				return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
			}

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

// IdempotencyCreateOrder implements [MiddlewareAdapterInterface].
func (m *middlewareAdapter) IdempotencyCreateOrder() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			idempotencyKey := c.Request().Header.Get("X-Idempotency-Key")
			if idempotencyKey == "" {
				err := utils.ErrIdempotencyKeyRequired
				m.logger.Errorf("[MiddlewareAdapter] IdempotencyCreateOrder: %v", err)
				return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
			}

			redisKey := fmt.Sprintf("idempotency:%s", idempotencyKey)

			success, err := m.redisClient.SetNX(c.Request().Context(), redisKey, "processing", 5*time.Minute).Result()
			if err != nil {
				m.logger.Errorf("[MiddlewareAdapter] IdempotencyCreateOrder: %v", err)
				return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
			}

			if !success {
				cachedJSON, err := m.redisClient.Get(c.Request().Context(), redisKey+":response").Result()
				if err != nil {
					m.logger.Errorf("[MiddlewareAdapter] IdempotencyCreateOrder: %v", err)
					if errors.Is(err, redis.Nil) {
						err := utils.ErrRequestProcessing
						return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
					}
					return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
				}

				var cachedResponse map[string]any
				if err := json.Unmarshal([]byte(cachedJSON), &cachedResponse); err != nil {
					m.logger.Errorf("[MiddlewareAdapter] IdempotencyCreateOrder: %v", err)
					return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
				}
				return c.JSON(http.StatusOK, response.ResponseSuccess(cachedResponse))

			}

			rec := &responseRecorder{ResponseWriter: c.Response().Writer, body: new(bytes.Buffer)}
			c.Response().Writer = rec

			err = next(c)

			if c.Response().Status == 200 || c.Response().Status == 201 {
				m.redisClient.Set(c.Request().Context(), redisKey+":response", rec.body.String(), 5*time.Minute)
			}

			return err
		}
	}
}
