package logger

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func MiddlewareLog() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			req := c.Request()
			res := c.Response()

			logrus.WithFields(logrus.Fields{
				"method":     req.Method,
				"uri":        req.RequestURI,
				"status":     res.Status,
				"latency":    time.Since(start).String(),
				"user_agent": req.UserAgent(),
				"remote_ip":  c.RealIP(),
			}).Info("HTTP Request")

			return err
		}
	}
}
