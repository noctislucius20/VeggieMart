package logger

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type LoggerInterface interface {
	RequestLogger(c echo.Context, v middleware.RequestLoggerValues) error
	Logger() *log.Logger
}

type logger struct {
}

// Logger implements [LoggerInterface].
func (l *logger) Logger() *log.Logger {
	newLog := log.New("product-service")

	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		newLog.SetLevel(log.DEBUG)
	case "WARN":
		newLog.SetLevel(log.WARN)
	case "ERROR":
		newLog.SetLevel(log.ERROR)
	default:
		newLog.SetLevel(log.INFO)
	}
	newLog.SetHeader(
		`${time_rfc3339} | ${level} | ${short_file}:${line} |`,
	)

	return newLog
}

// RequestLogger implements [LoggerInterface].
func (l *logger) RequestLogger(c echo.Context, v middleware.RequestLoggerValues) error {
	fields := fmt.Sprintf(`%s %d %s`, v.Method, v.Status, v.URI)

	switch {
	case v.Status >= 500:
		c.Logger().Error(fields)
	case v.Status >= 400:
		c.Logger().Warn(fields)
	default:
		c.Logger().Info(fields)
	}

	return nil
}

func NewLogger() LoggerInterface {
	logger := new(logger)

	return logger
}
