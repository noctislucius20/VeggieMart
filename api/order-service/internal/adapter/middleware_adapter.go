package adapter

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"order-service/config"
	"order-service/internal/adapter/handler/response"
	"order-service/internal/core/domain/entity"
	"order-service/utils"
	"order-service/utils/conv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type MiddlewareAdapterInterface interface {
	CheckToken() echo.MiddlewareFunc
	DistanceCheck() echo.MiddlewareFunc
	haversineDistance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64
}

type middlewareAdapter struct {
	cfg    *config.Config
	logger *log.Logger
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
				m.logger.Errorf("[MiddlewareAdapter-1] DistanceCheck: %s", "lat or lng required")
				return c.JSON(http.StatusBadRequest, response.ResponseFailed("lat or lng required"))
			}

			lat, lng, err := conv.ParseLatLngToFloat64(latParam, lngParam)
			if err != nil {
				m.logger.Errorf("[MiddlewareAdapter-2] DistanceCheck: %s", "lat or lng invalid")
				return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed("lat or lng invalid"))
			}

			latRef, lngRef, _ := conv.ParseLatLngToFloat64(m.cfg.App.LatitudeRef, m.cfg.App.LongitudeRef)
			distance := m.haversineDistance(latRef, lngRef, lat, lng)
			if distance > float64(m.cfg.App.MaxDistance) {
				m.logger.Errorf("[MiddlewareAdapter-3] DistanceCheck: %s", "distance too far")
				return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed("distance too far"))
			}

			return next(c)
		}
	}
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
